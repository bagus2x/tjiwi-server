package app

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	EInternal            = "internal_error"
	ENotFound            = "not_found"
	EBadRequest          = "bad_request"
	EUnauthorized        = "unauthorized"
	Econflict            = "conflict"
	EAccessTokenExpired  = "access_token_expired"
	ERefreshTokenExpired = "refresh_token_expired"
	EInvalidAccessToken  = "invalid_access_token"
	EInvalidRefreshToken = "invalid_refresh_token"
	EForbidden           = "forbidden"
)

type Error struct {
	Code     string
	Messages []string
	Err      error
}

func (e *Error) Error() string {
	var buff bytes.Buffer

	if e.Err != nil {
		buff.WriteString(e.Err.Error())
	} else if len(e.Messages) != 0 {
		fmt.Fprintf(&buff, "<%s> [%s]", e.Code, strings.Join(e.Messages, ", "))
	} else {
		fmt.Fprintf(&buff, "<%s>", e.Code)
	}

	return buff.String()
}

func NewError(err error, code string, msg ...string) error {
	return &Error{
		Code:     code,
		Messages: msg,
		Err:      err,
	}
}

func ErrorCode(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}

	return EInternal
}

func ErrorMessage(err error) []string {
	if err == nil {
		return []string{}
	} else if e, ok := err.(*Error); ok {
		if len(e.Messages) == 0 {
			return []string{strings.Title(strings.Replace(e.Code, "_", " ", -1))}
		}
		return e.Messages
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}

	return []string{"An internal error has occurred"}
}
