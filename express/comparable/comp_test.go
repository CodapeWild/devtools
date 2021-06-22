package comparable

import (
	"log"
	"testing"
)

func TestResetCompRsltString(t *testing.T) {
	log.Println(Comp_True, Comp_False)
	ResetCompRsltString()
	log.Println(Comp_True, Comp_False)
}
