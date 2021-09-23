package user

import "context"

type Service interface {
	SignIn(ctx context.Context, req *SignInRequest) (*SignInResponse, error)
	SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error)
	SignOut(ctx context.Context, userID int64) error
	GetUserByID(ctx context.Context, userID int64) (*GetUserResponse, error)
	SearchByUsername(ctx context.Context, username string) (SearchUsernameResponse, error)
	Update(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error)
	Delete(ctx context.Context, userID int64) error
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
	ExtractAccessToken(tokenStr string) (*AccessClaims, error)
}
