package article

import (
	"log"
	"testing"
)

var src = "abcdefghijk\nlm\nnopqistuvwxyz a bc de fghij klmnopq istuvwxyz 来自星 星的你 abcdefgh ijk lmnop qistuvwxyz abcd efghijklmnopqistuvwxyz\r\n"

func TestStrings(t *testing.T) {
	log.Println(FoldByMax(src, 24))
}

func TestCommonPrefixLen(t *testing.T) {
	log.Println(CommonPrefixLen("123", "132"))
	log.Println(CommonPrefixLen("123", "0132"))
}
