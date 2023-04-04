// Copyright (C) 2022 Sneller, Inc.
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

// Polymorphic comparison that works at a value level and results in [-1, 0, 1] outputs. It can be
// used to compare multiple data types in different lanes. In general the implementation does the
// following comparisons, in the respective order:
//
//   - NULL/BOOL values
//   - NUMBER values (both I64 and F64)
//   - STRING and TIMESTAMP values (comparison treat these the same as only bytes need to be compared)
//
// To perform the comparison, an ION data type is translated to an Internal type ID, which has two
// purposes:
//
//   - Define a sorting order (comparison of non-comparable values is simply `Order(B) - Order(A)`
//   - Define whether the comparison has sorting semantics or not
//
// Bit pattern:
//
//   - 0x0F - ORDERING rules (used by sorting comparisons only)
//   - 0x30 - IGNORED at the moment, set to the same value as 0x40 bit
//   - 0x40 - NON-COMPARABLE bit, represents types that cannot be compared
//   - 0x80 - SORTING SEMANTICS enabled if ALL bits in the predicate are 1

// Matching types only (this predicate doesn't provide sorting semantics)
CONST_DATA_U64(cmpv_predicate_matching_type, 0, $0x0403000202020100)
CONST_DATA_U64(cmpv_predicate_matching_type, 8, $0x7F7F7F7F7F7F7F04)
CONST_GLOBAL(cmpv_predicate_matching_type, $16)

// Compare with sorting semantics, NULLs are sorted before any other value.
CONST_DATA_U64(cmpv_predicate_sort_nulls_first, 0, $0x8483808282828180)
CONST_DATA_U64(cmpv_predicate_sort_nulls_first, 8, $0xFFFFFFFFFFFFFF84)
CONST_GLOBAL(cmpv_predicate_sort_nulls_first, $16)

// Compare with sorting semantics, NULLs are sorted after any other value.
CONST_DATA_U64(cmpv_predicate_sort_nulls_last, 0, $0x848380828282818F)
CONST_DATA_U64(cmpv_predicate_sort_nulls_last, 8, $0xFFFFFFFFFFFFFF84)
CONST_GLOBAL(cmpv_predicate_sort_nulls_last, $16)

