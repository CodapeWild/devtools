package calculate

import (
	"bytes"
	"log"
	"math/big"
)

type tag int

const (
	tag_eof tag = iota
	tag_operand
	tag_operator
)

var operators = map[string]int{
	"+": ADD,
	"-": SUB,
	"*": MUL,
	"/": DIV,
	"(": LEFT_BRACE,
	")": RIGHT_BRACE,
}

type calcLex struct {
	line    *bytes.Buffer
	preRead *bytes.Buffer
}

func newCalcLex(line []byte) *calcLex {
	return &calcLex{
		line:    bytes.NewBuffer(line),
		preRead: &bytes.Buffer{},
	}
}

func (this *calcLex) Lex(lval *calcSymType) int {
	for {
		t, key := this.next()
		switch t {
		case tag_eof:
			return int(tag_eof)
		case tag_operator:
			return operators[key]
		case tag_operand:
			r := &big.Rat{}
			if _, ok := r.SetString(key); !ok {
				log.Printf("bad operand %s\n", key)

				return int(tag_eof)
			}

			return NUM
		}
	}
}

func (_ *calcLex) Error(e string) {
	log.Println(e)
}

func (this *calcLex) next() (tag, string) {
	for {
		switch {
		case condition:

		}
	}
}
