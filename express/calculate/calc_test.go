package calculate

import (
	"log"
	"testing"
)

var (
	line1 = []byte("1+2+3")
)

func TestNext(t *testing.T) {
	lex := newCalcLex(line1)
	tag, data := lex.next()
	for tag != tag_eof {
		log.Println(data)
		tag, data = lex.next()
	}
}

func TestCalculate(t *testing.T) {
	calcParse(&calcLex{line: []byte("1+2+3")})
}

func init() {
	calcErrorVerbose = true
}
