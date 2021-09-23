package model

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Photo     sql.NullString
	Username  string
	Email     string
	Password  string
	IsDeleted bool
	Token     string
	CreatedAt int64
	UpdatedAt int64
}

func (p *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Password = string(bytes)

	return nil
}

func (p *User) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	return err == nil
}
