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

// GetAllDomains will return a slice of domains registered with the user.
func (c Client) ListAllDomains() ([]Domain, error) {
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

// AddDomain will add a specified domain name to ZEIT, either as an external or internal domain.
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

// TransferInDomain will initiate a domain transfer request from an external Registrar to ZEIT.
func (c Client) TransferInDomain(name, authCode string, expectedPrice int) (*Domain, error) {
	parameters := struct {
		Method        string `json:"method"`
		Name          string `json:"name"`
		AuthCode      string `json:"authCode"`
		ExpectedPrice int    `json:"expectedPrice"`
	}{"transfer-in", name, authCode, expectedPrice}
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

// VerifyDomain will check if the domain either has the correct nameservers for ZEIT defined or if the DNS TXT
// verification is set.
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

// GetDomain will return specific information for one domain.
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

// RemoveDomain will remove a domain from the ZEIT DNS server.
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

// CheckDomainAvailability will check if the specified domain is available for sale.
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

// CheckDomainPrice will check how much a domain will cost to purchase and will return the price and period of purchase.
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

// BuyDomain will buy a domain name at the expectedPrice.
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
	return nil
}
