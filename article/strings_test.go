package article

import (
	"fmt"
	"testing"
)

var src = "abcdefghijk\nlm\nnopqistuvwxyz a bc de fghij klmnopq istuvwxyz 来自星 星的你 abcdefgh ijk lmnop qistuvwxyz abcd efghijklmnopqistuvwxyz\r\n"

func TestStrings(t *testing.T) {
	fmt.Println(FoldByMax(src, 24))
}
