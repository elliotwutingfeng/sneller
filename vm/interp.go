// Copyright (C) 2023 Sneller, Inc.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"encoding/binary"
	"math/bits"
	"unsafe"

	"golang.org/x/sys/cpu"

	"github.com/SnellerInc/sneller/ion"
)

type opfn func(bc *bytecode, pc int) int

// ops, arg slots, etc. are encoded as 16-bit integers:
func bcword(buf []byte, pc int) uint {
	return uint(binary.LittleEndian.Uint16(buf[pc:]))
}

func bcword64(buf []byte, pc int) uint64 {
	return binary.LittleEndian.Uint64(buf[pc:])
}

func slotcast[T any](b *bytecode, slot uint) *T {
	buf := unsafe.Slice((*byte)(unsafe.Pointer(&b.vstack[0])), len(b.vstack)*8)
	ptr := unsafe.Pointer(&buf[slot])
	return (*T)(ptr)
}

func argptr[T any](b *bytecode, pc int) *T {
	return slotcast[T](b, bcword(b.compiled, pc))
}

func setvmrefB(dst *bRegData, src []vmref) {
	*dst = bRegData{}
	for i := 0; i < min(bcLaneCount, len(src)); i++ {
		dst.offsets[i] = src[i][0]
		dst.sizes[i] = src[i][1]
	}
}

func setvmref(dst *vRegData, src []vmref) {
	*dst = vRegData{}
	for i := 0; i < min(bcLaneCount, len(src)); i++ {
		dst.offsets[i] = src[i][0]
		dst.sizes[i] = src[i][1]
		if dst.sizes[i] == 0 {
			continue
		}
		mem := src[i].mem()
		dst.typeL[i] = byte(mem[0])
		dst.headerSize[i] = byte(ion.HeaderSizeOf(mem))
	}
}

func getTLVSize(length uint) uint {
	if length < 14 {
		return 1
	}

	if length < (1 << 7) {
		return 2
	}

	if length < (1 << 14) {
		return 3
	}

	if length < (1 << 21) {
		return 4
	}

	if length < (1 << 28) {
		return 5
	}

	return 6
}

func encodeSymbol(dst []byte, offset int, symbol ion.Symbol) int {
	if symbol < (1 << 7) {
		dst[offset+0] = byte(symbol) | 0x80
		return 1
	}

	if symbol < (1 << 14) {
		dst[offset+0] = byte((symbol >> 7) & 0x7F)
		dst[offset+1] = byte(symbol&0xFF) | 0x80
		return 2
	}

	if symbol < (1 << 21) {
		dst[offset+0] = byte((symbol >> 14) & 0x7F)
		dst[offset+1] = byte((symbol >> 7) & 0x7F)
		dst[offset+2] = byte(symbol&0xFF) | 0x80
		return 3
	}

	if symbol < (1 << 28) {
		dst[offset+0] = byte((symbol >> 21) & 0x7F)
		dst[offset+1] = byte((symbol >> 14) & 0x7F)
		dst[offset+2] = byte((symbol >> 7) & 0x7F)
		dst[offset+3] = byte(symbol&0xFF) | 0x80
		return 4
	}

	panic("encodeSymbol: symbol ID out of range")
}

func encodeTLVUnsafe(dst []byte, offset int, valueType ion.Type, length uint) int {
	tag := byte(valueType) << 4
	if length < 14 {
		dst[offset] = tag | byte(length)
		return 1
	}

	dst[offset] = tag | 0xE

	if length < (1 << 7) {
		dst[offset+1] = byte(length) | 0x80
		return 2
	}

	if length < (1 << 14) {
		dst[offset+1] = byte((length >> 7) & 0x7F)
		dst[offset+2] = byte(length&0xFF) | 0x80
		return 3
	}

	if length < (1 << 21) {
		dst[offset+1] = byte((length >> 14) & 0x7F)
		dst[offset+2] = byte((length >> 7) & 0x7F)
		dst[offset+3] = byte(length&0xFF) | 0x80
		return 4
	}

	if length < (1 << 28) {
		dst[offset+1] = byte((length >> 21) & 0x7F)
		dst[offset+2] = byte((length >> 14) & 0x7F)
		dst[offset+3] = byte((length >> 7) & 0x7F)
		dst[offset+4] = byte(length&0xFF) | 0x80
		return 5
	}

	panic("encodeTLVUnsafe: length too large")
}

