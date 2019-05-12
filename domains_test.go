package zeit

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/kochie/zeit-api-go/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var TestToken = os.Getenv("NOW_TEST_TOKEN")
const rootUrl = "https://zeit.api.co"

func TestClient_GetAllDomains(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testDomains := []Domain{{},{}}
	response, err := json.Marshal(&struct {
		Domains []Domain `json:"domains"`
	}{testDomains})
	a.Nil(err)

	header := http.Header{}
	header.Set("X-RateLimit-remaining","10")
	header.Set("X-RateLimit-limit","10")
	header.Set("X-RateLimit-reset","0")

	httpResponse := http.Response{
		Header: header,
		Body: ioutil.NopCloser(bytes.NewBuffer(response)),
		StatusCode: http.StatusOK,
	}

	mockHttpClient := mocks.NewMockHttpClient(ctrl)
	mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

	client := Client{
		TestToken,
		rootUrl,
		mockHttpClient,
		&rateLimit{},
		"",
	}
	domains, err := client.GetAllDomains()

	a.Nil(err, "Error should be nil")
	a.NotNil(domains, "Domain list should be defined")
	a.Equal(len(testDomains), len(domains), "Should be the same length")
}

func TestClient_AddDomain(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testDomains := []Domain{{},{}}
	response, err := json.Marshal(&struct {
		Domains []Domain `json:"domains"`
	}{testDomains})
	a.Nil(err)

	header := http.Header{}
	header.Set("X-RateLimit-remaining","10")
	header.Set("X-RateLimit-limit","10")
	header.Set("X-RateLimit-reset","0")

	httpResponse := http.Response{
		Header: header,
		Body: ioutil.NopCloser(bytes.NewBuffer(response)),
		StatusCode: http.StatusOK,
	}

	mockHttpClient := mocks.NewMockHttpClient(ctrl)
	mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

	client := Client{
		TestToken,
		rootUrl,
		mockHttpClient,
		&rateLimit{},
		"",
	}

	domainNames := []string{
		"test.com",
		"www.test.com",
	}
	for _, domainName := range domainNames {
		t.Run(domainName, func(t *testing.T) {
			domains, err := client.AddDomain(domainName)
			a.Nil(err, "Error should be nil")
			a.NotNil(domains, "Domain list should be defined")
			//a.Equal(len(testDomains), len(domains), "Should be the same length")
		})
	}

}

func TestClient_TransferInDomain(t *testing.T) {

}

func TestClient_VerifyDomain(t *testing.T) {

}

func TestClient_GetDomain(t *testing.T) {

}

func TestClient_RemoveDomain(t *testing.T) {

}

func TestClient_CheckDomainAvailability(t *testing.T) {

}

func TestClient_CheckDomainPrice(t *testing.T) {

}

func TestClient_BuyDomain(t *testing.T) {

}