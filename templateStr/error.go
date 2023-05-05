package templateStr

import (
	"errors"
	"fmt"
)

type TError struct {
	Name string
	Err error
}

func (err TError) Error() string {
	return fmt.Sprintf("%s: %v", err.Name, err.Err)
}

func NotFoundFunctionError(err string) TError {
	return TError{
		Name: "NotFoundFunctionError",
		Err: errors.New(err),
	}
}

func NotFoundVariableError(err string) TError {
	return TError{
		Name: "NotFoundVariableError",
		Err: errors.New(err),
	}
}

func NotAArrayError(err string) TError {
	return TError{
		Name: "NotAArrayError",
		Err: errors.New(err),
	}
}

func IndexError(err string) TError {
	return TError{
		Name: "IndexError",
		Err: errors.New(err),
	}
}

func BadComparatorError(err string) TError {
	return TError{
		Name: "BadComparatorError",
		Err: errors.New(err),
	}
}