func evalfiltergo(bc *bytecode, delims []vmref) int {
	i, j := 0, 0
	for i < len(delims) {
		next := delims[i:]
		next = next[:min(len(next), bcLaneCount)]
		apos := bc.auxpos
		mask := evalfiltergolanes(bc, next)
		if bc.err != 0 {
			return j
		}
		// compress delims + auxvals
		for k := range next {
			if mask == 0 {
				break
			}
			if (mask & 1) != 0 {
				for l := range bc.auxvals {
					bc.auxvals[l][j] = bc.auxvals[l][apos+k]
				}
				delims[j] = next[k]
				j++
			}
			mask >>= 1
		}
		i += len(next)
	}
	return j
}

func bcauxvalgo(bc *bytecode, pc int) int {
	dstv := argptr[vRegData](bc, pc)
	dstk := argptr[kRegData](bc, pc+2)
	auxv := bcword(bc.compiled, pc+4)
	lst := bc.auxvals[auxv][bc.auxpos:]
	lst = lst[:min(bcLaneCount, len(lst))]
	mask := uint16(0)
	for i := range lst {
		if lst[i][1] != 0 {
			mask |= (1 << i)
		}
	}
	setvmref(dstv, lst)
	dstk.mask = mask
	return pc + 6
}

func bcinitgo(bc *bytecode, pc int) int {
	delims := argptr[bRegData](bc, pc+0)
	delims.offsets = bc.vmState.delims.offsets
	delims.sizes = bc.vmState.delims.sizes

	bc.err = 0
	mask := argptr[kRegData](bc, pc+2)
	mask.mask = bc.vmState.validLanes.mask
	return pc + 4
}

func bcfindsymgo(bc *bytecode, pc int) int {
	dstv := argptr[vRegData](bc, pc)
	dstk := argptr[kRegData](bc, pc+2)
	srcb := argptr[bRegData](bc, pc+4)
	symbol, _, _ := ion.ReadLabel(bc.compiled[pc+6:])
	srck := argptr[kRegData](bc, pc+10)

	src := *srcb // may alias output
	srcmask := srck.mask
	retmask := uint16(0)

outer:
	for i := 0; i < bcLaneCount; i++ {
		start := src.offsets[i]
		width := src.sizes[i]
		dstv.offsets[i] = start
		dstv.sizes[i] = 0
		dstv.typeL[i] = 0
		dstv.headerSize[i] = 0
		if srcmask&(1<<i) == 0 {
			continue
		}
		mem := vmref{start, width}.mem()
		var sym ion.Symbol
		var err error
	symsearch:
		for len(mem) > 0 {
			sym, mem, err = ion.ReadLabel(mem)
			if err != nil {
				bc.err = bcerrCorrupt
				break outer
			}
			if sym > symbol {
				break symsearch
			}
			dstv.offsets[i] = start + width - uint32(len(mem))
			dstv.sizes[i] = uint32(ion.SizeOf(mem))
			dstv.typeL[i] = byte(mem[0])
			dstv.headerSize[i] = byte(ion.HeaderSizeOf(mem))
			if sym == symbol {
				retmask |= (1 << i)
				break symsearch
			}
			mem = mem[ion.SizeOf(mem):]
		}
	}
	dstk.mask = retmask
	return pc + 12
}

func bccmplti64immgo(bc *bytecode, pc int) int {
	dstk := argptr[kRegData](bc, pc)
	arg0 := argptr[i64RegData](bc, pc+2)
	arg1imm := int64(bcword64(bc.compiled, pc+4))
	srck := argptr[kRegData](bc, pc+12)

	mask := srck.mask
	retmask := uint16(0)
	for i := 0; i < bcLaneCount; i++ {
		if arg0.values[i] < arg1imm {
			retmask |= (1 << i)
		}
	}
	dstk.mask = retmask & mask
	return pc + 14
}

