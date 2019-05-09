package zeit_api_go

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

const TestToken = "QtneoQOn5LyifOuXbgG7mKAz"

func saveJsonFile(filename string, data interface{}) {
	f, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("couldn't create the file beacuse %s", err.Error()))
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(&data)
	if err != nil {
		panic(fmt.Sprintf("couldn't encode the data because %s", err.Error()))
	}
}

type HttpClientMock interface {
	Do(req *http.Request) (*http.Response, error)
}

func TestClient_GetAllDomains(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := something.NewMockMyInterface(mockCtrl)

	client := NewClient(TestToken)
	domains, err := client.GetAllDomains()

	a := assert.New(t)
	a.Nil(err, "Error should be nil")
	a.NotNil(domains, "Domain list should be defined")

	a.Equal(len(domains), 1)
}