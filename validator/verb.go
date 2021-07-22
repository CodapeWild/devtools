package validator

import (
	"regexp"
	"strings"

	"github.com/CodapeWild/devtools/charset"
)

var countries = map[string]*regexp.Regexp{
	"86": regexp.MustCompile(`^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\d{8}$`),
}

func IsPhone(countryCode, phone string) bool {
	return phone != "" && countries[strings.TrimLeft(countryCode, "+")].MatchString(phone)
}

var emailValidator = regexp.MustCompile(`^([0-9a-zA-Z]([-.\w]*[0-9a-zA-Z])*@([0-9a-zA-Z][-\w]*[0-9a-zA-Z]\.)+[a-zA-Z]{2,9})$`)

func IsEmail(email string) bool {
	return email != "" && emailValidator.MatchString(email)
}

type PswdLevel byte

const (
	Pswd_Lv1 PswdLevel = iota + 1
	Pswd_Lv2
	Pswd_Lv3
)

func hasProperLength(pswd string) bool {
	return len(pswd) >= 6 && len(pswd) <= 36
}

func hasSpecialChar(pswd string) bool {
	for _, v := range pswd {
		if v < 33 || v > 126 {
			return true
		}
	}

	return false
}

func hasSameNChar(pswd string, n int) bool {
	if n < 2 || n > len(pswd)/2 {
		return false
	}

	r := 1
	for i := 1; i < len(pswd); i++ {
		if pswd[i] == pswd[i-1] {
			if r++; r >= n {
				return true
			}
		} else {
			r = 1
		}
	}

	return false
}

func hasRepeatedLNSub(pswd string) bool {
	for i := 2; i <= len(pswd)/2; i++ {
		if charset.RepeatedLNSub(pswd, i) != "" {
			return true
		}
	}

	return false
}

func hasCharType(pswd string) int {
	var b byte
	for _, v := range pswd {
		switch {
		case (v >= '!' && v <= '/') || (v >= ':' && v <= '@') || (v >= '[' && v <= '`') || (v >= '{' && v <= '~'):
			b |= 1
		case v >= '0' && v <= '9':
			b |= 2
		case (v >= 'A' && v <= 'Z') || (v >= 'a' && v <= 'z'):
			b |= 4
		}
	}
	var c int
	for _, v := range []byte{1, 2, 4} {
		if b&v != 0 {
			c++
		}
	}

	return c
}

func IsPswdStrongEnough(pswd string, level PswdLevel) bool {
	if len(pswd) == 0 {
		return false
	}

	switch level {
	case Pswd_Lv1:
		return hasProperLength(pswd) && !hasSpecialChar(pswd) && !hasSameNChar(pswd, 3) && !hasRepeatedLNSub(pswd)
	case Pswd_Lv2:
		return hasProperLength(pswd) && !hasSpecialChar(pswd) && !hasSameNChar(pswd, 3) && !hasRepeatedLNSub(pswd) && hasCharType(pswd) >= 2
	case Pswd_Lv3:
		return hasProperLength(pswd) && !hasSpecialChar(pswd) && !hasSameNChar(pswd, 3) && !hasRepeatedLNSub(pswd) && hasCharType(pswd) == 3
	default:
		return false
	}
}

var nicknameValidator = regexp.MustCompile(`[\S2E80～33FFh]{6,15}`)

func IsNickname(nickname string) bool {
	return nickname != "" && nicknameValidator.MatchString(nickname)
}
