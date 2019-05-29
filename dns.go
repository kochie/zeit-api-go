package zeit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (c Client) ListDNSRecords(domain string) ([]Record, error) {
	endpoint := fmt.Sprintf("v2/domains/%s/records", domain)

	resp, err := c.makeAndDoRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var records []Record

	err = json.NewDecoder(resp.Body).Decode(&struct {
		Records *[]Record
	}{&records})
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (c Client) CreateDNSRecord(domain string, record *Record) (string, error) {
	if record == nil {
		return "", errors.New(ErrorNilRecord)
	}

	if record.Name == "@" {
		return "", errors.New(ErrorOrigin)
	}

	parameters := struct {
		Name       string `json:"name"`
		RecordType string `json:"type"`
		Value      string `json:"value"`
	}{record.Name, record.Type, strings.TrimSuffix(record.GetValue(), ".")}

	body, err := json.Marshal(parameters)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("v2/domains/%s/records", domain)
	resp, err := c.makeAndDoRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	defer closeResponseBody(resp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusBadRequest {
		requestError := BasicError{}
		err = json.NewDecoder(resp.Body).Decode(&struct {
			Error BasicError
		}{requestError})
		if err != nil {
			return "", err
		}
		return "", requestError
	}

	if resp.StatusCode == http.StatusConflict {
		conflictError := ConflictError{}
		err = json.NewDecoder(resp.Body).Decode(&struct {
			Error ConflictError
		}{conflictError})
		if err != nil {
			return "", err
		}
		log.Println(resp.Status, record.Name, record.Type, record.GetValue())
		return "", conflictError
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	var uid string
	err = json.NewDecoder(resp.Body).Decode(&struct {
		Uid *string
	}{&uid})

	if err != nil {
		return "", err
	}

	return uid, nil
}

func (c Client) RemoveDNSRecord(domain, recId string) error {
	endpoint := fmt.Sprintf("v2/domains/%s/records/%s", domain, recId)
	resp, err := c.makeAndDoRequest(http.MethodDelete, endpoint, nil)

	defer closeResponseBody(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}
