package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cnicolov/terraform-provider-spotinstadmin/client/common"
)

const (
	createAccoountRequestPath = "/setup/account"
	accountServiceBaseURL     = "https://api.spotinst.io"
)

type AccountService struct {
	httpClient *Client
}

func NewAccountService(token string) *AccountService {
	return &AccountService{
		httpClient: NewClient(accountServiceBaseURL, token),
	}
}

type Account struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	OrganizationID     string `json:"organizationId"`
	ProviderExternalID string `json:"providerExternalId,omitempty"`
}

type AccountNotFoundError struct {
	AccountID string
}

func (a *AccountNotFoundError) Error() string {
	return fmt.Sprintf("Account %s not found", a.AccountID)
}

func (as *AccountService) Create(name string) (*Account, error) {
	body := struct {
		Account struct {
			Name string `json:"name"`
		} `json:"account"`
	}{
		struct {
			Name string `json:"name"`
		}{Name: name},
	}

	req, err := as.httpClient.NewRequest(http.MethodPost, "/setup/account", &body)
	if err != nil {
		return nil, err
	}

	var v common.Response
	_, err = as.httpClient.Do(req, &v)
	if err != nil {
		return nil, err
	}
	if len(v.Response.Items) == 0 {
		return nil, errors.New("Couldn't create account")
	}
	var account Account

	fmt.Println(string(v.Response.Items[0]))

	err = json.Unmarshal(v.Response.Items[0], &account)

	return &account, err
}

func (as *AccountService) Get(id string) (*Account, error) {
	req, err := as.httpClient.NewRequest(http.MethodGet, "/setup/account", nil)
	if err != nil {
		return nil, err
	}

	var r common.Response

	_, err = as.httpClient.Do(req, &r)
	if err != nil {
		return nil, err
	}

	var al []*Account

	type getAccountResponse struct {
		Name               string `json:"name"`
		AccountID          string `json:"accountId"`
		OrganizationID     string `json:"organizationId"`
		ProviderExternalID string `json:"providerExternalId"`
	}

	for _, data := range r.Response.Items {

		var acc getAccountResponse
		err = json.Unmarshal(data, &acc)
		if err != nil {
			return nil, err
		}

		al = append(al, &Account{
			ID:                 acc.AccountID,
			Name:               acc.Name,
			OrganizationID:     acc.OrganizationID,
			ProviderExternalID: acc.ProviderExternalID,
		})

	}

	return filterAccountByID(id, al)
}

func filterAccountByID(id string, al []*Account) (*Account, error) {
	for _, a := range al {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, &AccountNotFoundError{AccountID: id}
}

func (as *AccountService) Delete(id string) error {
	req, err := as.httpClient.NewRequest(http.MethodDelete, fmt.Sprintf("/setup/account/%s", id), nil)
	v := make(map[string]interface{})
	_, err = as.httpClient.Do(req, &v)
	return err
}
