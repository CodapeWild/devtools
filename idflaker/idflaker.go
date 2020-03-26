package idflaker

import (
	"devtools/clock"
	"devtools/comerr"
	"encoding/base64"
	"encoding/binary"
	"sync"
)

/*
	64bits id flaker from left to right
	  1bit: sign
	 8bits: id
	11bits: sequence number
	42bits:	timestamp(millisecond)
*/

const (
	seq_mask = ^(-1 << 11)
	ts_mask  = ^(-1 << 42)
)

type IdFlaker struct {
	id  int64
	seq int64
	ts  int64
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

	return this.id | this.seq<<42 | this.ts
}

func (this *IdFlaker) until() {
	for {
		if now := clock.NowMillisec(); now != this.ts {
			this.ts = now
			break
		}
	}
}

func (this *IdFlaker) NextBase64Id(encode *base64.Encoding) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(this.NextInt64Id()))

	return encode.EncodeToString(buf)
}

func ParseInt64Id(flkId int64) (id, seq, ts int64) {
	return flkId >> 53, flkId >> 42 & seq_mask, flkId & ts_mask
}

func ParseBase64Id(flkId string, encode *base64.Encoding) (id int64, err error) {
	buf := make([]byte, 8)
	_, err = encode.Decode(buf, []byte(flkId))
	if err != nil {
		return 0, err
	}

	return int64(binary.BigEndian.Uint64(buf)), nil
}