// Input registers:
//
//   - K1  - input predicate
//   - SI  - [L] value base offset (this is a VIRT_BASE register)
//   - R8  - [R] value base offset (can be the same as VIRT_BASE)
//   - Z10 - [L] unpacked ION value offsets
//   - Z11 - [L] unpacked ION value lengths
//   - Z12 - [R] unpacked ION value offsets
//   - Z13 - [R] unpacked ION value lengths
//   - Z14 - [L] TLV byte, having zeros in unused bytes
//   - Z15 - [R] TLV byte, having zeros in unused bytes
//   - Z16 - comparison predicate, see cmpv_predicate_... for more details
//   - Z30 - predicate(bswap64)
//   - Z31 - constant(-1)
//
// Output registers:
//
//   - K1  - output predicate
//   - Z16 - comparison results as 32-bit clamped values to [-1, 0, 1].
//   - Z30 - predicate(bswap64)
//   - Z31 - constant(-1)
//
// Preserved registers:
//
//   - Z0:Z9 - Unchanged
//   - Z30   - predicate(bswap64)
//   - Z31   - constant(-1)
//
// Purpose of some registers:
//
//   - K1  - final predicate, having masked out lanes that couldn't be compared
//   - K2  - working predicate, always contains remaining lanes to compare
//   - Z24 - [L] value 8 content bytes (low)
//   - Z25 - [L] value 8 content bytes (high)
//   - Z26 - [R] value 8 content bytes (low)
//   - Z27 - [R] value 8 content bytes (high)
//
// Implementation Notes:
//
//   - Initial 8 content bytes of both left and right values are gathered before
//     jumping to value-specialized compare implementations. The reason is to hide
//     the latency of VPGATHERDQ as much as possible and to basically have the data
//     ready when needed.
TEXT fncmpv(SB), NOSPLIT|NOFRAME, $0
  KMOVB K1, K3
  VPXORD X24, X24, X24
  KSHIFTRW $8, K1, K2
  VPGATHERDQ 0(VIRT_BASE)(Y10*1), K3, Z24              // Z24 <- [L] value 8 content bytes (low)
  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

  VPXORD X26, X26, X26
  KMOVB K1, K3
  KMOVB K2, K4
  VEXTRACTI32X8 $1, Z10, Y20
  VEXTRACTI32X8 $1, Z12, Y21
  VPSRLD $4, Z14, Z18                                  // Z18 <- [L] ION type
  VPSRLD $4, Z15, Z19                                  // Z19 <- [R] ION type
  VPGATHERDQ 0(R8)(Y12*1), K3, Z26                     // Z26 <- [R] value 8 content bytes (low)
  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

  VPSHUFB Z18, Z16, Z22                                // Z22 <- [L] ION type converted to internal type
  VPSHUFB Z19, Z16, Z23                                // Z23 <- [R] ION type converted to internal type
  VPXORD X25, X25, X25
  VPABSD Z31, Z28                                      // Z28 <- dword(1)
  VPORD.Z Z23, Z22, K1, Z29                            // Z29 <- [L] and right internal types combined (for predicate calculations)
  VPSUBD Z23, Z22, Z16                                 // Z16 <- initial comparison results (with sorting semantics, at this point)
  VPGATHERDQ 0(VIRT_BASE)(Y20*1), K2, Z25              // Z25 <- [L] value 8 content bytes (high)
  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

  VPTESTNMD.BCST CONSTD_0x40(), Z29, K1, K1            // K1 <- updated predicate to only contain lanes that can be compared
  VPTESTNMD Z16, Z16, K1, K2                           // K2 <- lanes having compatible type, which means they can be value compared
  KANDNW K1, K2, K3                                    // K3 <- lanes not having compatible type (only useful for compare with sorting semantics)
  VPCMPD $VPCMP_IMM_LE, Z28, Z18, K2, K5               // K5 <- null/bool comparison predicate
  VPTESTNMD.BCST CONSTD_0x80(), Z29, K3, K3            // K3 <- lanes to be cleared from K1 in case they are not comparable and sorting semantics disabled
  VPXORD X27, X27, X27
  VPSUBD Z15, Z14, K5, Z16                             // Z16 <- merged comparison results from NULL/BOOL comparison
  KANDNW K2, K5, K2                                    // K2 <- comparable lanes, without nulls/bools
  KXORW K3, K1, K1                                     // K1 <- updated output predicate to follow sorting semantics
  VPCMPUD.BCST $VPCMP_IMM_LE, CONSTD_4(), Z18, K2, K3  // K3 <- number comparison predicate
  VPGATHERDQ 0(R8)(Y21*1), K4, Z27                     // Z27 <- [R] value 8 content bytes (high)
  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

  KTESTW K3, K3
  JZ compare_bytes                       // skip number comparison if there are no numbers

  // Number Comparison - I64/F64
  // ---------------------------

