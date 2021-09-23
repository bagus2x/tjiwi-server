package user

import (
	"github.com/bagus2x/tjiwi/app"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type AccessClaims struct {
	jwt.StandardClaims
	UserID int64
}

type RefreshClaims struct {
	jwt.StandardClaims
	UserID int64 `json:"userID"`
}

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type User struct {
	ID        int64  `json:"id"`
	Photo     string `json:"photo"`
	Username  string `json:"username"`
	Email     string `json:"email,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
}

type SignInRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail" validate:"required,gt=0,lte=255,excludesall= "`
	Password        string `json:"password" validate:"required,gt=0,lte=255"`
}

func (r *SignInRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type SignInResponse struct {
	Token Token `json:"token"`
	User  User  `json:"user"`
}

type SignUpRequest struct {
	Username string `json:"username" validate:"required,gte=5,lte=255,excludesall= "`
	Email    string `json:"email" validate:"required,gte=5,lte=255"`
	Password string `json:"password" validate:"required,gte=5,lte=255"`
}

func (r *SignUpRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type SignUpResponse struct {
	Token Token `json:"token"`
	User  User  `json:"user"`
}

type RefreshTokenRequest struct {
	Token string `json:"token"`
}

type RefreshTokenResponse struct {
	ID    int64 `json:"id"`
	Token Token `json:"token"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"id" validate:"required,gte=0"`
	Photo    string `json:"photo"`
	Username string `json:"username" validate:"required,gte=5,lte=255"`
	Email    string `json:"email" validate:"required,gte=5,lte=255"`
}

func (r *UpdateUserRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type UpdateUserResponse struct {
	ID        int64  `json:"id"`
	Photo     string `json:"photo"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	UpdatedAt int64  `json:"updatedAt"`
}

type GetUserResponse User

type SearchUsernameRequest struct {
	Username string `json:"username" validate:"required,lte=255,excludesall= "`
}

func (r *SearchUsernameRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type SearchUsernameResponse []GetUserResponse
