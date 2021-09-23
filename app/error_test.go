package app

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	var e error = &Error{
		Code:     "CODE1",
		Messages: []string{"ERROR1"},
		Err: &Error{
			Code:     "CODE2",
			Messages: []string{"ERROR2"},
			Err:      errors.New("WADUH!!"),
		},
	}
	_ = e

	t.Log(ErrorMessage(errors.New("wkwkkw")))
	t.Log(ErrorCode(errors.New("wkwkkw")))
}