func bccmpvi64immgo(bc *bytecode, pc int) int {
	retslot := argptr[i64RegData](bc, pc)
	retk := argptr[kRegData](bc, pc+2)
	argv := argptr[vRegData](bc, pc+4)
	imm := int64(bcword64(bc.compiled, pc+6))
	argk := argptr[kRegData](bc, pc+14)

	src := *argv // copied since argv may alias retslot
	mask := argk.mask
	retmask := uint16(0)
	for i := 0; i < bcLaneCount; i++ {
		if mask&(1<<i) == 0 {
			retslot.values[i] = 0
			continue
		}
		start := src.offsets[i]
		width := src.sizes[i]
		if width == 0 {
			retslot.values[i] = 0
			continue
		}
		mem := vmref{start, width}.mem()
		var rv int64
		switch ion.Type(src.typeL[i] >> 4) {
		case ion.FloatType:
			f, _, _ := ion.ReadFloat64(mem)
			if f < float64(imm) {
				rv = -1
			} else if f > float64(imm) {
				rv = 1
			}
			retmask |= (1 << i)
		case ion.IntType:
			j, _, _ := ion.ReadInt(mem)
			if j < imm {
				rv = -1
			} else if j > imm {
				rv = 1
			}
			retmask |= (1 << i)
		case ion.UintType:
			u, _, _ := ion.ReadUint(mem)
			if imm < 0 || u < uint64(imm) {
				rv = -1
			} else if u > uint64(imm) {
				rv = 1
			}
			retmask |= (1 << i)
		}
		retslot.values[i] = rv
	}
	retk.mask = retmask
	return pc + 16
}

func bctuplego(bc *bytecode, pc int) int {
	dstb := argptr[bRegData](bc, pc)
	dstk := argptr[kRegData](bc, pc+2)
	srcv := argptr[vRegData](bc, pc+4)
	srck := argptr[kRegData](bc, pc+6)

	src := *srcv
	mask := srck.mask
	retmask := uint16(0)
	for i := 0; i < bcLaneCount; i++ {
		dstb.offsets[i] = 0
		dstb.sizes[i] = 0
		if mask&(1<<i) == 0 ||
			ion.Type(src.typeL[i]>>4) != ion.StructType ||
			src.sizes[i] == 0 {
			continue
		}
		hdrsize := uint32(src.headerSize[i])
		dstb.offsets[i] = src.offsets[i] + hdrsize
		dstb.sizes[i] = src.sizes[i] - hdrsize
		retmask |= (1 << i)
	}
	dstk.mask = retmask
	return pc + 8
}

/*
// TODO: Temporarily disabled as it regresses one test for some reason.
func bcboxf64go(bc *bytecode, pc int) int {
	dst := argptr[vRegData](bc, pc)
	src := argptr[f64RegData](bc, pc+2)
	mask := argptr[kRegData](bc, pc+4).mask

	p := len(bc.scratch)
	want := 9 * bcLaneCount
	if cap(bc.scratch)-p < want {
		bc.err = bcerrMoreScratch
		return pc + 6
	}
	bc.scratch = bc.scratch[:p+want]
	mem := bc.scratch[p:]
	var buf ion.Buffer
	for i := 0; i < bcLaneCount; i++ {
		if mask&(1<<i) == 0 {
			dst.offsets[i] = 0
			dst.sizes[i] = 0
			dst.typeL[i] = 0
			dst.headerSize[i] = 0
			continue
		}
		buf.Set(mem[:0])
		buf.WriteCanonicalFloat(src.values[i])
		start, ok := vmdispl(mem)
		if !ok {
			panic("bad scratch buffer")
		}
		dst.offsets[i] = start
		dst.sizes[i] = uint32(len(buf.Bytes()))
		dst.typeL[i] = mem[0]
		dst.headerSize[i] = 1 // ints and floats always have 1-byte headers
		mem = mem[9:]
	}
	return pc + 6
}
*/

