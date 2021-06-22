package comparable

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"regexp"
	"strconv"
)

var (
	ErrUnrecognizedOperator = errors.New("unrecognized operator")
	ErrCompareIncomplete    = errors.New("relational expression not calculated")
)

var (
	Comp_True  Comparable = "6PZiy4ppq2"
	Comp_False Comparable = "JPlSjXf25s"
)

type Comparable string

func (lopr Comparable) Equal(ropr Comparable) Comparable {
	return convert(calc(lopr, ropr, "=="))
}

func (lopr Comparable) NotEqual(ropr Comparable) Comparable {
	return negate(lopr.Equal(ropr))
}

func (lopr Comparable) GreatThan(ropr Comparable) Comparable {
	return convert(calc(lopr, ropr, ">"))
}

func (lopr Comparable) LittleThan(ropr Comparable) Comparable {
	return convert(calc(lopr, ropr, "<"))
}

func (lopr Comparable) GreatEqualThan(ropr Comparable) Comparable {
	return convert(calc(lopr, ropr, ">="))
}

func (lopr Comparable) LittleEqualThan(ropr Comparable) Comparable {
	return convert(calc(lopr, ropr, "<="))
}

func (lopr Comparable) In(roprs ...Comparable) Comparable {
	for _, ropr := range roprs {
		if lopr.Equal(ropr) == Comp_True {
			return Comp_True
		}
	}

	return Comp_False
}

func (lopr Comparable) NotIn(roprs ...Comparable) Comparable {
	return negate(lopr.In(roprs...))
}

func (lopr Comparable) And(ropr Comparable) Comparable {
	return convert((lopr == Comp_True) && (ropr == Comp_True))
}

func (lopr Comparable) Or(ropr Comparable) Comparable {
	return convert((lopr == Comp_True) || (ropr == Comp_True))
}

func (lopr Comparable) string() string {
	return string(lopr)
}

func (lopr Comparable) float() (float64, error) {
	return strconv.ParseFloat(lopr.string(), 64)
}

func (lopr Comparable) int() (int64, error) {
	return strconv.ParseInt(lopr.string(), 10, 64)
}

func (lopr Comparable) boolean() (bool, error) {
	if lopr == Comp_True {
		return true, nil
	} else if lopr == Comp_False {
		return false, nil
	} else {
		return false, ErrCompareIncomplete
	}
}

// !!!important calling this function only after finish the whole expression calculation
func ResetCompRsltString() {
	buf := make([]byte, 30)
	rand.Read(buf)
	Comp_True = Comparable(base64.StdEncoding.EncodeToString(buf))
	rand.Read(buf)
	Comp_False = Comparable(base64.StdEncoding.EncodeToString(buf))
}

func Match(target, reg string) Comparable {
	re := regexp.MustCompile(reg)

	return convert(re.MatchString(target))
}

func convert(b bool) Comparable {
	if b {
		return Comp_True
	} else {
		return Comp_False
	}
}

func negate(comp Comparable) Comparable {
	if comp.string() == Comp_True.string() {
		return Comp_False
	} else if comp.string() == Comp_False.string() {
		return Comp_True
	} else {
		return comp
	}
}

func calc(lopr, ropr Comparable, operator string) bool {
	var (
		lf, rf *float64
		li, ri *int64
		ls, rs *string
		err    error
	)
	if *lf, err = lopr.float(); err == nil {
		if *rf, err = lopr.float(); err == nil {
			goto CALC
		}
	}
	if *li, err = lopr.int(); err == nil {
		if *ri, err = ropr.int(); err == nil {
			goto CALC
		}
	}
	*ls, *rs = lopr.string(), ropr.string()

CALC:
	switch operator {
	case "==":
		return (lf != nil && *lf == *rf) || (li != nil && *li == *ri) || (ls != nil && *ls == *rs)
	case ">":
		return (lf != nil && *lf > *rf) || (li != nil && *li > *ri) || (ls != nil && *ls > *rs)
	case ">=":
		return (lf != nil && *lf >= *rf) || (li != nil && *li >= *ri) || (ls != nil && *ls >= *rs)
	case "<":
		return (lf != nil && *lf < *rf) || (li != nil && *li < *ri) || (ls != nil && *ls < *rs)
	case "<=":
		return (lf != nil && *lf <= *rf) || (li != nil && *li <= *ri) || (ls != nil && *ls <= *rs)
	default:
		log.Println(ErrUnrecognizedOperator)
	}

	return false
}
