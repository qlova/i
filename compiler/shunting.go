package compiler

//Precedence returns i's precedence for the specified symbol.
func Precedence(symbol []byte) int {
	if symbol == nil {
		return -1
	}
	switch string(symbol) {
	case ",", ")", "]", "\n", "", "}", "in", ":", ";", "to":
		return -1

	case "|":
		return 0

	case "&":
		return 1

	case "=", "<", ">", "!":
		return 2

	case "+", "-":
		return 3

	case "*", "/", `%`:
		return 4

	case "^":
		return 5

	case "(", "[":
		return 6

	case ".":
		return 7

	default:
		return -2
	}
}

//Shunt shunts an expression with the next part of the expression. Emplying operators.
func (c *Compiler) Shunt(expression Expression, precedence int) (result Expression, err error) {
	result = expression

	//shunting:
	for peek := c.Peek(); Precedence(peek) >= precedence; {

		if Precedence(c.Peek()) <= -1 {
			break
		}

		precedence := Precedence(peek)
		symbol := peek
		if c.Scan() == nil {
			break
		}

		lhs := result
		rhs, err := c.scanExpression()
		if err != nil {
			return result, err
		}

		peek = c.Peek()
		for Precedence(peek) > precedence {
			rhs, err = c.Shunt(rhs, Precedence(peek))
			if err != nil {
				return result, err
			}

			peek = c.Peek()
		}

		if symbol.Is("=") && rhs.Type == Newline {
			return Expression{Switch.With(result.Type), nil, result}, nil
		}

		//Indexing
		if symbol.Is("[") && lhs.Type == Data {
			if !c.ScanIf(']') {
				return Expression{}, c.NewError("expecting ']'")
			}
			return c.Shunt(Expression{Byte, c.Symbol(lhs.Value, rhs.Value), nil}, Precedence(peek))
		}

		//Equality
		if symbol.Is("=") && lhs.Type == rhs.Type {
			if lhs.Type == Byte {
				return Expression{Bit, c.Same(lhs.Value, rhs.Value), nil}, nil
			}
		}

		return Expression{}, c.NewError("Operator ", string(symbol),
			" does not apply to ", string(result.Type), " and ", string(rhs.Type))
	}
	return result, nil
}