func bcretgo(bc *bytecode, pc int) int {
	bc.err = 0
	bc.vmState.outputLanes.mask = bc.vmState.validLanes.mask
	bc.auxpos += bits.OnesCount16(bc.vmState.outputLanes.mask)
	return pc
}

func bcretkgo(bc *bytecode, pc int) int {
	k := argptr[kRegData](bc, pc+0)
	bc.vmState.outputLanes.mask = k.mask
	bc.auxpos += bits.OnesCount16(bc.vmState.validLanes.mask)
	return pc + 2
}

func bcretbkgo(bc *bytecode, pc int) int {
	b := argptr[bRegData](bc, pc+0)
	k := argptr[kRegData](bc, pc+2)

	bc.vmState.delims.offsets = b.offsets
	bc.vmState.delims.sizes = b.sizes
	bc.vmState.outputLanes.mask = k.mask
	bc.auxpos += bits.OnesCount16(bc.vmState.validLanes.mask)
	return pc + 4
}

func bcretbhkgo(bc *bytecode, pc int) int {
	b := argptr[bRegData](bc, pc+0)
	k := argptr[kRegData](bc, pc+4)

	bc.vmState.delims.offsets = b.offsets
	bc.vmState.delims.sizes = b.sizes
	bc.vmState.outputLanes.mask = k.mask
	bc.auxpos += bits.OnesCount16(bc.vmState.validLanes.mask)

	return pc + 6
}

func bcretskgo(bc *bytecode, pc int) int {
	s := argptr[sRegData](bc, pc+0)
	k := argptr[kRegData](bc, pc+2)
	bc.vmState.sreg = *s
	bc.vmState.outputLanes = *k
	bc.auxpos += bits.OnesCount16(bc.vmState.validLanes.mask)
	return pc + 4
}

func init() {
	opinfo[opinit].portable = bcinitgo
	opinfo[opret].portable = bcretgo
	opinfo[opretk].portable = bcretkgo
	opinfo[opretbk].portable = bcretbkgo
	opinfo[opretsk].portable = bcretskgo
	opinfo[opretbhk].portable = bcretbhkgo
	opinfo[opauxval].portable = bcauxvalgo
	opinfo[opfindsym].portable = bcfindsymgo
	opinfo[opcmpvi64imm].portable = bccmpvi64immgo
	opinfo[opcmplti64imm].portable = bccmplti64immgo
	opinfo[optuple].portable = bctuplego
	opinfo[opbroadcast0k].portable = bcbroadcast0kgo
	opinfo[opbroadcast1k].portable = bcbroadcast1kgo
	opinfo[opfalse].portable = bcfalsego
	opinfo[opnotk].portable = bcnotkgo
	opinfo[opandk].portable = bcandkgo
	opinfo[opandnk].portable = bcandnkgo
	opinfo[opork].portable = bcorkgo
	opinfo[opxork].portable = bcxorkgo
	opinfo[opxnork].portable = bcxnorkgo
	// opinfo[opboxf64].portable = bcboxf64go
}

func evalfindgo(bc *bytecode, delims []vmref, stride int) {
	stack := bc.vstack
	var alt bytecode
	bc.scratch = bc.scratch[:len(bc.savedlit)] // reset scratch ONCE, here
	// convert stride to 64-bit words:
	stride = stride / int(unsafe.Sizeof(bc.vstack[0]))
	for len(delims) > 0 {
		mask := uint16(0xffff)
		lanes := bcLaneCount
		if len(delims) < lanes {
			mask >>= bcLaneCount - len(delims)
			lanes = len(delims)
		}
		bc.err = 0
		bc.vmState.validLanes.mask = mask
		bc.vmState.outputLanes.mask = mask
		setvmrefB(&bc.vmState.delims, delims)
		eval(bc, &alt, false)
		if bc.err != 0 {
			return
		}
		delims = delims[lanes:]
		bc.vstack = bc.vstack[stride:]
	}
	bc.vstack = stack
}

