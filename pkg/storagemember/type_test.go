package storagemember

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAja(t *testing.T) {
	req := UpdateStorMembRequest{
		ID:       0,
		IsAdmin:  false,
		IsActive: false,
	}
	assert.NoError(t, req.Validate())
}
