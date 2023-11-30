package account

// abc:123+def:a/b/c

// type ID string // [a-zA-Z0-9/_]

type Line string

func Term(t string) Line {
	return Line(t)
}

func Pair(p, q string) Line {
	return Line(p + ":" + q)
}

func Cat(p, q Line) Line {
	return Line(p + "+" + q)
}
