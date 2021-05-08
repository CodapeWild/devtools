package algorithm

type Permutable interface {
	NewEmpty() Permutable
	Len() int
	Copy() Permutable
	Append(value interface{}) Permutable
	At(i int) interface{}
	Set(i int, value interface{})
	Range(i, j int) Permutable
}

func permuteHandler(src Permutable, permutation Permutable, output chan interface{}) {
	if src == nil {
		return
	}

	if src.Len() == 0 {
		output <- permutation
	} else {
		for i := 0; i < src.Len(); i++ {
			pertmp := permutation.Copy().(Permutable)
			pertmp = pertmp.Append(src.At(i)).(Permutable)

			srctmp := src.Copy().(Permutable)
			srctmp.Set(i, srctmp.At(srctmp.Len()-1))

			permuteHandler(srctmp.Range(0, srctmp.Len()-1), pertmp, output)
		}
	}
}

func Permute(src Permutable, output chan interface{}) {
	permuteHandler(src, src.NewEmpty().(Permutable), output)
}

/*
	ints permutable array
*/
type IntsPermutable []int

func (this IntsPermutable) NewEmpty() Permutable {
	return IntsPermutable{}
}

func (this IntsPermutable) Len() int {
	return len(this)
}

func (this IntsPermutable) Copy() Permutable {
	tmp := make([]int, len(this))
	copy(tmp, this)

	return IntsPermutable(tmp)
}

func (this IntsPermutable) Append(value interface{}) Permutable {
	return append(this, value.(int))
}

func (this IntsPermutable) At(i int) interface{} {
	if i >= 0 && i < len(this) {
		return this[i]
	} else {
		return nil
	}
}

func (this IntsPermutable) Set(i int, value interface{}) {
	if i >= 0 && i < len(this) {
		this[i] = value.(int)
	}
}

func (this IntsPermutable) Range(i, j int) Permutable {
	if i >= 0 && i < len(this) && j >= i {
		return this[i:j]
	} else {
		return nil
	}
}
