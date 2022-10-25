package expr

// code generated by terms.go; DO NOT EDIT
import (
	"math/big"
	"strings"
	"unicode/utf8"
)

func simplifyClass0(src *Arithmetic, h Hint) Node {
	switch src.Op {
	case AddOp:
		// (add x (int "0")), "TypeOf(x, h) == (NumericType|MissingType)" -> x
		if x := src.Left; true {
			if _tmp001001, ok := (src.Right).(Integer); ok {
				if Integer(0).Equals(_tmp001001) {
					if TypeOf(x, h) == (NumericType | MissingType) {
						return x
					}
				}
			}
		}
		// (add (constant x) y), "_, ok := y.(Constant); !ok" -> (add y x)
		if x, ok := (src.Left).(Constant); ok {
			if y := src.Right; true {
				if _, ok := y.(Constant); !ok {
					return &Arithmetic{Op: AddOp, Left: y, Right: x}
				}
			}
		}
		// (add (add x (constant y)) (constant z)) -> (add x (add y z))
		if _tmp001000, ok := (src.Left).(*Arithmetic); ok && _tmp001000.Op == AddOp {
			if z, ok := (src.Right).(Constant); ok {
				if x := _tmp001000.Left; true {
					if y, ok := (_tmp001000.Right).(Constant); ok {
						return &Arithmetic{Op: AddOp, Left: x, Right: &Arithmetic{Op: AddOp, Left: y, Right: z}}
					}
				}
			}
		}
		// (add (add a (constant b)) (add c (constant d))) -> (add (add a c) (add b d))
		if _tmp001000, ok := (src.Left).(*Arithmetic); ok && _tmp001000.Op == AddOp {
			if _tmp001001, ok := (src.Right).(*Arithmetic); ok && _tmp001001.Op == AddOp {
				if a := _tmp001000.Left; true {
					if b, ok := (_tmp001000.Right).(Constant); ok {
						if c := _tmp001001.Left; true {
							if d, ok := (_tmp001001.Right).(Constant); ok {
								return &Arithmetic{Op: AddOp, Left: &Arithmetic{Op: AddOp, Left: a, Right: c}, Right: &Arithmetic{Op: AddOp, Left: b, Right: d}}
							}
						}
					}
				}
			}
		}
	case DivOp:
		// (div x (int "1")), "TypeOf(x, h) == (NumericType|MissingType)" -> x
		if x := src.Left; true {
			if _tmp001001, ok := (src.Right).(Integer); ok {
				if Integer(1).Equals(_tmp001001) {
					if TypeOf(x, h) == (NumericType | MissingType) {
						return x
					}
				}
			}
		}
		// (div _ (int "0")) -> (missing)
		if _tmp001001, ok := (src.Right).(Integer); ok {
			if Integer(0).Equals(_tmp001001) {
				return Missing{}
			}
		}
	case ModOp:
		// (mod _ (int "0")) -> (missing)
		if _tmp001001, ok := (src.Right).(Integer); ok {
			if Integer(0).Equals(_tmp001001) {
				return Missing{}
			}
		}
	case MulOp:
		// (mul x (int "1")), "TypeOf(x, h) == (NumericType|MissingType)" -> x
		if x := src.Left; true {
			if _tmp001001, ok := (src.Right).(Integer); ok {
				if Integer(1).Equals(_tmp001001) {
					if TypeOf(x, h) == (NumericType | MissingType) {
						return x
					}
				}
			}
		}
		// (mul (constant x) y), "_, ok := y.(Constant); !ok" -> (mul y x)
		if x, ok := (src.Left).(Constant); ok {
			if y := src.Right; true {
				if _, ok := y.(Constant); !ok {
					return &Arithmetic{Op: MulOp, Left: y, Right: x}
				}
			}
		}
		// (mul (mul x (constant y)) (constant z)) -> (mul x (mul y z))
		if _tmp001000, ok := (src.Left).(*Arithmetic); ok && _tmp001000.Op == MulOp {
			if z, ok := (src.Right).(Constant); ok {
				if x := _tmp001000.Left; true {
					if y, ok := (_tmp001000.Right).(Constant); ok {
						return &Arithmetic{Op: MulOp, Left: x, Right: &Arithmetic{Op: MulOp, Left: y, Right: z}}
					}
				}
			}
		}
		// (mul (mul a (constant b)) (mul c (constant d))) -> (mul (mul a c) (mul b d))
		if _tmp001000, ok := (src.Left).(*Arithmetic); ok && _tmp001000.Op == MulOp {
			if _tmp001001, ok := (src.Right).(*Arithmetic); ok && _tmp001001.Op == MulOp {
				if a := _tmp001000.Left; true {
					if b, ok := (_tmp001000.Right).(Constant); ok {
						if c := _tmp001001.Left; true {
							if d, ok := (_tmp001001.Right).(Constant); ok {
								return &Arithmetic{Op: MulOp, Left: &Arithmetic{Op: MulOp, Left: a, Right: c}, Right: &Arithmetic{Op: MulOp, Left: b, Right: d}}
							}
						}
					}
				}
			}
		}
	case SubOp:
		// (sub x (int "0")), "TypeOf(x, h) == (NumericType|MissingType)" -> x
		if x := src.Left; true {
			if _tmp001001, ok := (src.Right).(Integer); ok {
				if Integer(0).Equals(_tmp001001) {
					if TypeOf(x, h) == (NumericType | MissingType) {
						return x
					}
				}
			}
		}
	}
	return nil
}

