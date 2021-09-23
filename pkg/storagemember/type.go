package storagemember

import (
	"github.com/bagus2x/tjiwi/app"
	"github.com/go-playground/validator/v10"
)

type Storage struct {
	ID          int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Member struct {
	ID       int64  `json:"id"`
	Photo    string `json:"photo,omitempty"`
	Username string `json:"username,omitempty"`
}

type CreateStorMembRequest struct {
	StorageID int64 `json:"storageID" validate:"required,gte=0"`
	MemberID  int64 `json:"memberID" validate:"required,gte=0"`
	IsAdmin   bool  `json:"isAdmin" validate:"required"`
	IsActive  bool  `json:"isActive" validate:"required"`
}

func (r *CreateStorMembRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type CreateStorMembResponse struct {
	ID        int64 `json:"id"`
	StorageID int64 `json:"storageID"`
	MemberID  int64 `json:"memberID"`
	IsAdmin   bool  `json:"isAdmin"`
	IsActive  bool  `json:"isActive"`
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

type GetStorMembResponse struct {
	ID        int64   `json:"id"`
	Storage   Storage `json:"storage"`
	Member    Member  `json:"member"`
	IsAdmin   bool    `json:"isAdmin"`
	IsActive  bool    `json:"isActive"`
	CreatedAt int64   `json:"createdAt"`
	UpdatedAt int64   `json:"updatedAt"`
}

type UpdateStorMembRequest struct {
	ID       int64 `json:"id" validate:"required,gte=0"`
	IsAdmin  bool  `json:"isAdmin"`
	IsActive bool  `json:"isActive"`
}

func (r *UpdateStorMembRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type UpdateStorMembResponse struct {
	ID        int64 `json:"id"`
	IsAdmin   bool  `json:"isAdmin"`
	IsActive  bool  `json:"isActive"`
	UpdatedAt int64 `json:"updatedAt"`
}
