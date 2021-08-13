package cache

var defCache Cache = &MemCache{}

type Cache interface {
	Push(obj interface{}) bool
	Pop() interface{}
	Clear()
	Len() int
}

func UseCache(cache Cache) (restore func()) {
	old := defCache
	if cache != nil {
		defCache = cache
	}

	return func() {
		defCache = old
	}
}

func Push(obj interface{}) bool {
	return defCache.Push(obj)
}

func Pop() interface{} {
	return defCache.Pop()
}

func Clear() {
	defCache.Clear()
}

func Len() int {
	return defCache.Len()
}