dispatch_compare_number:
  VPBROADCASTD CONSTD_8(), Z17                         // Z17 <- dword(8)

  // make K2 contain only lanes without number comparison, will be used later to decide whether
  // we are done or whether there are more lanes (of different value types) to be compared
  KANDNW K2, K3, K2

  // let's test K2 here, so the branch predictor sees the flag early
  KTESTW K2, K2

  // byteswap each value and shift right in case of signed/unsigned int
  VPSUBD Z11, Z17, Z22
  VPSUBD Z13, Z17, Z28
  VPSLLD $3, Z22, Z22
  VPSLLD $3, Z28, Z28

  VEXTRACTI32X8 $1, Z22, Y23
  VPMOVZXDQ Y22, Z22
  VEXTRACTI32X8 $1, Z28, Y29
  VPMOVZXDQ Y28, Z28
  VPMOVZXDQ Y23, Z23
  VPMOVZXDQ Y29, Z29

  VPSHUFB Z30, Z24, Z20
  VPSHUFB Z30, Z25, Z21
  VPSRLVQ Z22, Z20, Z20
  VPSHUFB Z30, Z26, Z22
  VPSRLVQ Z23, Z21, Z21
  VPSHUFB Z30, Z27, Z23
  VPSRLD $30, Z31, Z17                                 // Z17 <- dword(3)
  VPSRLVQ Z28, Z22, Z22
  VPSRLVQ Z29, Z23, Z23

  // apply nagation to negative integers, which are stored as positive integers in ION
  VPCMPEQD Z17, Z18, K3, K4
  VPCMPEQD Z17, Z19, K3, K5

  VPXORQ X28, X28, X28
  KSHIFTRW $8, K4, K6
  VPSUBQ Z20, Z28, K4, Z20
  VPSUBQ Z21, Z28, K6, Z21

  KSHIFTRW $8, K5, K6
  VPSUBQ Z22, Z28, K5, Z22
  VPSUBQ Z23, Z28, K6, Z23

  // Now we have either a double precision floating point or int64 (per lane) in Z20|Z21 (left) and Z22|Z23 (right).
  // What we want is to compare floats with floats and integers with integers. Our canonical format is designed in
  // a way that we only use floating point in case that integer is not representable. This means that if a value is
  // floating point, but without a fraction, it's beyond a 64-bit integer range. This leads to a conclusion that if
  // there is an integer vs floating point, we convert the integer to floating point and compare floats.

  VPCMPUD $VPCMP_IMM_LE, Z17, Z18, K3, K4              // K4 <- [L] integer values (positive or negative)
  VPCMPUD $VPCMP_IMM_LE, Z17, Z19, K3, K5              // K5 <- [R] integer values (positive or negative)
  KANDW K4, K5, K6                                     // K6 <- integer values on both sides
  KANDNW K4, K6, K4                                    // K4 <- [L] integer values to convert to floats
  KANDNW K5, K6, K5                                    // K5 <- [R] integer values to convert to floats

  // Convert mixed integer/floating point values on both lanes to floating point
  VCVTQQ2PD Z20, K4, Z20
  KSHIFTRW $8, K4, K4
  VCVTQQ2PD Z22, K5, Z22
  KSHIFTRW $8, K5, K5
  VCVTQQ2PD Z21, K4, Z21
  VCVTQQ2PD Z23, K5, Z23

  KANDNW K3, K6, K5                                    // K5 <- floating point values on both sides
  KSHIFTRW $8, K3, K4                                  // K4 <- number predicate (high)
  KSHIFTRW $8, K5, K6                                  // K6 <- floating point values on both sides (high)

  VPANDQ.Z Z22, Z20, K5, Z28
  VPANDQ.Z Z23, Z21, K6, Z29
  VPMOVQ2M Z28, K5                                     // K5 <- floating point negative values (low)
  VPMOVQ2M Z29, K6                                     // K6 <- floating point negative values (high)

  VPXORD X28, X28, X28
  KUNPCKBW K5, K6, K5                                  // K5 <- floating point negative values (all)

  VPCMPQ $VPCMP_IMM_LT, Z22, Z20, K3, K0               // K0 <- less than (low)
  VPCMPQ $VPCMP_IMM_LT, Z23, Z21, K4, K6               // K6 <- less than (high)
  KUNPCKBW K0, K6, K6                                  // K6 <- less than (all)
  VMOVDQA32 Z31, K6, Z16                               // Z16 <- merge less than results (-1)

  VPCMPQ $VPCMP_IMM_GT, Z22, Z20, K3, K0               // K0 <- greater than (low)
  VPCMPQ $VPCMP_IMM_GT, Z23, Z21, K4, K6               // K6 <- greater than (high)
  KUNPCKBW K0, K6, K6                                  // K6 <- greater than (all)
  VPABSD Z31, K6, Z16                                  // Z16 <- merge greater than results (1)

  VPSUBD Z16, Z28, K5, Z16                             // Z16 <- results with corrected floating point comparison
  JZ next

  // Bytes Comparison - Unsymbolize
  // ------------------------------

