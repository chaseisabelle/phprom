package funcs

import (
	"errors"
	"fmt"
	"github.com/chaseisabelle/goresp"
)

func Err(msg string, etc ...interface{}) error {
	if len(etc) == 0 {
		return errors.New(msg)
	}

	return fmt.Errorf(msg, etc...)
}

func RErr(msg string, etc ...interface{}) *goresp.Error {
	return goresp.NewError(Err(msg, etc...))
}

func RErrs(msg string, etc ...interface{}) []goresp.Value {
	return []goresp.Value{RErr(msg, etc...)}
}

func RRErrs(msg string, etc ...interface{}) ([]goresp.Value, error) {
	return RErrs(msg, etc...), Err(msg, etc...)
}

func ERErrs(err error) ([]goresp.Value, error) {
	return []goresp.Value{goresp.NewError(err)}, err
}