func simplifyClass1(src *Builtin, h Hint) Node {
	switch src.Func {
	case Abs:
		if len(src.Args) == 1 {
			// (abs (number x)) -> "(*Rational)(new(big.Rat).Abs(x.rat()))"
			if x, ok := (src.Args[0]).(number); ok {
				return (*Rational)(new(big.Rat).Abs(x.rat()))
			}
		}
	case CharLength:
		if len(src.Args) == 1 {
			// (char_length (concat x y)) -> (add (char_length x) (char_length y))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Concat && len(_tmp001000.Args) == 2 {
				if x := _tmp001000.Args[0]; true {
					if y := _tmp001000.Args[1]; true {
						return &Arithmetic{Op: AddOp, Left: Call(CharLength, x), Right: Call(CharLength, y)}
					}
				}
			}
			// (char_length (string x)) -> (int "utf8.RuneCountInString(string(x))")
			if x, ok := (src.Args[0]).(String); ok {
				return Integer(utf8.RuneCountInString(string(x)))
			}
			// (char_length (lower x)) -> (char_length x)
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(CharLength, x)
				}
			}
			// (char_length (upper x)) -> (char_length x)
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(CharLength, x)
				}
			}
		}
	case Concat:
		if len(src.Args) == 2 {
			// (concat (string x) (string y)) -> (string "x + y")
			if x, ok := (src.Args[0]).(String); ok {
				if y, ok := (src.Args[1]).(String); ok {
					return String(x + y)
				}
			}
			// (concat (string x) (concat (string y) z)) -> (concat (string "x + y") z)
			if x, ok := (src.Args[0]).(String); ok {
				if _tmp001001, ok := (src.Args[1]).(*Builtin); ok && _tmp001001.Func == Concat && len(_tmp001001.Args) == 2 {
					if y, ok := (_tmp001001.Args[0]).(String); ok {
						if z := _tmp001001.Args[1]; true {
							return Call(Concat, String(x+y), z)
						}
					}
				}
			}
			// (concat (concat x (string a)) (string b)) -> (concat x (string "a + b"))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Concat && len(_tmp001000.Args) == 2 {
				if b, ok := (src.Args[1]).(String); ok {
					if x := _tmp001000.Args[0]; true {
						if a, ok := (_tmp001000.Args[1]).(String); ok {
							return Call(Concat, x, String(a+b))
						}
					}
				}
			}
			// (concat x (string "\"\"")) -> (assert_str x)
			if x := src.Args[0]; true {
				if _tmp001001, ok := (src.Args[1]).(String); ok {
					if String("").Equals(_tmp001001) {
						return Call(AssertIonType, x, Integer(0x8))
					}
				}
			}
			// (concat (upper x) (string y)), "isUpper(string(y))" -> (upper (concat x y))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if y, ok := (src.Args[1]).(String); ok {
					if x := _tmp001000.Args[0]; true {
						if isUpper(string(y)) {
							return Call(Upper, Call(Concat, x, y))
						}
					}
				}
			}
			// (concat (upper x) (upper y)) -> (upper (concat x y))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if _tmp001001, ok := (src.Args[1]).(*Builtin); ok && _tmp001001.Func == Upper && len(_tmp001001.Args) == 1 {
					if x := _tmp001000.Args[0]; true {
						if y := _tmp001001.Args[0]; true {
							return Call(Upper, Call(Concat, x, y))
						}
					}
				}
			}
			// (concat (lower x) (string y)), "isLower(string(y))" -> (lower (concat x y))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if y, ok := (src.Args[1]).(String); ok {
					if x := _tmp001000.Args[0]; true {
						if isLower(string(y)) {
							return Call(Lower, Call(Concat, x, y))
						}
					}
				}
			}
			// (concat (lower x) (lower y)) -> (lower (concat x y))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if _tmp001001, ok := (src.Args[1]).(*Builtin); ok && _tmp001001.Func == Lower && len(_tmp001001.Args) == 1 {
					if x := _tmp001000.Args[0]; true {
						if y := _tmp001001.Args[0]; true {
							return Call(Lower, Call(Concat, x, y))
						}
					}
				}
			}
			// (concat (string x) (upper y)), "isUpper(string(x))" -> (upper (concat x y))
			if x, ok := (src.Args[0]).(String); ok {
				if _tmp001001, ok := (src.Args[1]).(*Builtin); ok && _tmp001001.Func == Upper && len(_tmp001001.Args) == 1 {
					if y := _tmp001001.Args[0]; true {
						if isUpper(string(x)) {
							return Call(Upper, Call(Concat, x, y))
						}
					}
				}
			}
			// (concat (string x) (lower y)), "isLower(string(x))" -> (lower (concat x y))
			if x, ok := (src.Args[0]).(String); ok {
				if _tmp001001, ok := (src.Args[1]).(*Builtin); ok && _tmp001001.Func == Lower && len(_tmp001001.Args) == 1 {
					if y := _tmp001001.Args[0]; true {
						if isLower(string(x)) {
							return Call(Lower, Call(Concat, x, y))
						}
					}
				}
			}
		}
	case Contains:
		if len(src.Args) == 2 {
			// (contains (string x) (string y)) -> "Bool(strings.Contains(string(x), string(y)))"
			if x, ok := (src.Args[0]).(String); ok {
				if y, ok := (src.Args[1]).(String); ok {
					return Bool(strings.Contains(string(x), string(y)))
				}
			}
			// (contains (upper x) (string y)), "isUpper(string(y))" -> (contains_ci x y)
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if y, ok := (src.Args[1]).(String); ok {
					if x := _tmp001000.Args[0]; true {
						if isUpper(string(y)) {
							return Call(ContainsCI, x, y)
						}
					}
				}
			}
			// (contains (upper _) (string y)), "!isUpper(string(y))" -> (bool "false")
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if y, ok := (src.Args[1]).(String); ok {
					if !isUpper(string(y)) {
						return Bool(false)
					}
				}
			}
			// (contains (lower x) (string y)), "isLower(string(y))" -> (contains_ci x y)
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if y, ok := (src.Args[1]).(String); ok {
					if x := _tmp001000.Args[0]; true {
						if isLower(string(y)) {
							return Call(ContainsCI, x, y)
						}
					}
				}
			}
			// (contains (lower _) (string y)), "!isLower(string(y))" -> (bool "false")
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if y, ok := (src.Args[1]).(String); ok {
					if !isLower(string(y)) {
						return Bool(false)
					}
				}
			}
		}
	case DateExtractDay:
		if len(src.Args) == 1 {
			// (date_extract_day (ts x)) -> (int "x.Value.Day()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Day())
			}
		}
	case DateExtractHour:
		if len(src.Args) == 1 {
			// (date_extract_hour (ts x)) -> (int "x.Value.Hour()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Hour())
			}
		}
	case DateExtractMicrosecond:
		if len(src.Args) == 1 {
			// (date_extract_microsecond (ts x)) -> (int "x.Value.Nanosecond() / 1000")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Nanosecond() / 1000)
			}
		}
	case DateExtractMillisecond:
		if len(src.Args) == 1 {
			// (date_extract_millisecond (ts x)) -> (int "x.Value.Nanosecond() / 1000000")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Nanosecond() / 1000000)
			}
		}
	case DateExtractMinute:
		if len(src.Args) == 1 {
			// (date_extract_minute (ts x)) -> (int "x.Value.Minute()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Minute())
			}
		}
	case DateExtractMonth:
		if len(src.Args) == 1 {
			// (date_extract_month (ts x)) -> (int "x.Value.Month()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Month())
			}
		}
	case DateExtractQuarter:
		if len(src.Args) == 1 {
			// (date_extract_quarter (ts x)) -> (int "x.Value.Quarter()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Quarter())
			}
		}
	case DateExtractSecond:
		if len(src.Args) == 1 {
			// (date_extract_second (ts x)) -> (int "x.Value.Second()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Second())
			}
		}
	case DateExtractYear:
		if len(src.Args) == 1 {
			// (date_extract_year (ts x)) -> (int "x.Value.Year()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Year())
			}
		}
	case Lower:
		if len(src.Args) == 1 {
			// (lower (string x)) -> (string "strings.ToLower(string(x))")
			if x, ok := (src.Args[0]).(String); ok {
				return String(strings.ToLower(string(x)))
			}
		}
	case Ltrim:
		if len(src.Args) == 1 {
			// (ltrim (rtrim x)) -> (trim x)
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Rtrim && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Trim, x)
				}
			}
			// (ltrim inner:(trim _)) -> inner
			if inner, ok := (src.Args[0]).(*Builtin); ok && inner.Func == Trim && len(inner.Args) == 1 {
				return inner
			}
			// (ltrim inner:(ltrim _)) -> inner
			if inner, ok := (src.Args[0]).(*Builtin); ok && inner.Func == Ltrim && len(inner.Args) == 1 {
				return inner
			}
			// (ltrim (upper x)) -> (upper (ltrim x))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Upper, Call(Ltrim, x))
				}
			}
			// (ltrim (upper x)) -> (upper (ltrim x))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Upper, Call(Ltrim, x))
				}
			}
		}
	case ObjectSize:
		if len(src.Args) == 1 {
			// (object_size (list l)) -> "Integer(len(l.Values))"
			if l, ok := (src.Args[0]).(*List); ok {
				return Integer(len(l.Values))
			}
			// (object_size (struct s)) -> "Integer(len(s.Fields))"
			if s, ok := (src.Args[0]).(*Struct); ok {
				return Integer(len(s.Fields))
			}
			// (object_size (missing)) -> (missing)
			if _, ok := (src.Args[0]).(Missing); ok {
				return Missing{}
			}
			// (object_size (null)) -> (null)
			if _, ok := (src.Args[0]).(Null); ok {
				return Null{}
			}
		}
	case Rtrim:
		if len(src.Args) == 1 {
			// (rtrim (ltrim x)) -> (trim x)
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Ltrim && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Trim, x)
				}
			}
			// (rtrim inner:(rtrim _)) -> inner
			if inner, ok := (src.Args[0]).(*Builtin); ok && inner.Func == Rtrim && len(inner.Args) == 1 {
				return inner
			}
			// (rtrim inner:(trim _)) -> inner
			if inner, ok := (src.Args[0]).(*Builtin); ok && inner.Func == Trim && len(inner.Args) == 1 {
				return inner
			}
			// (rtrim (lower x)) -> (lower (rtrim x))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Lower, Call(Rtrim, x))
				}
			}
			// (rtrim (lower x)) -> (lower (rtrim x))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Lower, Call(Rtrim, x))
				}
			}
		}
	case Sign:
		if len(src.Args) == 1 {
			// (sign (number x)) -> "(*Rational)(new(big.Rat).SetInt64(int64(x.rat().Sign())))"
			if x, ok := (src.Args[0]).(number); ok {
				return (*Rational)(new(big.Rat).SetInt64(int64(x.rat().Sign())))
			}
		}
	case Substring:
		if len(src.Args) == 2 {
			// (substring s (int "1")), "TypeOf(s, h) == StringType|MissingType" -> s
			if s := src.Args[0]; true {
				if _tmp001001, ok := (src.Args[1]).(Integer); ok {
					if Integer(1).Equals(_tmp001001) {
						if TypeOf(s, h) == StringType|MissingType {
							return s
						}
					}
				}
			}
			// (substring (string s) (int start)) -> "staticSubstr(s, start, 1<<21)"
			if s, ok := (src.Args[0]).(String); ok {
				if start, ok := (src.Args[1]).(Integer); ok {
					return staticSubstr(s, start, 1<<21)
				}
			}
			// (substring s x) -> (substring s x (int "1<<21"))
			if s := src.Args[0]; true {
				if x := src.Args[1]; true {
					return Call(Substring, s, x, Integer(1<<21))
				}
			}
			// (substring (lower x) off) -> (lower (substring x off))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if off := src.Args[1]; true {
					if x := _tmp001000.Args[0]; true {
						return Call(Lower, Call(Substring, x, off))
					}
				}
			}
			// (substring (upper x) off) -> (upper (substring x off))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if off := src.Args[1]; true {
					if x := _tmp001000.Args[0]; true {
						return Call(Upper, Call(Substring, x, off))
					}
				}
			}
		}
		if len(src.Args) == 3 {
			// (substring (string s) (int start) (int len)) -> "staticSubstr(s, start, len)"
			if s, ok := (src.Args[0]).(String); ok {
				if start, ok := (src.Args[1]).(Integer); ok {
					if len, ok := (src.Args[2]).(Integer); ok {
						return staticSubstr(s, start, len)
					}
				}
			}
			// (substring (lower x) off len) -> (lower (substring x off len))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if off := src.Args[1]; true {
					if len := src.Args[2]; true {
						if x := _tmp001000.Args[0]; true {
							return Call(Lower, Call(Substring, x, off, len))
						}
					}
				}
			}
			// (substring (upper x) off len) -> (upper (substring x off len))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if off := src.Args[1]; true {
					if len := src.Args[2]; true {
						if x := _tmp001000.Args[0]; true {
							return Call(Upper, Call(Substring, x, off, len))
						}
					}
				}
			}
		}
	case ToUnixEpoch:
		if len(src.Args) == 1 {
			// (to_unix_epoch (ts x)) -> (int "x.Value.Unix()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.Unix())
			}
		}
	case ToUnixMicro:
		if len(src.Args) == 1 {
			// (to_unix_micro (ts x)) -> (int "x.Value.UnixMicro()")
			if x, ok := (src.Args[0]).(*Timestamp); ok {
				return Integer(x.Value.UnixMicro())
			}
		}
	case Trim:
		if len(src.Args) == 1 {
			// (trim inner:(trim _)) -> inner
			if inner, ok := (src.Args[0]).(*Builtin); ok && inner.Func == Trim && len(inner.Args) == 1 {
				return inner
			}
			// (trim (lower x)) -> (lower (trim x))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Lower, Call(Trim, x))
				}
			}
			// (trim (upper x)) -> (upper (trim x))
			if _tmp001000, ok := (src.Args[0]).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
				if x := _tmp001000.Args[0]; true {
					return Call(Upper, Call(Trim, x))
				}
			}
		}
	case Upper:
		if len(src.Args) == 1 {
			// (upper (string x)) -> (string "strings.ToUpper(string(x))")
			if x, ok := (src.Args[0]).(String); ok {
				return String(strings.ToUpper(string(x)))
			}
		}
	}
	return nil
}