compare_bytes:
  // Convert base offsets to absolute 64-bit pointers.
  VEXTRACTI32X8 $1, Z10, Y14
  VEXTRACTI32X8 $1, Z12, Y15
  VPMOVZXDQ Y10, Z10
  VPMOVZXDQ Y14, Z14

  VPBROADCASTQ VIRT_BASE, Z20
  VPBROADCASTQ R8, Z21
  VPMOVZXDQ Y12, Z12
  VPMOVZXDQ Y15, Z15

  VPADDQ Z20, Z10, Z10
  VPADDQ Z20, Z14, Z14
  VPADDQ Z21, Z12, Z12
  VPADDQ Z21, Z15, Z15

  // To continue comparing string and timestamp values, we have to first "unsymbolize".
  VPSRLD $29, Z31, Z17                                 // Z17 <- dword(7)
  VPCMPEQD Z17, Z18, K2, K3                            // K3 <- [L] lanes that are symbols
  VPCMPEQD Z17, Z19, K2, K4                            // K4 <- [R] lanes that are symbols

  KORTESTW K3, K4
  JZ skip_unsymbolize                                  // don't unsymbolize if there are no symbols

  VPMOVQD Z24, Y20
  VPMOVQD Z25, Y21
  VPMOVQD Z26, Y22
  VPMOVQD Z27, Y23
  VINSERTI32X8 $1, Y21, Z20, Z20                       // Z20 <- [L] 4 bytes
  VINSERTI32X8 $1, Y23, Z22, Z21                       // Z21 <- [R] 4 bytes

  VPBROADCASTD CONSTD_4(), Z17                         // Z17 <- dword(4)
  VBROADCASTI32X4 bswap32<>+0(SB), Z22                 // Z22 <- predicate(bswap32)

  MOVQ bytecode_symtab+0(VIRT_BCPTR), BX               // BX <- Symbol table
  VPBROADCASTD bytecode_symtab+8(VIRT_BCPTR), Z23      // Z23 <- Number of symbols in symbol table

  VPSUBD Z11, Z17, Z28
  VPSUBD Z13, Z17, Z29
  VPSHUFB Z22, Z20, Z20
  VPSHUFB Z22, Z21, Z21
  VPSLLD $3, Z28, Z28
  VPSLLD $3, Z29, Z29
  VPSRLVD Z28, Z20, Z28                                // Z28 <- [L] SymbolIDs
  VPSRLVD Z29, Z21, Z29                                // Z29 <- [R] SymbolIDs

  // only unsymbolize lanes where id < len(symtab)
  VPCMPUD $VPCMP_IMM_LT, Z23, Z28, K3, K3              // K3 <- [L] symbols that are in symtab
  KMOVW K3, K5
  VPGATHERDQ 0(BX)(Y28*8), K5, Z20                     // Z20 <- [L] vmrefs of symbols (low)

  KSHIFTRW $8, K3, K6
  VEXTRACTI32X8 $1, Z28, Y28
  VPGATHERDQ 0(BX)(Y28*8), K6, Z21                     // Z21 <- [L] vmrefs of symbols (high)

  VPCMPUD $VPCMP_IMM_LT, Z23, Z29, K4, K4              // K4 <- [R] symbols that are in symtab
  KMOVW K4, K5
  KSHIFTRW $8, K4, K6
  VPGATHERDQ 0(BX)(Y29*8), K5, Z22                     // Z22 <- [R] vmrefs of symbols (low)

  VEXTRACTI32X8 $1, Z29, Y29
  VPGATHERDQ 0(BX)(Y29*8), K6, Z23                     // Z23 <- [R] vmrefs of symbols (high)

  VPBROADCASTD CONSTD_14(), Z17
  KSHIFTRW $8, K3, K6

  BC_MERGE_VMREFS_TO_VALUE_PTR_LEN(IN_OUT(Z10), IN_OUT(Z14), IN_OUT(Z11), IN(Z20), IN(Z21), IN(K3), IN(K6), Z28, Y28, Z29, Y29)
  BC_CALC_VALUE_HLEN_ALT(OUT(Z28), IN(Z11), IN(K3), IN(Z31), IN(Z17), Z29, K5)

  KMOVW K3, K5
  VPSUBD Z28, Z11, K3, Z11
  VEXTRACTI32X8 $1, Z28, Y29
  VPMOVZXDQ Y28, Z28
  VPMOVZXDQ Y29, Z29

  VPADDQ Z28, Z10, K5, Z10
  VPGATHERQQ 0(Z10*1), K5, Z24
  VPADDQ Z29, Z14, K6, Z14
  VPGATHERQQ 0(Z14*1), K6, Z25

  KSHIFTRW $8, K4, K6

  BC_MERGE_VMREFS_TO_VALUE_PTR_LEN(IN_OUT(Z12), IN_OUT(Z15), IN_OUT(Z13), IN(Z22), IN(Z23), IN(K4), IN(K6), Z28, Y28, Z29, Y29)
  BC_CALC_VALUE_HLEN_ALT(OUT(Z28), IN(Z13), IN(K4), IN(Z31), IN(Z17), Z29, K5)

  KMOVW K4, K5
  VPSUBD Z28, Z13, K4, Z13
  VEXTRACTI32X8 $1, Z28, Y29
  VPMOVZXDQ Y28, Z28
  VPMOVZXDQ Y29, Z29

  VPADDQ Z28, Z12, K5, Z12
  VPGATHERQQ 0(Z12*1), K5, Z26
  VPADDQ Z29, Z15, K6, Z15
  VPGATHERQQ 0(Z15*1), K6, Z27

