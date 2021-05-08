package clock

import (
	"devtools/comerr"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NowMillisec() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func ParseSeconds(s string) (int, error) {
	p := strings.Split(s, ":")
	if len(p) > 3 {
		return 0, comerr.ErrParamInvalid
	}
	opters := []int{1, 60, 3600}
	j := 0
	sec := 0
	for i := len(p) - 1; i >= 0; i-- {
		if k, err := strconv.Atoi(p[i]); err != nil {
			return 0, err
		} else {
			sec += k * opters[j]
			j++
		}
	}

	return sec, nil
}

func FormatSeconds(sec int) string {
	if sec <= 0 {
		return "00:00:00"
	} else {
		h := sec / 3600
		m := (sec - h*3600) / 60

		return fmt.Sprintf("%0.2d:%0.2d:%0.2d", h, m, sec-h*3600-m*60)
	}
}
