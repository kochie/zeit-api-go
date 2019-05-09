package zeit_api_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	token := "AABBCC"
	assert.NotNil(t, NewClient(token), "Client was not created properly")
}