func simplifyClass2(src *Comparison, h Hint) Node {
	switch src.Op {
	case Equals:
		// (eq x x), "TypeOf(x, h)&MissingType == 0" -> (bool "true")
		if x := src.Left; true {
			if x.Equals(src.Right) {
				if TypeOf(x, h)&MissingType == 0 {
					return Bool(true)
				}
			}
		}
		// (eq (upper _) (string lit)), "!isUpper(string(lit))" -> (bool "false")
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
			if lit, ok := (src.Right).(String); ok {
				if !isUpper(string(lit)) {
					return Bool(false)
				}
			}
		}
		// (eq (lower _) (string lit)), "!isLower(string(lit))" -> (bool "false")
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
			if lit, ok := (src.Right).(String); ok {
				if !isLower(string(lit)) {
					return Bool(false)
				}
			}
		}
		// (eq (upper x) (string lit)), "isUpper(string(lit))" -> (equals_ci x lit)
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
			if lit, ok := (src.Right).(String); ok {
				if x := _tmp001000.Args[0]; true {
					if isUpper(string(lit)) {
						return Call(EqualsCI, x, lit)
					}
				}
			}
		}
		// (eq (lower x) (string lit)), "isLower(string(lit))" -> (equals_ci x lit)
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
			if lit, ok := (src.Right).(String); ok {
				if x := _tmp001000.Args[0]; true {
					if isLower(string(lit)) {
						return Call(EqualsCI, x, lit)
					}
				}
			}
		}
	case Greater:
		// (gt (ts x) (ts y)) -> (bool "y.Value.Before(x.Value)")
		if x, ok := (src.Left).(*Timestamp); ok {
			if y, ok := (src.Right).(*Timestamp); ok {
				return Bool(y.Value.Before(x.Value))
			}
		}
	case GreaterEquals:
		// (gte x x), "TypeOf(x, h)&MissingType == 0" -> (bool "true")
		if x := src.Left; true {
			if x.Equals(src.Right) {
				if TypeOf(x, h)&MissingType == 0 {
					return Bool(true)
				}
			}
		}
		// (gte (ts x) (ts y)) -> (bool "y.Value.Before(x.Value) || x.Value == y.Value")
		if x, ok := (src.Left).(*Timestamp); ok {
			if y, ok := (src.Right).(*Timestamp); ok {
				return Bool(y.Value.Before(x.Value) || x.Value == y.Value)
			}
		}
	case Ilike:
		// (ilike x (string pat)), "!strings.ContainsAny(string(pat), \"%_\")" -> (equals_ci x pat)
		if x := src.Left; true {
			if pat, ok := (src.Right).(String); ok {
				if !strings.ContainsAny(string(pat), "%_") {
					return Call(EqualsCI, x, pat)
				}
			}
		}
		// (ilike x (string pat)), "term, ok := isSubstringSearchPattern(string(pat)); ok" -> (contains_ci x "String(term)")
		if x := src.Left; true {
			if pat, ok := (src.Right).(String); ok {
				if term, ok := isSubstringSearchPattern(string(pat)); ok {
					return Call(ContainsCI, x, String(term))
				}
			}
		}
	case Like:
		// (like x (string pat)), "!strings.ContainsAny(string(pat), \"%_\")" -> (eq x pat)
		if x := src.Left; true {
			if pat, ok := (src.Right).(String); ok {
				if !strings.ContainsAny(string(pat), "%_") {
					return &Comparison{Op: Equals, Left: x, Right: pat}
				}
			}
		}
		// (like x (string pat)), "term, ok := isSubstringSearchPattern(string(pat)); ok" -> (contains x "String(term)")
		if x := src.Left; true {
			if pat, ok := (src.Right).(String); ok {
				if term, ok := isSubstringSearchPattern(string(pat)); ok {
					return Call(Contains, x, String(term))
				}
			}
		}
		// (like (upper x) (string pat)), "isUpper(string(pat))" -> (ilike x pat)
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
			if pat, ok := (src.Right).(String); ok {
				if x := _tmp001000.Args[0]; true {
					if isUpper(string(pat)) {
						return &Comparison{Op: Ilike, Left: x, Right: pat}
					}
				}
			}
		}
		// (like (lower x) (string pat)), "isLower(string(pat))" -> (ilike x pat)
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
			if pat, ok := (src.Right).(String); ok {
				if x := _tmp001000.Args[0]; true {
					if isLower(string(pat)) {
						return &Comparison{Op: Ilike, Left: x, Right: pat}
					}
				}
			}
		}
		// (like (upper _) (string pat)), "!isUpper(string(pat))" -> (bool "false")
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
			if pat, ok := (src.Right).(String); ok {
				if !isUpper(string(pat)) {
					return Bool(false)
				}
			}
		}
		// (like (lower _) (string pat)), "!isLower(string(pat))" -> (bool "false")
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
			if pat, ok := (src.Right).(String); ok {
				if !isLower(string(pat)) {
					return Bool(false)
				}
			}
		}
	case Less:
		// (lt (ts x) (ts y)) -> (bool "x.Value.Before(y.Value)")
		if x, ok := (src.Left).(*Timestamp); ok {
			if y, ok := (src.Right).(*Timestamp); ok {
				return Bool(x.Value.Before(y.Value))
			}
		}
	case LessEquals:
		// (lte x x), "TypeOf(x, h)&MissingType == 0" -> (bool "true")
		if x := src.Left; true {
			if x.Equals(src.Right) {
				if TypeOf(x, h)&MissingType == 0 {
					return Bool(true)
				}
			}
		}
		// (lte (ts x) (ts y)) -> (bool "x.Value.Before(y.Value) || x.Value == y.Value")
		if x, ok := (src.Left).(*Timestamp); ok {
			if y, ok := (src.Right).(*Timestamp); ok {
				return Bool(x.Value.Before(y.Value) || x.Value == y.Value)
			}
		}
	case NotEquals:
		// (neq x x), "TypeOf(x, h)&MissingType == 0" -> (bool "true")
		if x := src.Left; true {
			if x.Equals(src.Right) {
				if TypeOf(x, h)&MissingType == 0 {
					return Bool(true)
				}
			}
		}
		// (neq (upper x) (string lit)), "isUpper(string(lit))" -> "&Not{Call(EqualsCI, x, lit)}"
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Upper && len(_tmp001000.Args) == 1 {
			if lit, ok := (src.Right).(String); ok {
				if x := _tmp001000.Args[0]; true {
					if isUpper(string(lit)) {
						return &Not{Call(EqualsCI, x, lit)}
					}
				}
			}
		}
		// (neq (lower x) (string lit)), "isLower(string(lit))" -> "&Not{Call(EqualsCI, x, lit)}"
		if _tmp001000, ok := (src.Left).(*Builtin); ok && _tmp001000.Func == Lower && len(_tmp001000.Args) == 1 {
			if lit, ok := (src.Right).(String); ok {
				if x := _tmp001000.Args[0]; true {
					if isLower(string(lit)) {
						return &Not{Call(EqualsCI, x, lit)}
					}
				}
			}
		}
	}
	return nil
}

