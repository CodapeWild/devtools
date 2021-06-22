package command

import (
	"bytes"
	"devtools/comerr"
	"devtools/file"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"reflect"
	"time"
)

type Command struct {
	*exec.Cmd
	buf *bytes.Buffer
}

func CmdFromPath(path string, args ...string) (*Command, error) {
	if !file.IsFileExists(path) {
		return nil, comerr.ErrFileNotExists
	}

	cmd := &Command{Cmd: exec.Command(path, args...)}
	cmd.buf = bytes.NewBuffer(nil)
	cmd.Stdout = cmd.buf

	return cmd, nil
}

func CmdFromCmd(cmd string, args ...string) (*Command, error) {
	if path, err := exec.LookPath(cmd); err != nil {
		return nil, err
	} else {
		return CmdFromPath(path, args...)
	}
}

func (this *Command) CombineOutErr() *Command {
	if this.buf == nil || this.Stdout == nil {
		this.buf = bytes.NewBuffer(nil)
		this.Stdout = this.buf
	}
	this.Stderr = this.Stdout

	return this
}

func (this *Command) RunWithTimeout(timeout time.Duration) ([]byte, error) {
	err := this.Start()
	if err != nil {
		return nil, err
	}

	after := time.AfterFunc(timeout, func() {
		err = this.Process.Kill()
	})

	cmderr := this.Wait()

	if !after.Stop() && err == nil {
		err = comerr.ErrProcessOvertime
	}

	output, err := ioutil.ReadAll(this.buf)
	if err != nil {
		return nil, err
	}

	return output, cmderr
}

func (this *Command) RunWithJsonOut(out interface{}) error {
	if rv := reflect.ValueOf(out); rv.Kind() != reflect.Ptr || rv.IsNil() {
		return comerr.ErrTypeInvalid
	}

	err := this.Run()
	if err != nil {
		return err
	}

	var output []byte
	if output, err = ioutil.ReadAll(this.buf); err != nil {
		return err
	}

	return json.Unmarshal(output, out)
}
