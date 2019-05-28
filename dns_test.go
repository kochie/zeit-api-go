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

	domainNames := map[string]*Record{
		"test.com": {
			Name:  "",
			Type:  RecordTypeCNAME,
			Value: "test.com",
			Id:    "123456",
		},
		"www.test.com": {
			Name:  "",
			Type:  RecordTypeCNAME,
			Value: "test.com",
			Id:    "123456",
		},
		"test.com.au": {
			Name:  "foo",
			Type:  RecordTypeCNAME,
			Value: "test.com",
			Id:    "123456",
		},
	}

	for domainName, record := range domainNames {
		response, err := json.Marshal(&struct {
			Uid string `json:"uid"`
		}{record.Id})
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
			uid, conflictError, err := client.CreateDNSRecord(domainName, record)
			a.Nil(err, "Error should be nil")
			a.Nil(conflictError, "Should not have a conflict")
			a.NotNil(uid, "uid should be defined")
			a.Equal(record.Id, uid, "uid should be the same")
		})
	}

	incorrectDomainRecords := []struct {
		origin      string
		record      *Record
		errorString string
		mock        bool
	}{
		{"test.com", &Record{
			Name:  "@",
			Type:  RecordTypeCNAME,
			Value: "test.com",
			Id:    "123456",
		}, ErrorOrigin, false},
		{"www.test.com", &Record{
			Name:  "@",
			Type:  RecordTypeCNAME,
			Value: "test.com",
			Id:    "123456",
		}, ErrorOrigin, false},
		{"test.com", nil, ErrorNilRecord, false},
		{"www.test.com", nil, ErrorNilRecord, false},
	}

	for _, incorrectDomain := range incorrectDomainRecords {
		if incorrectDomain.mock {
			response, err := json.Marshal(&struct {
				Uid string `json:"uid"`
			}{incorrectDomain.record.Id})
			a.Nil(err)

			httpResponse := makeResponse(response)
			mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)
		}

		client := Client{
			TestToken,
			rootUrl,
			mockHttpClient,
			&rateLimit{},
			"",
		}

		t.Run(incorrectDomain.origin, func(t *testing.T) {
			uid, conflictError, err := client.CreateDNSRecord(incorrectDomain.origin, incorrectDomain.record)
			a.Error(err, incorrectDomain.errorString)
			a.Nil(conflictError, "Should not have a conflict")
			a.Empty(uid, "uid should be empty")
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
