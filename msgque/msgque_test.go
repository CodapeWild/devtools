package msgque

import (
	"log"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	log.SetFlags(log.LstdFlags | log.LstdFlags)
	log.SetOutput(os.Stdout)

	os.Exit(t.Run())
}
