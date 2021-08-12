package httpext

import (
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer os.Exit(m.Run())

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func TestFileServer(t *testing.T) {
	fsrv := &FileServer{}
	log.Println(http.ListenAndServe(":8080", fsrv))
}
