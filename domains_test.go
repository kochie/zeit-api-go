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

func makeResponse(response []byte, statusCode int) http.Response {
	header := http.Header{}
	header.Set("X-RateLimit-Remaining", "10")
	header.Set("X-RateLimit-Limit", "10")
	header.Set("X-RateLimit-Reset", "0")
	httpResponse := http.Response{
		Header:     header,
		Body:       ioutil.NopCloser(bytes.NewBuffer(response)),
		StatusCode: statusCode,
	}
	return httpResponse
}

func TestClient_GetAllDomains(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testDomains := []Domain{{}, {}}
	response, err := json.Marshal(&struct {
		Domains []Domain `json:"domains"`
	}{testDomains})
	a.Nil(err)

	httpResponse := makeResponse(response, http.StatusOK)

	mockHttpClient := mocks.NewMockHttpClient(ctrl)
	mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

	client := Client{
		TestToken,
		rootUrl,
		mockHttpClient,
		&rateLimit{},
		"",
	}
	domains, err := client.ListAllDomains()

	a.Nil(err, "Error should be nil")
	a.NotNil(domains, "Domain list should be defined")
	a.Equal(len(testDomains), len(domains), "Should be the same length")
}

func TestClient_AddDomain(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testDomains := []Domain{{}, {}}
	response, err := json.Marshal(&struct {
		Domains []Domain `json:"domains"`
	}{testDomains})
	a.Nil(err)

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := []string{
		"test.com",
		"www.test.com",
	}

	for _, domainName := range domainNames {
		httpResponse := makeResponse(response, http.StatusOK)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			domains, err := client.AddDomain(domainName)
			a.Nil(err, "Error should be nil")
			a.NotNil(domains, "Domain list should be defined")
		})
	}
}

func TestClient_TransferInDomain(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testDomains := []Domain{{}, {}}
	response, err := json.Marshal(&struct {
		Domains []Domain `json:"domains"`
	}{testDomains})
	a.Nil(err)

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := []string{
		"test.com",
		"www.test.com",
	}

	for _, domainName := range domainNames {
		httpResponse := makeResponse(response, http.StatusOK)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			domains, err := client.TransferInDomain(domainName, "fakeAuthCode", 10)
			a.Nil(err, "Error should be nil")
			a.NotNil(domains, "Domain list should be defined")
		})
	}
}

func TestClient_VerifyDomain(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := map[string]Domain{
		"test.com":     {},
		"www.test.com": {},
	}

	for domainName, domain := range domainNames {
		response, err := json.Marshal(&struct {
			Domain Domain `json:"domain"`
		}{domain})
		a.Nil(err)

		httpResponse := makeResponse(response, http.StatusOK)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			domain, err := client.VerifyDomain(domainName)
			a.Nil(err, "Error should be nil")
			a.NotNil(domain, "Domain list should be defined")
		})
	}
}

func TestClient_GetDomain(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := map[string]Domain{
		"test.com":     {},
		"www.test.com": {},
	}

	for domainName, testDomain := range domainNames {
		response, err := json.Marshal(&struct {
			Domain Domain `json:"domain"`
		}{testDomain})
		a.Nil(err)

		httpResponse := makeResponse(response, http.StatusOK)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			domain, err := client.GetDomain(domainName)
			a.Nil(err, "Error should be nil")
			a.NotNil(domain, "Domain list should be defined")
			a.Equal(testDomain, *domain, "Domain should be equal")
		})
	}
}

func TestClient_RemoveDomain(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testDomains := []Domain{{}, {}}
	response, err := json.Marshal(&struct {
		Domains []Domain `json:"domains"`
	}{testDomains})
	a.Nil(err)

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := []string{
		"test.com",
		"www.test.com",
	}

	for _, domainName := range domainNames {
		httpResponse := makeResponse(response, http.StatusOK)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			domainsUid, err := client.RemoveDomain(domainName)
			a.Nil(err, "Error should be nil")
			a.NotNil(domainsUid, "Domain list should be defined")
		})
	}
}

func TestClient_CheckDomainAvailability(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := map[string]bool{
		"test.com":     true,
		"www.test.com": false,
	}

	for domainName, available := range domainNames {
		response, err := json.Marshal(&struct {
			Available bool `json:"available"`
		}{available})
		a.Nil(err)
		httpResponse := makeResponse(response, http.StatusOK)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			domainAvailability, err := client.CheckDomainAvailability(domainName)
			a.Nil(err, "Error should be nil")
			a.Equal(available, domainAvailability, "Domain list should be defined")
		})
	}
}

func TestClient_CheckDomainPrice(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := map[string][2]int{
		"test.com":     {1, 2},
		"www.test.com": {1, 2},
	}

	for domainName, domainValue := range domainNames {
		response, err := json.Marshal(&struct {
			Price  int `json:"price"`
			Period int `json:"period"`
		}{domainValue[0], domainValue[1]})
		a.Nil(err)
		httpResponse := makeResponse(response, http.StatusOK)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			price, period, err := client.CheckDomainPrice(domainName)
			a.Nil(err, "Error should be nil")
			a.Equal(domainValue[0], price, "Price should be the same")
			a.Equal(domainValue[1], period, "Period should be the same")
		})
	}
}

func TestClient_BuyDomain(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := map[string]int{
		"test.com":     10,
		"www.test.com": 20,
	}

	for domainName, expectedPrice := range domainNames {
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusOK}, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			err := client.BuyDomain(domainName, expectedPrice)
			a.Nil(err, "Error should be nil")
		})
	}
}
