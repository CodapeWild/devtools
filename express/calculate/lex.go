package calculate

import (
	"bytes"
	"log"
	"math/big"
	"strings"
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
	line    []byte
	i       int
	operand bytes.Buffer
}

func (this *calcLex) Lex(lval *calcSymType) int {
	log.Println("#########")
	for {
		op, data := this.next()
		switch op {
		case tag_eof:
			return int(tag_eof)
		case tag_operand:
			lval.num = &big.Rat{}
			if _, ok := lval.num.SetString(data); !ok {
				log.Println("set data to operand failed")

				return int(tag_eof)
			}
		case tag_operator:
			return operators[data]
		default:
			log.Println("unrecognized lexicon")
		}
	}
}

func (_ *calcLex) Error(e string) {
	log.Println(e)
}

func (this *calcLex) next() (tag, string) {
	for _, b := range this.line[this.i:] {
		switch b {
		case '\t', '\n', '\r':
			this.i++
		case '+', '-', '*', '/':
			if this.operand.Len() != 0 {
				operand := strings.TrimSpace(this.operand.String())
				this.operand = bytes.Buffer{}
				if len(operand) != 0 {
					return tag_operand, operand
				}
			}
			this.i++

			return tag_operator, string(b)
		default:
			if err := this.operand.WriteByte(b); err != nil {
				log.Panicln(err.Error())
			}
			this.i++
		}
	}

	if this.operand.Len() != 0 {
		operand := strings.TrimSpace(this.operand.String())
		this.operand = bytes.Buffer{}
		if len(operand) != 0 {
			return tag_operand, operand
		}
	}

	return tag_eof, ""
}
