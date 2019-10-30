package algorithm

import (
	"devtools/comerr"
	"errors"
	"reflect"
)

var (
	targetNotInSeq = errors.New("target item not included in sequence")
)

func BinarySearch(data OrderedData, target interface{}) (int, error) {
	if data == nil || data.Len() == 0 || reflect.TypeOf(data.Data(0)).Kind() != reflect.TypeOf(target).Kind() {
		return 0, comerr.ParamInvalid
	}

	var (
		l = 0
		r = data.Len() - 1
		m = 0
	)
	for {
		if l > r {
			return 0, targetNotInSeq
		}

		m = (l + r) / 2
		if data.Greater(m, target) {
			r = m - 1
		} else if data.Less(m, target) {
			l = m + 1
		} else {
			break
		}
	}

	return m, nil
}
