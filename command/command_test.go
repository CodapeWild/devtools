package command

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	str := "ffprobe -v quiet -print_format json -show_streams -show_format s3.flv"
	cmds := strings.Split(str, " ")
	output := make(map[string]interface{})
	if err := (&Command{exec.Command(cmds[0], cmds[1:]...)}).StdoutJson(&output); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(output)
	}
}
