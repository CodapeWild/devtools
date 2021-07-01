package algorithm

type InterfsAdaptor interface {
	ToInterfaceAry() []interface{}
}

type ToInterfaceAryFunc func() []interface{}

func (this ToInterfaceAryFunc) ToInterfaceAry() []interface{} {
	return this()
}

func converter(racc RandomAccess) []interface{} {
	convert := make([]interface{}, racc.Len())
	k, v := racc.Next()
	for v != nil {
		if i, ok := k.(int); ok {
			convert[i] = v
		}
		k, v = racc.Next()
	}

	return convert
}

func IntsToInterfaceAry(data []int) []interface{} {
	return converter(NewIntsRandAcc(data))
}

func Float64sToInterfaceAry(data []float64) []interface{} {
	return converter(NewFloatsRandAcc(data))
}

func StringsToInterfaceAry(data []string) []interface{} {
	return converter(NewStringsRandAcc(data))
}
