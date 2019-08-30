package i

//InSymbol reads stdin and returns a string up to and excluding 'symbol'.
func InSymbol(ctx Context, symbol rune) string {
	s, err := Stdin.ReadString(byte(symbol))
	if err != nil {
		ctx.Throw(1, err.Error())
		return ""
	}
	return s[:len(s)-1]
}
