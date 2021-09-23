package service

import (
	"context"
	"strings"
	"time"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/bagus2x/tjiwi/pkg/user"
	"github.com/golang-jwt/jwt"
)

type service struct {
	userRepo user.Repository
	cfg      *config.Config
}

func New(userRepo user.Repository, cfg *config.Config) user.Service {
	return &service{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *service) SignIn(ctx context.Context, req *user.SignInRequest) (*user.SignInResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	p, err := s.userRepo.FindByUsernameOrEmail(ctx, req.UsernameOrEmail, req.UsernameOrEmail)
	if app.ErrorCode(err) == app.ENotFound || p.IsDeleted {
		return nil, app.NewError(nil, app.ENotFound, "User not found")
	} else if err != nil {
		return nil, err
	}

	isMatch := p.ComparePasswords(req.Password)
	if !isMatch {
		return nil, app.NewError(nil, app.EBadRequest, "Password does not match")
	}

	accessToken, err := s.createAccessToken(p.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.CreateRefreshToken(p.ID)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateToken(ctx, p.ID, refreshToken)
	if err != nil {
		return nil, err
	}

	res := &user.SignInResponse{
		Token: user.Token{
			Access:  accessToken,
			Refresh: refreshToken,
		},
		User: user.User{
			ID:        p.ID,
			Photo:     p.Photo.String,
			Username:  p.Username,
			Email:     p.Email,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
	}

	return res, nil
}

func (s *service) SignUp(ctx context.Context, req *user.SignUpRequest) (*user.SignUpResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	p, err := s.userRepo.FindByUsernameOrEmail(ctx, req.Username, req.Email)
	if err != nil && app.ErrorCode(err) != app.ENotFound {
		return nil, err
	}
	if p != nil {
		if p.Username == req.Username {
			return nil, app.NewError(nil, app.Econflict, "Username already exist")
		}
		if p.Email == req.Email {
			return nil, app.NewError(nil, app.Econflict, "Email already exist")
		}
	}

	p = &model.User{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err = p.HashPassword()
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.createAccessToken(p.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.CreateRefreshToken(p.ID)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateToken(ctx, p.ID, refreshToken)
	if err != nil {
		return nil, app.NewError(err, app.EInternal)
	}

	res := &user.SignUpResponse{
		Token: user.Token{
			Access:  accessToken,
			Refresh: refreshToken,
		},
		User: user.User{
			ID:        p.ID,
			Photo:     p.Photo.String,
			Username:  p.Username,
			Email:     p.Email,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
	}

	return res, nil
}

func (s *service) GetUserByID(ctx context.Context, userID int64) (*user.GetUserResponse, error) {
	p, err := s.userRepo.FindByID(ctx, userID)
	if app.ErrorCode(err) == app.ENotFound {
		return nil, app.NewError(err, app.ENotFound, "User does not exist")
	} else if err != nil {
		return nil, err
	}

	res := user.GetUserResponse{
		ID:        userID,
		Photo:     p.Photo.String,
		Username:  p.Username,
		Email:     p.Email,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.CreatedAt,
	}

	return &res, nil
}

func (s *service) SearchByUsername(ctx context.Context, username string) (user.SearchUsernameResponse, error) {
	u, err := s.userRepo.MatchByUsername(ctx, username)
	if app.ErrorCode(err) == app.ENotFound {
		return nil, app.NewError(err, app.ENotFound, "User does not exist")
	} else if err != nil {
		return nil, err
	}

	var usernames = make(user.SearchUsernameResponse, 0)

	for _, usr := range u {
		usernames = append(usernames, user.GetUserResponse{
			Photo:    usr.Photo.String,
			ID:       usr.ID,
			Username: usr.Username,
		})
	}

	return usernames, err
}

func (s *service) RefreshToken(ctx context.Context, req *user.RefreshTokenRequest) (*user.RefreshTokenResponse, error) {
	claims, err := s.ExtractRefreshToken(req.Token)
	if err != nil {
		return nil, err
	}

	p, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, app.NewError(err, app.EInternal)
	}

	if p.Token != req.Token {
		return nil, app.NewError(nil, app.EInvalidRefreshToken)
	}

	accessToken, err := s.createAccessToken(claims.UserID)
	if err != nil {
		return nil, err
	}

	refereshToken, err := s.CreateRefreshToken(claims.UserID)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateToken(ctx, claims.UserID, refereshToken)
	if err != nil {
		return nil, app.NewError(err, app.EInternal)
	}

	res := user.RefreshTokenResponse{
		ID: claims.UserID,
		Token: user.Token{
			Access:  accessToken,
			Refresh: refereshToken,
		},
	}

	return &res, nil
}

func (s *service) Update(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	p := model.User{
		ID:        req.ID,
		Photo:     db.NewNullString(req.Photo, req.Photo != ""),
		Username:  req.Username,
		Email:     req.Email,
		UpdatedAt: time.Now().Unix(),
	}

	err = s.userRepo.Update(ctx, &p)
	if app.ErrorCode(err) == app.ENotFound {
		return nil, app.NewError(err, "User does not exist")
	} else if err != nil {
		return nil, err
	}

	res := user.UpdateUserResponse{
		ID:        p.ID,
		Photo:     p.Photo.String,
		Username:  p.Username,
		Email:     p.Email,
		UpdatedAt: time.Now().Unix(),
	}

	return &res, nil
}

func (s *service) Delete(ctx context.Context, userID int64) error {
	err := s.userRepo.SoftDelete(ctx, userID, true)
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(err, app.ENotFound, "User does not exist")
	}

	return err
}

func (s *service) SignOut(ctx context.Context, userID int64) error {
	err := s.userRepo.UpdateToken(ctx, userID, "")
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(err, app.ENotFound, "User does not exist")
	}

	return err
}

func (s *service) createAccessToken(userID int64) (string, error) {
	claims := user.AccessClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(s.cfg.AccessTokenLifetime())).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.AccessTokenKey()))
}

func (s *service) ExtractAccessToken(tokenStr string) (*user.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &user.AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, app.NewError(nil, app.EInternal, "Unexpected signing method")
		}
		return []byte(s.cfg.AccessTokenKey()), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, app.NewError(err, app.EAccessTokenExpired, "Access token has expired")
		}
		return nil, app.NewError(err, app.EInvalidAccessToken, "Invalid access token")
	}

	if claims, ok := token.Claims.(*user.AccessClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, app.NewError(err, app.EInvalidAccessToken, "Invalid access token")
}

func (s *service) CreateRefreshToken(userID int64) (string, error) {
	claims := user.AccessClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(s.cfg.RefreshTokenLifetime())).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.RefreshTokenKey()))
}

func (s *service) ExtractRefreshToken(tokenStr string) (*user.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &user.RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, app.NewError(nil, app.EInternal, "Unexpected signing method")
		}
		return []byte(s.cfg.RefreshTokenKey()), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, app.NewError(err, app.ERefreshTokenExpired, "Refresh token has expired")
		}
		return nil, app.NewError(err, app.EInvalidRefreshToken, "Invalid refresh token")
	}

	if claims, ok := token.Claims.(*user.RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, app.NewError(err, app.EInvalidRefreshToken, "Invalid refresh token")
}
