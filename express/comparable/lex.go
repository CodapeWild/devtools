package comparable

import (
	"bytes"
	"log"
	"unicode/utf8"
)

const eof = 0

var singles = map[rune]int{
	';': SEMICOLON,
	'{': LEFT_CURLY_BRACKET,
	'}': RIGHT_CULY_BRACKET,
	',': COMMA,
	'>': GREAT,
	'<': LITTLE,
	'(': LEFT_BRACE,
	')': RIGHT_BRACE,
}

var multiples = map[string]int{
	"&&":     AND,
	"||":     OR,
	"==":     EQUAL,
	"!=":     NOT_EQUAL,
	">=":     GREAT_EQUAL,
	"<=":     LITTLE_EQUAL,
	"in":     IN,
	"not_in": NOT_IN,
	"match":  MATCH,
}

type compLex struct {
	line []byte
	peek rune
}

// return operator and Comparable
func (this *compLex) Lex(lval *compSymType) int {
	for {
		r := this.next()
		switch r {
		case ';', '{', '}', ',', '(', ')':
			return singles[r]
		case '>', '<':
			if n := this.next(); n == '=' {
				return multiples[string([]rune{r, n})]
			} else {
				this.peek = n

				return singles[r]
			}
		case '&', '|', '=':
			if n := this.next(); n == r {
				return multiples[string([]rune{r, n})]
			} else {
				this.peek = n

			}
		case '!':
		case 'i', 'n', 'm':
		case ' ', '\t', '\n', 'r':
		default:
			// operand either be wrapped with "", `` or can not contain any spaces
			return this.comp(r, lval)
		}
	}
}

func (this *compLex) Error(e string) {

}

func (this *compLex) next() rune {
	if this.peek != eof {
		r := this.peek
		this.peek = eof

		return r
	}

	if len(this.line) == 0 {
		return eof
	}

	r, size := utf8.DecodeRune(this.line)
	this.line = this.line[size:]
	if r == utf8.RuneError && size == 1 {
		log.Println("invalid utf8 character")

		return this.next()
	}

	return r
}

func (this *compLex) comp(r rune, lval *compSymType) int {
	if r == '"' || r == '`' {
		if i := bytes.IndexByte(this.line, byte(r)); i < 0 {
			return eof
		} else {
			lval.comp = Comparable(string(this.line[:i]))
			this.line = this.line[i+1:]

			return COMPARABLE
		}
	} else {

	}
}
