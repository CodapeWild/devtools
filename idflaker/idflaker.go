package idflaker

import (
	"devtools/clock"
	"devtools/comerr"
	"errors"
	"strconv"
	"sync"
)

/*
	int64 id flaker
	  1bit: sign
	 8bits: id
	42bits:	timestamp(millisecond)
	11bits: sequence number
*/

const (
	seq_mask = ^(-1 << 11)
	ts_mask  = ^(-1 << 42) << 11
)

var (
	StrIdInvalid = errors.New("invalid string id for idflaker")
)

type IdFlaker struct {
	id  int64
	ts  int64
	seq int64
	sync.Mutex
}

func NewIdFlaker(id int64) (*IdFlaker, error) {
	if id < 0 || id > 255 {
		return nil, comerr.ParamInvalid
	}

	return &IdFlaker{id: id << 53}, nil
}

func (this *IdFlaker) NextInt64Id() int64 {
	this.Lock()
	defer this.Unlock()

	if now := clock.NowMillisec(); now == this.ts {
		this.seq++
		if this.seq &= seq_mask; this.seq == 0 {
			this.until()
		}
	} else {
		this.ts = now
		this.seq = 0
	}

	return this.id | this.ts<<11 | this.seq
}

func (this *IdFlaker) NextStrId() string {
	return strconv.FormatInt(this.NextInt64Id(), 10)
}

func ParseInt64Id(numId int64) (id, ts, seq int64) {
	id = numId >> 53
	ts = (numId & ts_mask) >> 11
	seq = numId & seq_mask

	return
}

func ParseStrId(strId string) (id, ts, seq int64, err error) {
	if strId == "" || strId == "0" {
		return 0, 0, 0, StrIdInvalid
	}

	var numId int64
	if numId, err = strconv.ParseInt(strId, 10, 64); err != nil {
		return
	} else {
		id, ts, seq = ParseInt64Id(numId)

		return
	}
}

func (this *IdFlaker) until() {
	for {
		if now := clock.NowMillisec(); now != this.ts {
			this.ts = now
			break
		}
	}
}
