package clock

import (
	"log"
	"testing"
	"time"
)

func TestSeconds(t *testing.T) {
	log.Println(FormatSeconds(3784))
	log.Println(ParseSeconds("01:03:04"))
}

func TestParse(t *testing.T) {
	log.Println(ParseUnixSec(1619771900).Format(time.ANSIC))
	log.Println(ParseUnixMillsec(1619771900 * 1000).Format(time.ANSIC))
}