func simplifyClass3(src *IsKey, h Hint) Node {
	switch src.Key {
	case IsFalse:
		// (is_false (bool x)) -> (bool "!x")
		if x, ok := (src.Expr).(Bool); ok {
			return Bool(!x)
		}
		// (is_false x), "TypeOf(x, h)&BoolType == 0" -> (bool "false")
		if x := src.Expr; true {
			if TypeOf(x, h)&BoolType == 0 {
				return Bool(false)
			}
		}
	case IsMissing:
		// (is_missing (missing)) -> (bool "true")
		if _, ok := (src.Expr).(Missing); ok {
			return Bool(true)
		}
		// (is_missing (constant _)) -> (bool "false")
		if _, ok := (src.Expr).(Constant); ok {
			return Bool(false)
		}
		// (is_missing x), "miss(x, h)" -> (bool "true")
		if x := src.Expr; true {
			if miss(x, h) {
				return Bool(true)
			}
		}
		// (is_missing x), "TypeOf(x, h)&MissingType == 0" -> (bool "false")
		if x := src.Expr; true {
			if TypeOf(x, h)&MissingType == 0 {
				return Bool(false)
			}
		}
	case IsNotFalse:
		// (is_not_false (bool x)) -> (bool x)
		if x, ok := (src.Expr).(Bool); ok {
			return Bool(x)
		}
		// (is_not_false x), "TypeOf(x, h)&BoolType == 0" -> (bool "true")
		if x := src.Expr; true {
			if TypeOf(x, h)&BoolType == 0 {
				return Bool(true)
			}
		}
	case IsNotMissing:
		// (is_not_missing (missing)) -> (bool "false")
		if _, ok := (src.Expr).(Missing); ok {
			return Bool(false)
		}
		// (is_not_missing (constant _)) -> (bool "true")
		if _, ok := (src.Expr).(Constant); ok {
			return Bool(true)
		}
		// (is_not_missing x), "TypeOf(x, h) == MissingType" -> (bool "false")
		if x := src.Expr; true {
			if TypeOf(x, h) == MissingType {
				return Bool(false)
			}
		}
		// (is_not_missing x), "TypeOf(x, h)&MissingType == 0" -> (bool "true")
		if x := src.Expr; true {
			if TypeOf(x, h)&MissingType == 0 {
				return Bool(true)
			}
		}
	case IsNotNull:
		// (is_not_null (null)) -> (bool "false")
		if _, ok := (src.Expr).(Null); ok {
			return Bool(false)
		}
		// (is_not_null x), "null(x, h)" -> (bool "false")
		if x := src.Expr; true {
			if null(x, h) {
				return Bool(false)
			}
		}
		// (is_not_null x), "TypeOf(x, h)&NullType == 0" -> (bool "true")
		if x := src.Expr; true {
			if TypeOf(x, h)&NullType == 0 {
				return Bool(true)
			}
		}
	case IsNotTrue:
		// (is_not_true (bool x)) -> (bool "!x")
		if x, ok := (src.Expr).(Bool); ok {
			return Bool(!x)
		}
		// (is_not_true x), "TypeOf(x, h)&BoolType == 0" -> (bool "true")
		if x := src.Expr; true {
			if TypeOf(x, h)&BoolType == 0 {
				return Bool(true)
			}
		}
	case IsNull:
		// (is_null (null)) -> (bool "true")
		if _, ok := (src.Expr).(Null); ok {
			return Bool(true)
		}
		// (is_null x), "null(x, h)" -> (bool "true")
		if x := src.Expr; true {
			if null(x, h) {
				return Bool(true)
			}
		}
		// (is_null x), "TypeOf(x, h)&NullType == 0" -> (bool "false")
		if x := src.Expr; true {
			if TypeOf(x, h)&NullType == 0 {
				return Bool(false)
			}
		}
	case IsTrue:
		// (is_true (bool x)) -> (bool x)
		if x, ok := (src.Expr).(Bool); ok {
			return Bool(x)
		}
		// (is_true x), "TypeOf(x, h)&BoolType == 0" -> (bool "false")
		if x := src.Expr; true {
			if TypeOf(x, h)&BoolType == 0 {
				return Bool(false)
			}
		}
	}
	return nil
}