func evalsplatgo(bc *bytecode, indelims, outdelims []vmref, perm []int32) (int, int) {
	ipos, opos := 0, 0
	var alt bytecode
	for ipos < len(indelims) && opos < len(outdelims) {
		next := indelims[ipos:]
		mask := uint16(0xffff)
		if len(next) < bcLaneCount {
			mask >>= bcLaneCount - len(next)
		}
		setvmrefB(&bc.vmState.delims, indelims[ipos:])
		bc.vmState.validLanes.mask = mask
		bc.vmState.outputLanes.mask = 0
		eval(bc, &alt, true)
		if bc.err != 0 {
			return 0, 0
		}
		retmask := bc.vmState.outputLanes.mask
		output := opos
		lanes := min(bcLaneCount, len(next))
		for i := 0; i < lanes; i++ {
			if (retmask & (1 << i)) == 0 {
				continue
			}
			start := bc.vmState.sreg.offsets[i]
			width := bc.vmState.sreg.sizes[i]
			slice := vmref{start, width}.mem()
			for len(slice) > 0 {
				if output == len(outdelims) || output == len(perm) {
					// need to return early
					return ipos, opos
				}
				s := ion.SizeOf(slice)
				outdelims[output] = vmref{start, uint32(s)}
				perm[output] = int32(i + ipos)
				output++
				slice = slice[s:]
				start += uint32(s)
			}
		}
		// checkpoint splat
		opos = output
		ipos += lanes
	}
	return ipos, opos
}

func evalfiltergolanes(bc *bytecode, delims []vmref) uint16 {
	if len(delims) > bcLaneCount {
		panic("invalid len(delims) for evalfiltergolanes")
	}
	mask := uint16(0xffff)
	mask >>= bcLaneCount - len(delims)
	var alt bytecode
	setvmrefB(&bc.vmState.delims, delims)
	bc.vmState.validLanes.mask = mask
	bc.vmState.outputLanes.mask = 0
	eval(bc, &alt, true)
	if bc.err != 0 {
		return 0
	}
	return bc.vmState.outputLanes.mask
}

func evalprojectgo(bc *bytecode, delims []vmref, dst []byte, symbols []syminfo) (int, int) {
	offset := 0
	capacity := len(dst)
	rowsProcessed := 0

	dstDisp, ok := vmdispl(dst)
	if !ok {
		return 0, 0
	}

	var alt bytecode

	for rowsProcessed < len(delims) {
		initialDstLength := offset
		n := min(len(delims)-rowsProcessed, bcLaneCount)

		for lane := 0; lane < n; lane++ {
			bc.vmState.delims.offsets[lane] = uint32(delims[rowsProcessed+lane][0])
			bc.vmState.delims.sizes[lane] = uint32(delims[rowsProcessed+lane][1])
		}

		mask := uint16((1 << n) - 1)
		bc.err = 0
		bc.vmState.validLanes.mask = mask

		eval(bc, &alt, true)

		if bc.err != 0 {
			return initialDstLength, rowsProcessed
		}

		mask = bc.vmState.outputLanes.mask

		for lane := 0; lane < n; lane++ {
			if (mask & (uint16(1) << lane)) == 0 {
				continue
			}

			contentSize := uint(0)
			for i := 0; i < len(symbols); i++ {
				// Only account for the value if the value is not MISSING.
				vals := slotcast[vRegData](bc, uint(symbols[i].slot))
				if symbols[i].size != 0 && vals.sizes[lane] != 0 {
					contentSize += uint(symbols[i].size) + uint(vals.sizes[lane])
				}
			}

			headerSize := getTLVSize(contentSize)
			structSize := headerSize + contentSize

			if uint(capacity-offset) < structSize {
				return initialDstLength, rowsProcessed
			}

			offset += encodeTLVUnsafe(dst, offset, ion.StructType, contentSize)

			// Update delims with the projection output
			delims[rowsProcessed+lane][0] = uint32(dstDisp + uint32(offset))
			delims[rowsProcessed+lane][1] = uint32(contentSize)

			for i := 0; i < len(symbols); i++ {
				// Only serialize Key+Value if the value is not MISSING.
				vals := slotcast[vRegData](bc, uint(symbols[i].slot))
				if symbols[i].size != 0 && vals.sizes[lane] != 0 {
					valOffset := vals.offsets[lane]
					valLength := vals.sizes[lane]

					offset += encodeSymbol(dst, offset, symbols[i].value)
					offset += copy(dst[offset:], vmm[valOffset:valOffset+valLength])
				}
			}
		}

		rowsProcessed += n
	}

	return offset, rowsProcessed
}

