package zeit

import (
	"github.com/golang/mock/gomock"
	"github.com/kochie/zeit-api-go/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	token := "AABBCC"
	assert.NotNil(t, NewClient(token), "Client was not created properly")
}

func TestRateLimit(t *testing.T) {
	token := ""
	a := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHttpClient := mocks.NewMockHttpClient(ctrl)
	rl := &rateLimit{}

	client := Client{
		token:      token,
		rootUrl:    rootUrl,
		httpClient: mockHttpClient,
		rateLimit:  rl,
	}

	a.NotNil(client)
}
