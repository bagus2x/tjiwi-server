package basepaper

import (
	"github.com/bagus2x/tjiwi/app"
	"github.com/go-playground/validator/v10"
)

type Params struct {
	StorageID      *int64  `form:"storage_id"`
	Gsm            *int64  `form:"gsm"`
	Width          *int64  `form:"width"`
	Io             *int64  `form:"io"`
	MaterialNumber *int64  `form:"material"`
	Location       *string `form:"location"`
	NextCursor     *int64  `form:"cursor"`
	Limit          *int64  `form:"limit"`
	Direction      *string `form:"dir"`
}

type Cursor struct {
	Next     int64 `json:"next"`
	Previous int64 `json:"previous"`
}

type AddBasePaperRequest struct {
	StorageID      int64 `json:"storageID" validate:"required"`
	Gsm            int64 `json:"gsm" validate:"required"`
	Width          int64 `json:"width" validate:"required"`
	Io             int64 `json:"io" validate:"required"`
	MaterialNumber int64 `json:"materialNumber" validate:"required"`
	Quantity       int64 `json:"quantity" validate:"required"`
}

func (r *AddBasePaperRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	return app.ValidateAndTranslate(validate, err)
}

type AddBasePaperResponse struct {
	ID             int64 `json:"id"`
	StorageID      int64 `json:"storageID"`
	Gsm            int64 `json:"gsm"`
	Width          int64 `json:"width"`
	Io             int64 `json:"io"`
	MaterialNumber int64 `json:"materialNumber"`
	Quantity       int64 `json:"quantity"`
	CreatedAt      int64 `json:"createdAt"`
	UpdatedAt      int64 `json:"updatedAt"`
}

type GetBasePaperResponse struct {
	ID             int64  `json:"id"`
	StorageID      int64  `json:"storageID"`
	Gsm            int64  `json:"gsm"`
	Width          int64  `json:"width"`
	Io             int64  `json:"io"`
	MaterialNumber int64  `json:"materialNumber"`
	Location       string `json:"location,omitempty"`
	Quantity       int64  `json:"quantity"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
}

type GetBasePapersResponse struct {
	Cursor     Cursor                  `json:"cursor"`
	BasePapers []*GetBasePaperResponse `json:"basePapers"`
}

type MoveToStorageRequest struct {
	ID       int64  `json:"id"`
	Location string `json:"location"`
	Quantity int64  `json:"quantity"`
}

type MoveToStorageResponse struct {
	ID             int64  `json:"id"`
	StorageID      int64  `json:"storageID"`
	Gsm            int64  `json:"gsm"`
	Width          int64  `json:"width"`
	Io             int64  `json:"io"`
	MaterialNumber int64  `json:"materialNumber"`
	Quantity       int64  `json:"quantity"`
	Location       string `json:"location"`
	CreatedAt      int64  `json:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt"`
}

type DeliverBasePaperRequest struct {
	ID       int64 `json:"id"`
	Quantity int64 `json:"quantity"`
	MemberID int64 `json:"memberID"`
}

type DeliverBasePaperResponse struct {
	ID       int64 `json:"id"`
	Quantity int64 `json:"quantity"`
	MemberID int64 `json:"memberID"`
}
