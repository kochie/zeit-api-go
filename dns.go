package zeit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

	records := struct {
		Records []Record `json:"records"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return nil, err
	}

	return records.Records, nil
}

func (c Client) CreateDNSRecord(domain, name, recordType, value string) (string, error) {
	parameters := struct {
		Name       string `json:"name"`
		RecordType string `json:"type"`
		Value      string `json:"value"`
	}{name, recordType, value}
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
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	record := struct {
		Uid string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&record)
	if err != nil {
		return "", err
	}
	return record.Uid, nil
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
