package algorithm

type PreviousIter interface {
	Previous() interface{}
}

type NextIter interface {
	Next() interface{}
}

type Iterator interface {
	PreviousIter
	NextIter
	Current() interface{}
}

type IntAryIter struct {
	data []int
	i    int
}

func (this *IntAryIter) Previous() interface{} {
	if this.i == 0 {
		return nil
	}

	this.i--

	return this.data[this.i]
}

func (this *IntAryIter) Next() interface{} {
	if this.i == len(this.data)-1 {
		return nil
	}

	this.i++

	return this.data[this.i]
}

func (this *IntAryIter) Current() interface{} {
	if len(this.data) != 0 && this.i < len(this.data) {
		return this.data[this.i]
	} else {
		return nil
	}
}

type StringIter struct {
	s string
	i int
}

func (this *StringIter) Previous() interface{} {
	if this.i == 0 {
		return nil
	}

	this.i--

	return this.s[this.i]
}

func (this *StringIter) Next() interface{} {
	if this.i == len(this.s)-1 {
		return nil
	}

	this.i++

	return this.s[this.i]
}

func (this *StringIter) Current() interface{} {
	if len(this.s) != 0 && this.i < len(this.s) {
		return this.s[this.i]
	} else {
		return nil
	}
}

type StringsIter struct {
	ss []string
	i  int
}

func (this *StringsIter) Previous() interface{} {
	if this.i == 0 {
		return nil
	}

	this.i--

	return this.ss[this.i]
}

func (this *StringsIter) Next() interface{} {
	if this.i == len(this.ss)-1 {
		return nil
	}

	this.i++

	return this.ss[this.i]
}

func (this *StringsIter) Current() interface{} {
	if len(this.ss) != 0 && this.i < len(this.ss) {
		return this.ss[this.i]
	} else {
		return nil
	}
}
