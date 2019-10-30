package command

import (
	"devtools/comerr"
	"encoding/json"
	"os/exec"
	"reflect"
)

type Command struct {
	*exec.Cmd
}

func (this *Command) StdoutJson(output interface{}) error {
	if output == nil || reflect.TypeOf(output).Kind() != reflect.Ptr {
		return comerr.ParamInvalid
	}

	if stdout, err := this.StdoutPipe(); err != nil {
		return err
	} else if err = this.Start(); err != nil {
		return err
	} else if err = json.NewDecoder(stdout).Decode(output); err != nil {
		return err
	} else {
		return this.Wait()
	}
}
