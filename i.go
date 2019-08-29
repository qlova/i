package i

import (
	"math/big"
	"strconv"
)

//Bit is the 'i' boolean type.
type Bit = bool

//String is the 'i' string type.
type String = string

//Error is the 'i' error type.
type Error struct {
	Code    int
	Message string
}

type context struct {
	errors []Error
}

//Context is the 'i' context for error handling and globals.
type Context struct {
	*context
}

//Errors returns and clears the error stack.
func (context Context) Errors() []Error {
	defer func() {
		context.errors = nil
	}()
	return context.errors
}

//NewContext returns a new 'i' context.
func NewContext() Context {
	return Context{new(context)}
}

//Throw throws an error onto the stack.
func (ctx *Context) Throw(code int, message string) {
	ctx.errors = append(ctx.errors, Error{
		Code:    code,
		Message: message,
	})
}

//Integer is the 'i' integer type.
type Integer struct {
	int64
	large big.Int
}

//IndexList returns an index to a list with the specified length.
func IndexList(index Integer, length int) int {
	return int(index.Mod(NewInteger(length-1)).Int64()) + 1
}

//IndexArray returns an index to an array with the specified length.
func IndexArray(index Integer, length int) int {
	return int(index.Mod(NewInteger(length)).Int64())
}

//Atoi returns a String converted to an Integer.
func Atoi(ctx Context, s String) Integer {
	i, err := strconv.Atoi(s)
	if err != nil {
		l, ok := big.NewInt(0).SetString(s, 10)
		if !ok {
			ctx.Throw(1, "invalid integer")
			return Integer{}
		}
		return Integer{large: *l}
	}
	return Integer{int64: int64(i)}
}

//NewInteger returns an 'i' integer from an int64.
func NewInteger(i int) Integer {
	return Integer{int64: int64(i)}
}

//Neg returns -i
func (a Integer) Neg() Integer {
	if a.large.Bits() == nil {
		a.large.Neg(&a.large)
		return a
	}
	return Integer{int64: 0 - a.int64}
}

//Add returns a + b
func (a Integer) Add(b Integer) Integer {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:

		c := a.int64 + b.int64

		if (c > a.int64) == (b.int64 > 0) {
			return Integer{int64: c}
		}

		//Overflow.
		a.large.Add(big.NewInt(a.int64), big.NewInt(b.int64))
		return a
	case aSmall:
		a.large.Add(big.NewInt(a.int64), &b.large)
		return a
	case bSmall:
		a.large.Add(&a.large, big.NewInt(b.int64))
		return a
	default:
		a.large.Add(&a.large, &b.large)
		return a
	}
}

//Sub returns a - b
func (a Integer) Sub(b Integer) Integer {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:

		c := a.int64 - b.int64

		if (c < a.int64) == (b.int64 > 0) {
			return Integer{int64: c}
		}

		//Overflow.
		a.large.Sub(big.NewInt(a.int64), big.NewInt(b.int64))
		return a
	case aSmall:
		a.large.Sub(big.NewInt(a.int64), &b.large)
		return a
	case bSmall:
		a.large.Sub(&a.large, big.NewInt(b.int64))
		return a
	default:
		a.large.Sub(&a.large, &b.large)
		return a
	}
}

//Mul returns a * b
func (a Integer) Mul(b Integer) Integer {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:

		if a.int64 == 0 || b.int64 == 0 {
			return Integer{}
		}

		c := a.int64 * b.int64
		if (c < 0) == ((a.int64 < 0) != (b.int64 < 0)) {
			if c/b.int64 == a.int64 {
				return Integer{int64: c}
			}
		}

		//Overflow.
		a.large.Mul(big.NewInt(a.int64), big.NewInt(b.int64))
		return a
	case aSmall:
		a.large.Mul(big.NewInt(a.int64), &b.large)
		return a
	case bSmall:
		a.large.Mul(&a.large, big.NewInt(b.int64))
		return a
	default:
		a.large.Mul(&a.large, &b.large)
		return a
	}
}

