package app

import (
	"net/http"
)

type Success struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type Failure struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code     string   `json:"code"`
	Messages []string `json:"messages"`
}

func Status(err error) int {
	switch ErrorCode(err) {
	case EBadRequest:
		return http.StatusBadRequest
	case EUnauthorized, EAccessTokenExpired, ERefreshTokenExpired, EInvalidAccessToken, EInvalidRefreshToken:
		return http.StatusUnauthorized
	case EForbidden:
		return http.StatusForbidden
	case ENotFound:
		return http.StatusNotFound
	case Econflict:
		return http.StatusConflict
	case EInternal:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}
