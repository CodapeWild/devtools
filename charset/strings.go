package charset

import (
	"bytes"
	"log"
	"reflect"
	"strings"
	"text/template"

	"github.com/CodapeWild/devtools/comerr"
)

func FoldByMax(src string, maxLineChars int) string {
	if len(src) <= maxLineChars {
		return src
	}

	var (
		p, q int // p: last write stop point; q: last blank point;
		dst  = ""
	)
	for i := 0; i < len(src); i++ {
		if src[i] == ' ' {
			if i-p < maxLineChars {
				q = i
			} else {
				if p > q {
					dst += src[p:i] + "\r\n"
					p = i + 1
				} else {
					dst += src[p:q] + "\r\n"
					p = q + 1
					if i-p >= maxLineChars {
						dst += src[p:i] + "\r\n"
						p = i + 1
					}
				}
				q = i
			}
		} else if src[i] == '\n' {
			if i-p > maxLineChars && q > p {
				dst += src[p:q] + "\r\n"
				p = q + 1
			}
			dst += src[p : i+1]
			p = i + 1
		}
	}
	if p < q {
		dst += src[p:q] + "\r\n"
		p = q + 1
	}
	dst += src[p:]

	log.Println(dst)

	return dst
}

func GoTemplateReplace(tmpl string, param interface{}) (string, error) {
	if rv := reflect.ValueOf(param); rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "", comerr.ErrParamInvalid
		} else {
			rv = rv.Elem()
		}
		if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Map {
			return "", comerr.ErrParamInvalid
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := template.Must(template.New("").Parse(tmpl)).Execute(buf, param); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// find repeated substing of length n
func RepeatedLNSub(s string, count int) string {
	if count <= 0 || count > len(s)/2 {
		return ""
	}

	for k := range s {
		if k >= len(s)-count-1 {
			return ""
		}
		if strings.Contains(s[k+count:], s[k:k+count]) {
			return s[k : k+count]
		}
	}

	return ""
}

func CommonPrefix(s1, s2 string) string {
	var i = 0
	for ; i < len(s1) && i < len(s2) && s1[i] == s2[i]; i++ {
	}

	return s1[:i]
}

func FindUInt(s string) (num uint, end int, err error) {
	stop := false
	for k, v := range s {
		if v >= '0' && v <= '9' {
			num = num*10 + uint(v-'0')
			stop = true
		} else {
			if stop {
				return num, k, nil
			}
		}
	}

	return num, 0, comerr.ErrNotFound
}
