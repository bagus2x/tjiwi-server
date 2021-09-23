package model

type History struct {
	ID        int64
	Storage   Storage
	BasePaper BasePaper
	Member    User
	Status    string
	Affected  int64
	CreatedAt int64
}
