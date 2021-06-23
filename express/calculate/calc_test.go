package calculate

import (
	"log"
	"testing"
)

func TestNext(t *testing.T) {
	lex := &calcLex{line: []byte("1+2+3")}
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
