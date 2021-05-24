package command

import (
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestCmd(t *testing.T) {
	output, err := exec.Command("smartctl", "/dev/disk1s1", "-s").CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(string(output))
}

func TestCommand(t *testing.T) {
	// cmd, err := CmdFromCmd("smartctl", "-a", "/dev/disk1s1")
	cmd, err := CmdFromCmd("smartctl", "/dev/disk1s1", "-s")
	// cmd, err := CmdFromCmd("smartctl", "/dev/disk1s1")
	// cmd, err := CmdFromCmd("smartctl")
	if err != nil {
		log.Panicln(err.Error())
	}

	// output, err := cmd.RunWithTimeout(3 * time.Second)
	output, err := cmd.CombineOutErr().RunWithTimeout(3 * time.Second)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(string(output))
}
