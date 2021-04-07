package comerr

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
)

var (
	NilErr = ComErr("nil error")
)

type ComErr string

func (this *ComErr) Mark() {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknow"
		line = 0
	}
	*this = ComErr(fmt.Sprintf("\n[%s]\n%s : %d\n[err_mark:]%s", time.Now().Format("2006-01-02 15:04:05"), file, line, *this))
}

func (this *ComErr) Show() {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknow"
		line = 0
	}
	fmt.Printf("%c[1;0;35m###\n[%s]\n%s : %d\n[err_show:]\n<<%s\n>>\n###\n%c[0m", 0x1B, time.Now().Format("2006-01-02 15:04:05"), file, line, *this, 0x1B)
}

func ToComErr(errStr string) ComErr {
	if errStr == "" {
		return NilErr
	}

	return ComErr(errStr)
}

func ContextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return fmt.Errorf("process canceled")
	case context.DeadlineExceeded:
		return fmt.Errorf("deadline is exceeded")
	default:
		return nil
	}
}

func LogError(err error) error {
	if err != nil {
		log.Println(err)
	}

	return err
}
