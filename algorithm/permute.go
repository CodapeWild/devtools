package algorithm

type Permutable interface {
	Len() int
	NewEmpty() interface{}
	Copy() interface{}
	Append(interface{}) interface{}
	At(int) interface{}
	Remove(int) interface{}
}

func permuteHandler(src Permutable, permutation Permutable, output chan interface{}) {
	if src.Len() == 0 {
		output <- permutation
	} else {
		for i := 0; i < src.Len(); i++ {
			ptmp := permutation.Copy().(Permutable)
			ptmp = ptmp.Append(src.At(i)).(Permutable)

			srctmp := src.Copy().(Permutable)
			srctmp = srctmp.Remove(i).(Permutable)

			permuteHandler(srctmp, ptmp, output)
		}
	}
}

func Permute(src Permutable, output chan interface{}) {
	permuteHandler(src, src.NewEmpty().(Permutable), output)
}

type IntsPermutable []int

func (this IntsPermutable) Len() int {
	return len(this)
}

func (this IntsPermutable) NewEmpty() interface{} {
	return IntsPermutable([]int{})
}

func (this IntsPermutable) Copy() interface{} {
	tmp := make([]int, len(this))
	copy(tmp, this)

	return IntsPermutable(tmp)
}

func (this IntsPermutable) Append(a interface{}) interface{} {
	return append(this, a.(int))
}

func (this IntsPermutable) At(i int) interface{} {
	return this[i]
}

func (this IntsPermutable) Remove(i int) interface{} {
	this[i] = this[len(this)-1]

	return this[:len(this)-1]
}
