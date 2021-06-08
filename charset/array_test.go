package charset

import (
	"log"
	"testing"
)

func TestCharSet(t *testing.T) {
	set1 := []string{"1", "2", "3"}
	set2 := []string{"2", "2", "3"}
	log.Println("contain: ", Contains(set1, "3"))
	log.Println("unique: ", Unique(set2))
	log.Println("intersect: ", Intersect(set1, set2))
	log.Println("intersect: ", Intersect(set1, []string{}))
	log.Println("intersect: ", Intersect(nil, []string{}))
	log.Println("intersect: ", Intersect([]string{}, nil))
	log.Println("intersect: ", Intersect(nil, set2))
	log.Println("union: ", Union(set1, set2))
	log.Println("differ: ", Differ(set1, set2))
}