skip_unsymbolize:

  // Bytes Comparison - Prepare
  // --------------------------

  VPSUBD Z13, Z11, K2, Z16                             // Z16 <- merged length comparison
  VPMINUD Z13, Z11, Z28                                // Z28 <- length iterator (min(left, right) length) (decreasing)

  // Bytes Comparison - Vector
  // -------------------------

  // We keep K2 alive - it's not really necessary in the current implementation, but it's
  // likely we would want to extend this to support lists and structs in the future.
  // Additionally - to prevent bugs triggered by empty strings that have arbitrary offsets,
  // but zero lengths, we filter them here. Any string that has zero length would be already
  // compared before entering vector or scalar loop.

  VPTESTMD Z28, Z28, K2, K3                            // K3 <- comparison predicate (values having non-zero length)
  VPBROADCASTD CONSTD_8(), Z17                         // Z17 <- dword(8)

  // Avoid gathering bytes that we have already gathered. The idea is to use the 8
  // bytes of each lane that we have already gathered, and to do some computation
  // that we do meanwhile gathering inside the loop here (as otherwise we would
  // have to avoid doing any computations meanwhile gathering).

  VPMINUD.Z Z17, Z28, K3, Z23                          // Z23 <- number of bytes to compare (max 8).

  VPSUBD Z23, Z17, Z22                                 // Z22 <- number of bytes to discard (8 - Z23).
  VEXTRACTI32X8 $1, Z23, Y21
  VPSUBD Z23, Z28, K3, Z28                             // Z28 <- adjusted length to compare
  VPSLLD $3, Z22, Z22                                  // Z22 <- number of bits to discard

  VPMOVZXDQ Y23, Z20                                   // Z20 <- number of bytes to compare as QWORD (low)
  VPMOVZXDQ Y21, Z21                                   // Z21 <- number of bytes to compare as QWORD (high)
  VPADDQ Z20, Z10, Z10                                 // Z10 <- [L] advance left pointer (low)
  VPADDQ Z21, Z14, Z14                                 // Z14 <- [L] advance left pointer (high)
  VPADDQ Z20, Z12, Z12                                 // Z12 <- [R] advance left pointer (low)
  VPADDQ Z21, Z15, Z15                                 // Z15 <- [R] advance left pointer (high)

  VEXTRACTI32X8 $1, Z22, Y23
  VPMOVZXDQ Y22, Z22
  VPMOVZXDQ Y23, Z23
  JMP compare_bytes_vector_after_gather

  // The idea is to keep using vector loop unless the number of lanes gets too low.
  // The initial eight bytes are always compared in this vector loop to prevent
  // going scalar in case that those eight bytes determine the results of all lanes.

