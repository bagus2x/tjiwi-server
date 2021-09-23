package storage

import (
	"github.com/bagus2x/tjiwi/app"
	"github.com/go-playground/validator/v10"
)

type Supervisor struct {
	ID       int64  `json:"id"`
	Photo    string `json:"photo,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

type CreateStorageRequest struct {
	SupervisorID int64  `json:"supervisorID" validate:"required,gte=0"`
	Name         string `json:"name" validate:"required,gte=5,lte=255"`
	Description  string `json:"description" validate:"lte=255"`
}

func (r *CreateStorageRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type CreateStorageResponse struct {
	ID           int64  `json:"id"`
	SupervisorID int64  `json:"supervisorID"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CreatedAt    int64  `json:"int64"`
	UpdatedAt    int64  `json:"updatedAt"`
}

type FindStorageResponse struct {
	ID          int64      `json:"id"`
	Supervisor  Supervisor `json:"supervisor"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   int64      `json:"createdAt"`
	UpdatedAt   int64      `json:"updatedAt"`
}

type UpdateStorageRequest struct {
	ID          int64  `json:"id" validate:"required,gte=0"`
	Name        string `json:"name" validate:"required,gte=5,lte=255"`
	Description string `json:"description" validate:"lte=255"`
}

func (r *UpdateStorageRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type UpdateStorageResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updatedAt"`
}
