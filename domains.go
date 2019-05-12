package zeit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type User struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	CustomerId string `json:"customerId"`
}

type Aliases struct {
	Id      string
	Alias   string
	Created *Time
}

type Certs struct {
	Id      string
	Cns     []string
	Created *Time
}

type Domain struct {
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	ServiceType         string    `json:"serviceType"`
	NsVerifiedAt        *Time     `json:"nsVerifiedAt"`
	TxtVerifiedAt       *Time     `json:"txtVerifiedAt"`
	CdnEnabled          bool      `json:"cdnEnabled"`
	CreatedAt           *Time     `json:"createdAt"`
	ExpiresAt           *Time     `json:"expiresAt"`
	BoughtAt            *Time     `json:"boughtAt"`
	VerifiedRecord      string    `string:"verifiedRecord"`
	Verified            bool      `json:"verified"`
	Nameservers         []string  `json:"nameservers"`
	IntendedNameservers []string  `json:"intendedNameservers"`
	Creator             User      `json:"creator"`
	Suffix              bool      `json:"suffix,omitempty"`
	Aliases             []Aliases `json:"aliases,omitempty"`
	Certs               []Certs   `json:"certs,omitempty"`
}

type GetError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type VerificationError struct {
	GetError
	NsVerification struct {
		Name                string   `json:"name"`
		Nameservers         []string `json:"nameservers"`
		IntendedNameservers []string `json:"intendedNameservers"`
	} `json:"nsVerification"`
	TxtVerification struct {
		Name               string   `json:"name"`
		Values             []string `json:"values"`
		VerificationRecord string   `json:"verificationRecord"`
	} `json:"txtVerification"`
}

func (c Client) GetAllDomains() ([]Domain, error) {
	resp, err := c.makeAndDoRequest(http.MethodGet, "v4/domains", nil)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	domains := struct {
		Domains []Domain `json:"domains"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&domains)
	if err != nil {
		return nil, err
	}

	return domains.Domains, nil
}

func (c Client) AddDomain(name string) (*Domain, error) {
	parameters := struct {
		Name string `json:"name"`
	}{name}
	body, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	resp, err := c.makeAndDoRequest(http.MethodPost, "v4/domains", bytes.NewBuffer(body))

	defer closeResponseBody(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	domain := Domain{}
	err = json.NewDecoder(resp.Body).Decode(&domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (c Client) TransferInDomain(method, name, authCode string, expectedPrice int) (*Domain, error) {
	parameters := struct {
		Method        string `json:"method"`
		Name          string `json:"name"`
		AuthCode      string `json:"authCode"`
		ExpectedPrice int    `json:"expectedPrice"`
	}{method, name, authCode, expectedPrice}
	body, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	resp, err := c.makeAndDoRequest(http.MethodPost, "v4/domains", bytes.NewBuffer(body))
	defer closeResponseBody(resp)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	domain := Domain{}
	err = json.NewDecoder(resp.Body).Decode(&domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (c Client) VerifyDomain(name string) (*Domain, *VerificationError, error) {
	endpoint := fmt.Sprintf("v4/domains/%s/verify", name)
	resp, err := c.makeAndDoRequest(http.MethodPost, endpoint, nil)
	defer closeResponseBody(resp)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New(resp.Status)
	}
	verification := struct {
		Domain *Domain
		Error  *VerificationError
	}{nil, nil}
	err = json.NewDecoder(resp.Body).Decode(&verification)
	if err != nil {
		return nil, nil, err
	}
	return verification.Domain, verification.Error, nil
}

func (c Client) GetDomain(name string) (*Domain, *GetError, error) {
	endpoint := fmt.Sprintf("v4/domains/%s", name)
	resp, err := c.makeAndDoRequest(http.MethodGet, endpoint, nil)
	defer closeResponseBody(resp)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New(resp.Status)
	}
	verification := struct {
		Domain *Domain
		Error  *GetError
	}{nil, nil}
	err = json.NewDecoder(resp.Body).Decode(&verification)
	if err != nil {
		return nil, nil, err
	}
	return verification.Domain, verification.Error, nil
}

func (c Client) RemoveDomain(name string) (string, error) {
	endpoint := fmt.Sprintf("v4/domains/%s", name)
	resp, err := c.makeAndDoRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	response := struct {
		Uid string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	return response.Uid, nil
}

func (c Client) CheckDomainAvailability(name string) (bool, error) {
	endpoint := fmt.Sprintf("v4/domains/status?%s", name)
	resp, err := c.makeAndDoRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}
	response := struct {
		Available bool
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}
	return response.Available, nil
}

func (c Client) CheckDomainPrice(name string) (int, int, error) {
	endpoint := fmt.Sprintf("v4/domains/price?%s", name)
	resp, err := c.makeAndDoRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, 0, errors.New(resp.Status)
	}
	response := struct {
		Price  int
		Period int
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, 0, err
	}
	return response.Price, response.Period, nil
}

func (c Client) BuyDomain(name string, expectedPrice int) error {
	parameters := struct {
		Name          string `json:"name"`
		ExpectedPrice int    `json:"expectedPrice"`
	}{name, expectedPrice}
	body, err := json.Marshal(parameters)
	if err != nil {
		return err
	}
	resp, err := c.makeAndDoRequest(http.MethodPost, "v4/domains/buy", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	response := struct {
		Price  int
		Period int
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}