compare_bytes_vector:
  // NOTE: don't clean Z24:Z27 - we are just updating bytes that need to be updated,
  // but we want to keep those that we are not comparing here (as we could need them
  // if we want to compare lists / structs).

  KMOVB K3, K4
  KSHIFTRW $8, K3, K5
  VPGATHERQQ 0(Z10*1), K4, Z24                         // Z24 <- [L] bytes to compare (low)

  KMOVB K3, K4
  VPMINUD Z17, Z28, Z20                                // Z20 <- number of bytes to compare (max 8).
  VPGATHERQQ 0(Z14*1), K5, Z25                         // Z25 <- [L] bytes to compare (high)

  KSHIFTRW $8, K3, K5
  VPSUBD Z20, Z17, Z22                                 // Z22 <- number of bytes to discard (8 - Z21).
  VPSUBD Z20, Z28, K3, Z28                             // Z28 <- adjusted length to compare

  VEXTRACTI32X8 $1, Z20, Y21                           // Z20 <- number of bytes to compare as DWORD (high)
  VPMOVZXDQ Y20, Z20                                   // Z20 <- number of bytes to compare as QWORD (low)
  VPMOVZXDQ Y21, Z21                                   // Z21 <- number of bytes to compare as QWORD (high)

  VPSLLD $3, Z22, Z22                                  // Z22 <- number of bits to discard
  VPGATHERQQ 0(Z12*1), K4, Z26                         // Z26 <- [R] bytes to compare (low)

  VPADDQ Z20, Z10, Z10                                 // Z10 <- [L] advance left pointer (low)
  VPADDQ Z21, Z14, Z14                                 // Z14 <- [L] advance left pointer (high)
  VPADDQ Z20, Z12, Z12                                 // Z12 <- [R] advance left pointer (low)
  VPGATHERQQ 0(Z15*1), K5, Z27                         // Z27 <- [R] bytes to compare (high)

  VEXTRACTI32X8 $1, Z22, Y23
  VPADDQ Z21, Z15, Z15                                 // Z15 <- [R] advance left pointer (high)

  VPMOVZXDQ Y22, Z22
  VPMOVZXDQ Y23, Z23

compare_bytes_vector_after_gather:
  // To compare bytes, we have to byteswap and eliminate bytes we don't compare:
  //
  //   I -> [HH GG FF EE DD CC BB AA] (8 input bytes)
  //   0 -> [AA BB CC DD EE FF GG HH] (0 bytes to discard => 8 bytes compared)
  //   1 -> [00 AA BB CC DD EE FF GG] (1 byte  to discard => 7 bytes compared)
  //   2 -> [00 00 AA BB CC DD EE FF] (2 bytes to discard => 6 bytes compared)
  //   3 -> [00 00 00 AA BB CC DD EE] (3 bytes to discard => 5 bytes compared)
  //   4 -> [00 00 00 00 AA BB CC DD] (4 bytes to discard => 4 bytes compared)
  //   5 -> [00 00 00 00 00 AA BB CC] (5 bytes to discard => 3 bytes compared)
  //   6 -> [00 00 00 00 00 00 AA BB] (6 bytes to discard => 2 bytes compared)
  //   7 -> [00 00 00 00 00 00 00 AA] (7 bytes to discard => 1 byte  compared)
  KSHIFTRW $8, K3, K4

  VPSLLVQ Z22, Z24, Z20                                // Z20 <- [L] bytes to compare (low)
  VPSLLVQ Z23, Z25, Z21                                // Z21 <- [L] bytes to compare (high)
  VPSHUFB Z30, Z20, Z20                                // Z20 <- [L] byteswapped quadword to compare (low)
  VPSHUFB Z30, Z21, Z21                                // Z21 <- [L] byteswapped quadword to compare (high)

  VPSLLVQ Z22, Z26, Z22                                // Z22 <- [R] bytes to compare (low)
  VPSLLVQ Z23, Z27, Z23                                // Z23 <- [R] bytes to compare (high)
  VPSHUFB Z30, Z22, Z22                                // Z22 <- [R] byteswapped quadword to compare (low)
  VPSHUFB Z30, Z23, Z23                                // Z23 <- [R] byteswapped quadword to compare (high)

  VPCMPQ $VPCMP_IMM_NE, Z22, Z20, K3, K5               // K5 <- lanes having values that aren't equal (low)
  VPCMPQ $VPCMP_IMM_NE, Z23, Z21, K4, K4               // K4 <- lanes having values that aren't equal (high)
  KUNPCKBW K5, K4, K5                                  // K5 <- lanes having values that aren't equal (all lanes)
  KANDNW K3, K5, K3                                    // K3 <- lanes to continue being compared

  VPCMPUQ $VPCMP_IMM_LT, Z22, Z20, K5, K0              // K0 <- lanes where the comparison is less than (low)
  VPCMPUQ $VPCMP_IMM_LT, Z23, Z21, K4, K6              // K6 <- lanes where the comparison is less than (high)
  VPTESTMD Z28, Z28, K3, K3                            // K3 <- lanes to continue being compared (where length is non-zero)
  KUNPCKBW K0, K6, K6                                  // K6 <- lanes where the comparison is less than (all lanes)
  VMOVDQA32 Z31, K6, Z16                               // Z16 <- merge less than results (-1)
  KMOVW K3, BX

  VPCMPUQ $VPCMP_IMM_GT, Z22, Z20, K5, K5              // K5 <- lanes where the comparison is greater than (low)
  VPCMPUQ $VPCMP_IMM_GT, Z23, Z21, K4, K4              // K6 <- lanes where the comparison is greater than (high)
  POPCNTL BX, DX                                       // DX <- number of remaining lanes to process
  KUNPCKBW K5, K4, K5                                  // K5 <- lanes where the comparison is greater than (all lanes)
  VPABSD Z31, K5, Z16                                  // Z16 <- merge greater than results (1)

  TESTL BX, BX
  JZ next

  // Go to scalar loop if the number of lanes to compare gets low
  CMPL DX, $BC_SCALAR_PROCESSING_LANE_COUNT
  JHI compare_bytes_vector

  VMOVDQU32 Z10, BC_SPILL_AREA(0)                      // [L] content pointer (QWORD) (low)
  VMOVDQU32 Z14, BC_SPILL_AREA(64)                     // [L] content pointer (QWORD) (high)
  VMOVDQU32 Z12, BC_SPILL_AREA(128)                    // [R] content pointer (QWORD) (low)
  VMOVDQU32 Z15, BC_SPILL_AREA(192)                    // [R] content pointer (QWORD) (high)
  VMOVDQU32 Z28, BC_SPILL_AREA(256)                    // min(left, right) length (DWORD)
  VMOVDQU32 Z16, BC_SPILL_AREA(320)                    // comparison results (DWORD)

  MOVQ $-1, R13
  JMP compare_bytes_scalar_lane

  // Bytes Comparison (Scalar)
  // -------------------------

