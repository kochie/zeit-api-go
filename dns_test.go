package zeit

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/kochie/zeit-api-go/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClient_ListDNSRecords(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := []string{
		"test.com",
		"www.test.com",
	}

	for _, domainName := range domainNames {
		records := []Record{{}, {}, {}}
		response, err := json.Marshal(&struct {
			Records []Record `json:"records"`
		}{records})
		a.Nil(err)

		httpResponse := makeResponse(response)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			dnsRecords, err := client.ListDNSRecords(domainName)
			a.Nil(err, "Error should be nil")
			a.NotNil(dnsRecords, "dns record list should be defined")
			a.Equal(len(records), len(dnsRecords), "records should be the same length")
		})
	}
}

func TestClient_CreateDNSRecord(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := map[string]struct {
		Name       string
		RecordType string
		Value      string
		Uid        string
	}{
		"test.com": {
			"@",
			"CNAME",
			"test.com",
			"12345",
		},
		"www.test.com": {
			"@",
			"CNAME",
			"test.com",
			"12345",
		},
	}

	for domainName, record := range domainNames {
		response, err := json.Marshal(&struct {
			Uid string `json:"uid"`
		}{record.Uid})
		a.Nil(err)

		httpResponse := makeResponse(response)
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			uid, err := client.CreateDNSRecord(domainName, record.Name, record.RecordType, record.Value)
			a.Nil(err, "Error should be nil")
			a.NotNil(uid, "uid should be defined")
			a.Equal(record.Uid, uid, "uid should be the same")
		})
	}
}

func TestClient_RemoveDNSRecord(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := mocks.NewMockHttpClient(ctrl)

	domainNames := map[string]string{
		"test.com":     "1234556",
		"www.test.com": "545656",
	}

	for domainName, recordId := range domainNames {
		mockHttpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(nil),
		}, nil)

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(domainName, func(t *testing.T) {
			err := client.RemoveDNSRecord(domainName, recordId)
			a.Nil(err, "Error should be nil")
		})
	}
}
