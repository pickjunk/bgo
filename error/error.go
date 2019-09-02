package error

import (
	"fmt"

	"github.com/rs/zerolog"
)

// BusinessError struct
type BusinessError struct {
	Code int
	Msg  string
}

func (e BusinessError) Error() string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, e.Code, e.Msg)
}

// Throw a BusinessError with panic
// Note: make sure to use this func in the call stack
// which has registered a defer to catch the panic
func Throw(code int, msg string) {
	panic(&BusinessError{
		Code: code,
		Msg:  msg,
	})
}

// SystemError struct
type SystemError struct {
	Event *zerolog.Event
	Err   error
}

func (e SystemError) Error() string {
	return e.Err.Error()
}