compare_bytes_scalar_diff:
  KMOVQ K4, CX
  TZCNTQ CX, CX
  MOVBLZX 0(R14)(CX*1), R14
  MOVBLZX 0(R15)(CX*1), R15
  SUBL R15, R14
  MOVL R14, BC_SPILL_AREA_INDEX(320, DX*4)

  TESTL BX, BX
  JE compare_bytes_scalar_done

compare_bytes_scalar_lane:
  TZCNTL BX, DX                                        // DX - Index of the lane to process
  BLSRL BX, BX                                         // clear the index of the iterator

  MOVL BC_SPILL_AREA_INDEX(256, DX*4), CX              // min(left, right) length
  MOVQ BC_SPILL_AREA_INDEX(0, DX*8), R14               // [L] slice pointer (absolute address)
  MOVQ BC_SPILL_AREA_INDEX(128, DX*8), R15             // [R] slice pointer (absolute address)

  SUBL $64, CX
  JCS compare_bytes_tail

compare_bytes_scalar_iter:                             // main compare loop that processes 64 bytes at once
  VMOVDQU8 0(R14), Z20
  VMOVDQU8 0(R15), Z21
  VPCMPB $VPCMP_IMM_NE, Z21, Z20, K4
  KTESTQ K4, K4
  JNE compare_bytes_scalar_diff

  ADDQ $64, R14
  ADDQ $64, R15
  SUBL $64, CX
  JA compare_bytes_scalar_iter

compare_bytes_tail:                                    // tail loop that processes up to 64 bytes at once
  SHLXQ CX, R13, CX
  NOTQ CX
  KMOVQ CX, K4                                         // K4 <- LSB mask of bits to process (valid characters)

  VMOVDQU8.Z 0(R14), K4, Z20
  VMOVDQU8.Z 0(R15), K4, Z21

  VPCMPB $VPCMP_IMM_NE, Z21, Z20, K4
  KTESTQ K4, K4
  JNE compare_bytes_scalar_diff

  // Comparable slices have the same content, which means that `leftLen-rightLen` is the result
  // This result was already precalculated before entering the scalar loop, so we don't have to
  // calculate and store it again.

  TESTL BX, BX
  JNE compare_bytes_scalar_lane

compare_bytes_scalar_done:
  VMOVDQU32 BC_SPILL_AREA(320), Z16

next:
  VPABSD Z31, Z17                                      // Z17 <- dword(1)
  VPMAXSD Z31, Z16, Z16
  VPMINSD.Z Z17, Z16, K1, Z16

  RET
