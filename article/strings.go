package article

import (
	"bytes"
	"devtools/comerr"
	"log"
	"reflect"
	"strings"
	"text/template"
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
	if param == nil {
		return "", comerr.ParamInvalid
	} else {
		t := reflect.TypeOf(param)
		k := t.Kind()
		if k == reflect.Ptr {
			k = t.Elem().Kind()
		}
		if k != reflect.Struct && k != reflect.Map {
			return "", comerr.ParamInvalid
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

func Distinct(input []string) []string {
	output := make([]string, len(input))
	j := 0
	for i := 0; i < len(input); i++ {
		for _, v := range output {
			if v == input[i] {
				goto NEXT
			}
		}
		output[j] = input[i]
		j++
	NEXT:
	}

	return output[:j]
}

func CommonPrefix(s1, s2 string) string {
	var i = 0
	for ; i < len(s1) && i < len(s2) && s1[i] == s2[i]; i++ {
	}

	return s1[:i]
}
