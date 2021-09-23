package model

type StorageMember struct {
	ID        int64
	Storage   Storage
	Member    User
	IsAdmin   bool
	IsActive  bool
	IsDeleted bool
	CreatedAt int64
	UpdatedAt int64
}
