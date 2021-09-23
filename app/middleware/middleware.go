package middleware

import (
	stormemb "github.com/bagus2x/tjiwi/pkg/storagemember"
	"github.com/bagus2x/tjiwi/pkg/user"
)

type Middleware struct {
	userService     user.Service
	storMembService stormemb.Service
}

func New(userService user.Service, storMembService stormemb.Service) *Middleware {
	return &Middleware{
		userService:     userService,
		storMembService: storMembService,
	}
}
