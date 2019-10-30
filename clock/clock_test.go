package clock

import (
	"log"
	"testing"
)

func TestClock(t *testing.T) {
	log.Println(FormatSeconds(3784))
	log.Println(ParseSeconds("01:03:04"))
}