//Div returns a / b
func (a Integer) Div(b Integer) Integer {
	if b.Equals(Integer{}) {
		if a.Equals(Integer{}) {
			return Integer{int64: 1}
		}
		return Integer{int64: 0}
	}

	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:

		c := a.int64 / b.int64
		if (c < 0) == ((a.int64 < 0) != (b.int64 < 0)) {
			return Integer{int64: c}
		}

		//Overflow.
		a.large.Div(big.NewInt(a.int64), big.NewInt(b.int64))
		return a
	case aSmall:
		a.large.Div(big.NewInt(a.int64), &b.large)
		return a
	case bSmall:
		a.large.Div(&a.large, big.NewInt(b.int64))
		return a
	default:
		a.large.Div(&a.large, &b.large)
		return a
	}
}

//Pow returns a ** b
func (a Integer) Pow(b Integer) Integer {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:
		a.large.Exp(big.NewInt(a.int64), big.NewInt(b.int64), nil)
		return a
	case aSmall:
		a.large.Exp(big.NewInt(a.int64), &b.large, nil)
		return a
	case bSmall:
		a.large.Exp(&a.large, big.NewInt(b.int64), nil)
		return a
	default:
		a.large.Exp(&a.large, &b.large, nil)
		return a
	}
}

//Mod returns a % b
func (a Integer) Mod(b Integer) Integer {
	if b.Equals(Integer{}) {
		if a.Equals(Integer{}) {
			return Integer{int64: 1}
		}
		return Integer{int64: 0}
	}

	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:

		c := a.int64 % b.int64
		return Integer{int64: c}

	case aSmall:
		a.large.Mod(big.NewInt(a.int64), &b.large)
		return a
	case bSmall:
		a.large.Mod(&a.large, big.NewInt(b.int64))
		return a
	default:
		a.large.Mod(&a.large, &b.large)
		return a
	}
}

//Compare returns a == b.
func (a Integer) Compare(b Integer) int {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:
		switch true {
		case a.int64 < b.int64:
			return -1
		case a.int64 > b.int64:
			return 1
		default:
			return 0
		}
	case aSmall:
		return big.NewInt(a.int64).Cmp(&b.large)
	case bSmall:
		return a.large.Cmp(big.NewInt(b.int64))
	default:
		return a.large.Cmp(&b.large)
	}
}

//Equals returns a == b.
func (a Integer) Equals(b Integer) Bit {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:
		return a.int64 == b.int64
	case aSmall:
		return big.NewInt(a.int64).Cmp(&b.large) == 0
	case bSmall:
		return a.large.Cmp(big.NewInt(b.int64)) == 0
	default:
		return a.large.Cmp(&b.large) == 0
	}
}

//GreaterThan returns a > b.
func (a Integer) GreaterThan(b Integer) Bit {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:
		return a.int64 > b.int64
	case aSmall:
		return big.NewInt(a.int64).Cmp(&b.large) > 0
	case bSmall:
		return a.large.Cmp(big.NewInt(b.int64)) > 0
	default:
		return a.large.Cmp(&b.large) > 0
	}
}

//LessThan returns a < b.
func (a Integer) LessThan(b Integer) Bit {
	aSmall, bSmall := a.large.Bits() == nil, b.large.Bits() == nil
	switch true {
	case aSmall && bSmall:
		return a.int64 < b.int64
	case aSmall:
		return big.NewInt(a.int64).Cmp(&b.large) < 0
	case bSmall:
		return a.large.Cmp(big.NewInt(b.int64)) < 0
	default:
		return a.large.Cmp(&b.large) < 0
	}
}

//To returns a + norm(a - b).
func (a Integer) To(b Integer) Integer {
	if a.GreaterThan(b) {
		return a.Sub(Integer{int64: 1})
	}
	return a.Add(Integer{int64: 1})
}

//SetupTo returns a setup for a to loop.
func SetupTo(a, b Integer) (Integer, Integer) {
	if a.GreaterThan(b) {
		return a, b.Sub(Integer{int64: 1})
	}
	return a, b.Add(Integer{int64: 1})
}

//SetupStep returns a setup for a stepped loop.
func SetupStep(to, step Integer) (Integer, Integer, Integer) {
	return step, step, to
}

//CompareStep compares a stepped loop's step.
func (a Integer) CompareStep(to, step Integer) bool {
	if a.GreaterThan(to) {
		return false
	}
	return true
}

func (a Integer) String() String {
	if a.large.Bits() == nil {
		return strconv.FormatInt(a.int64, 10)
	}
	return a.large.String()
}

//Int64 returns the integer as an int64.
func (a Integer) Int64() int64 {
	if a.large.Bits() == nil {
		return a.int64
	}
	return a.large.Int64()
}