func evaldedupgo(bc *bytecode, delims []vmref, hashes []uint64, tree *radixTree64, slot int) int {
	var alt bytecode
	indelims := delims
	dout := 0
	for len(indelims) > 0 {
		n := min(len(indelims), bcLaneCount)
		mask := uint16(0xffff)
		if n < bcLaneCount {
			mask >>= bcLaneCount - n
		}
		setvmrefB(&bc.vmState.delims, indelims[:n])
		bc.vmState.validLanes.mask = mask
		bc.vmState.outputLanes.mask = 0
		apos := bc.auxpos
		eval(bc, &alt, true)
		if bc.err != 0 {
			return 0
		}
		outhashes := slotcast[hRegData](bc, uint(slot))
		outmask := bc.vmState.outputLanes.mask
		for i := 0; i < n; i++ {
			if outmask&(1<<i) == 0 || tree.Offset(outhashes.lo[i]) >= 0 {
				continue // lane not active or tree contains hash already
			}
			delims[dout] = indelims[i]
			hashes[dout] = outhashes.lo[i]
			// compress auxvals
			for j, lst := range bc.auxvals {
				bc.auxvals[j][dout] = lst[apos+i]
			}
			dout++
		}
		indelims = indelims[n:]
	}
	return dout
}

func evalaggregatego(bc *bytecode, delims []vmref, aggregateDataBuffer []byte) int {
	var alt bytecode
	ret := 0
	bc.vmState.aggPtr = unsafe.Pointer(&aggregateDataBuffer[0])
	for len(delims) > 0 {
		n := min(len(delims), bcLaneCount)
		mask := uint16(0xffff)
		if n < bcLaneCount {
			mask >>= bcLaneCount - n
		}
		setvmrefB(&bc.vmState.delims, delims)
		bc.vmState.validLanes.mask = mask
		bc.vmState.outputLanes.mask = 0
		eval(bc, &alt, true)
		if bc.err != 0 {
			return ret
		}
		delims = delims[n:]
		ret += n
	}
	return ret
}

//go:noescape
func bcenter(bc *bytecode, k7 uint16)

// eval evaluates bc and uses alt as scratch space
// for evaluating unimplemented opcodes via the assembly interpreter
func eval(bc, alt *bytecode, resetScratch bool) {
	l := len(bc.compiled)
	pc := 0
	if resetScratch {
		bc.scratch = bc.scratch[:len(bc.savedlit)]
	}
	for pc < l && bc.err == 0 {
		op := bcop(bcword(bc.compiled, pc))
		pc += 2
		fn := opinfo[op].portable
		if fn != nil {
			pc = fn(bc, pc)
		} else if cpu.X86.HasAVX512 {
			pc = runSingle(bc, alt, pc, bc.vmState.validLanes.mask)
		} else {
			bc.err = bcerrNotSupported
			break
		}
	}
}

// run a single bytecode instruction @ pc
func runSingle(bc, alt *bytecode, pc int, k7 uint16) int {
	// copy over everything except the compiled bytestream:
	compiled := alt.compiled
	*alt = *bc

	alt.compiled = compiled[:0]

	// create a new compiled bytestream with the single instr + return
	width := bcwidth(bc, pc-2)
	alt.compiled = append(alt.compiled, bc.compiled[pc-2:pc+width]...)
	alt.compiled = append(alt.compiled, byte(opret), byte(opret>>8))

	// evaluate the bytecode and copy back the error state
	bcenter(alt, k7)
	bc.err = alt.err
	bc.errpc = alt.errpc
	bc.errinfo = alt.errinfo
	bc.scratch = alt.scratch // copy back len(scratch)
	return pc + width
}
