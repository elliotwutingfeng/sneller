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

package vm

//go:generate go run -tags genrewrite genrewrite_main.go -o simplify1.go simplify.rules
//go:generate gofmt -w simplify1.go

func dependsOn(a, b *value) bool {
	for _, arg := range a.args {
		if arg == b || dependsOn(arg, b) {
			return true
		}
	}
	return false
}

// given a, b, produce (a AND b, true) or (nil, false)
// by setting the mask of one of the ops to the result
// of the other op
func conjoin(p *prog, a, b *value) (*value, bool) {
	if a == b {
		return nil, false
	}
	try := func(p *prog, a, b *value) (*value, bool) {
		// a must be conjunctive and not be a dependency of b
		if ssainfo[a.op].disjunctive || dependsOn(b, a) {
			return nil, false
		}
		// a must have a mask arg that is the same as b's (and be non-nil)
		mask := a.maskarg()
		if mask == nil || mask != b.maskarg() {
			return nil, false
		}
		// success: produce a with b as its mask
		return p.dup(a).setmask(b), true
	}
	// try both orderings:
	v, ok := try(p, a, b)
	if !ok {
		return try(p, b, a)
	}
	return v, true
}

func isfalse(v *value) (*value, bool) {
	depth := 3 // limit polynomial complexity
	for v != nil && depth > 0 {
		if v.op == skfalse {
			return v, true
		}
		if ssainfo[v.op].disjunctive {
			break
		}
		v = v.maskarg()
		depth--
	}
	return nil, false
}

// eliminateMovs replaces `mov` arguments of `v` with `mov` sources. The purpose of
// this optimization is to completely remove instructions such as `smakevk`, which
// are only used by SSA to create a value with a custom mask associated with it.
func eliminateMovs(p *prog, v *value) bool {
	info := &ssainfo[v.op]
	changed := false

	for i, arg := range v.args {
		if arg.op == smakevk {
			if (info.argtypes[i] & stValue) != 0 {
				v.args[i] = arg.args[0]
				changed = true
			} else if (info.argtypes[i] & stBool) != 0 {
				v.args[i] = arg.args[1]
				changed = true
			}
		}
	}

	return changed
}

func rewrite(p *prog, v *value) (*value, bool) {
	info := &ssainfo[v.op]

	// (op ... false) -> false
	// when the op is conjunctive and returns
	// either stValue or stBool or both
	// (since the type of false is stValueMasked)
	m := v.maskarg()
	if vfalse, ok := isfalse(m); ok {
		if !info.disjunctive && info.rettype&^stValueMasked == 0 {
			return vfalse, true
		}
		// set mask directly to false
		// if we detect a false in the arg chain
		if m.op != skfalse {
			v.setmask(vfalse)
		}
	}

	// apply rewriting until we
	// reach a fixed point
	var opt bool
	v, opt = rewrite1(p, v)
	any := false
	for opt {
		any = true
		v, opt = rewrite1(p, v)
	}
	return v, any
}

// simplify iteratively simplifies
// the program p until it reaches
// a fixed point
//
// see simplify.rules
func (p *prog) simplify(pi *proginfo) {
	var rewrote []*value
	for {
		changed := false
		values := p.values
		rewrote = shrink(rewrote, len(values))
		// reverse-postorder guarantees that
		// optimizations are applied bottom-up,
		// which ought to minimize the number of
		// passes we need to reach a fixed point here
		ord := p.order(pi)
		for _, v := range ord {
			for i, arg := range v.args {
				if rewrote[arg.id] != nil {
					v.args[i] = rewrote[arg.id]
				}
			}

			if eliminateMovs(p, v) {
				changed = true
			}

			out, ok := rewrite(p, v)
			if ok {
				changed = true
				if out != v {
					rewrote[v.id] = out
				}
			}
		}
		if !changed {
			return
		}
		if rewrote[p.ret.id] != nil {
			p.ret = rewrote[p.ret.id]
		}
		pi.invalidate()
	}
}