func simplifyClass4(src *Logical, h Hint) Node {
	switch src.Op {
	case OpAnd:
		// (and x x), "TypeOf(x, h) == LogicalType" -> x
		if x := src.Left; true {
			if x.Equals(src.Right) {
				if TypeOf(x, h) == LogicalType {
					return x
				}
			}
		}
		// (and (bool x) y), "x && TypeOf(y, h) == LogicalType" -> y
		if x, ok := (src.Left).(Bool); ok {
			if y := src.Right; true {
				if x && TypeOf(y, h) == LogicalType {
					return y
				}
			}
		}
		// (and x (bool y)), "y && TypeOf(x, h) == LogicalType" -> x
		if x := src.Left; true {
			if y, ok := (src.Right).(Bool); ok {
				if y && TypeOf(x, h) == LogicalType {
					return x
				}
			}
		}
		// (and (bool x) (bool y)) -> (bool "x && y")
		if x, ok := (src.Left).(Bool); ok {
			if y, ok := (src.Right).(Bool); ok {
				return Bool(x && y)
			}
		}
	case OpOr:
		// (or x x), "TypeOf(x, h) == LogicalType" -> x
		if x := src.Left; true {
			if x.Equals(src.Right) {
				if TypeOf(x, h) == LogicalType {
					return x
				}
			}
		}
		// (or (bool x) y), "!x", "TypeOf(y, h) == LogicalType" -> y
		if x, ok := (src.Left).(Bool); ok {
			if y := src.Right; true {
				if !x {
					if TypeOf(y, h) == LogicalType {
						return y
					}
				}
			}
		}
		// (or x (bool y)), "!y", "TypeOf(x, h) == LogicalType" -> x
		if x := src.Left; true {
			if y, ok := (src.Right).(Bool); ok {
				if !y {
					if TypeOf(x, h) == LogicalType {
						return x
					}
				}
			}
		}
		// (or (bool x) (bool y)) -> (bool "x || y")
		if x, ok := (src.Left).(Bool); ok {
			if y, ok := (src.Right).(Bool); ok {
				return Bool(x || y)
			}
		}
	case OpXnor:
		// (xnor (bool x) (bool y)) -> (bool "x == y")
		if x, ok := (src.Left).(Bool); ok {
			if y, ok := (src.Right).(Bool); ok {
				return Bool(x == y)
			}
		}
	case OpXor:
		// (xor (bool x) (bool y)) -> (bool "x != y")
		if x, ok := (src.Left).(Bool); ok {
			if y, ok := (src.Right).(Bool); ok {
				return Bool(x != y)
			}
		}
	}
	return nil
}
func simplify1(src Node, h Hint) Node {
	switch src := src.(type) {
	case *Arithmetic:
		return simplifyClass0(src, h)
	case *Builtin:
		return simplifyClass1(src, h)
	case *Comparison:
		return simplifyClass2(src, h)
	case *IsKey:
		return simplifyClass3(src, h)
	case *Logical:
		return simplifyClass4(src, h)
	}
	return nil
}
