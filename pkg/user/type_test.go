package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignInrequest(t *testing.T) {
	u := SignInRequest{
		UsernameOrEmail: "bagusganteng",
		Password:        "wdwdwd",
	}
	assert.NoError(t, u.Validate())
}
