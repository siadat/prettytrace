package testpkg

type S struct{}

func (S) f1() {
	panic("panicking in S's func f1()")
}

func (*S) f2() {
	panic("panicking in S's func f2()")
}

func F1() {
	var s = S{}
	s.f1()
}

func F2() {
	var s = &S{}
	s.f2()
}

func DivideByZero() {
	divideByZero()
}

func divideByZero() {
	var x = 1
	var y = 0
	var _ = x / y
}
