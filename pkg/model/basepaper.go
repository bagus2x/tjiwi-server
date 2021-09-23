package model

type BasePaper struct {
	ID             int64
	Storage        Storage
	Gsm            int64
	Width          int64
	Io             int64
	MaterialNumber int64
	Quantity       int64
	Location       string
	IsDeleted      bool
	CreatedAt      int64
	UpdatedAt      int64
}
