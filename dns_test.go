package zeit

import (
	"encoding/json"
	"errors"
	"fmt"
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
			uid, err := client.CreateDNSRecord(domainName, record)
			a.Nil(err, "Error should be nil")
			a.NotNil(uid, "uid should be defined")
			a.Equal(record.Id, uid, "uid should be the same")
		})
	}

	badRequests := []struct {
		domain     string
		response   error
		statusCode int
		record     *Record
	}{
		{
			"www.test.com",
			BasicError{},
			http.StatusBadRequest,
			&Record{},
		},
		{
			"foo.test.com",
			ConflictError{},
			http.StatusConflict,
			&Record{},
		},
		{
			"foo.test.com",
			BasicError{
				"invalid_value",
				"Invalid record value: \"0 issue letsencrypt.org\"",
			},
			http.StatusBadRequest,
			&Record{},
		},
		{"test.com", errors.New(ErrorOrigin), -1, &Record{
			Name:  "@",
			Type:  RecordTypeCNAME,
			Value: "test.com",
			Id:    "123456",
		}},
		{"www.test.com", errors.New(ErrorOrigin), -1, &Record{
			Name:  "@",
			Type:  RecordTypeCNAME,
			Value: "test.com",
			Id:    "123456",
		}},
		{"test.com", errors.New(ErrorNilRecord), -1, nil},
		{"www.test.com", errors.New(ErrorNilRecord), -1, nil},
	}

	for _, badRequest := range badRequests {
		t.Run(fmt.Sprintf("%s_%d", badRequest.domain, badRequest.statusCode), func(t *testing.T) {
			if badRequest.statusCode > 0 {
				response, err := json.Marshal(&struct {
					Error error
				}{badRequest.response})
				a.Nil(err, "should create response")

				httpResponse := makeResponse(response, badRequest.statusCode)
				mockHttpClient.EXPECT().Do(gomock.Any()).Return(&httpResponse, nil)
			}

			client := Client{
				TestToken,
				rootUrl,
				mockHttpClient,
				&rateLimit{},
				"",
			}

			uid, err := client.CreateDNSRecord(badRequest.domain, badRequest.record)
			a.Error(err, badRequest.response.Error())
			a.IsType(badRequest.response, err)
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
