package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/cnicolov/terraform-provider-spotinstadmin/client"
	"github.com/cnicolov/terraform-provider-spotinstadmin/client/common"
)

const (
	createAccoountRequestPath = "/setup/account"
	accountServiceBaseURL     = "https://api.spotinst.io"
)

// Service is a client for creating accounts
type Service struct {
	httpClient *client.Client
}

// New creates new accounts service client
func New(token string) *Service {
	log.Println("Initializing accounts service")
	return &Service{
		httpClient: client.New(accountServiceBaseURL, token),
	}
}

// Account represesnts Spotinst account in API
type Account struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	OrganizationID     string `json:"organizationId"`
	ProviderExternalID string `json:"providerExternalId,omitempty"`
}

// AccountNotFoundError is raised when looking up account
// fails because there's no such account in Spotinst:w
type AccountNotFoundError struct {
	AccountID string
}

func (a *AccountNotFoundError) Error() string {
	return fmt.Sprintf("Account %s not found", a.AccountID)
}

// Create creates accoount in Spotinst
func (as *Service) Create(name, iamRole, externalID string) (*Account, error) {

	body := map[string]map[string]string{
		"account": {"name": name},
	}

	log.Printf("Making request %v\n", body)
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

	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 5)

	err = as.setupCloudCredentials(account.ID, iamRole, externalID)

	if err != nil {
		_ = as.Delete(account.ID)
		return nil, err
	}

	return &account, nil
}

func (as *Service) setupCloudCredentials(accountID, iamRole, externalID string) error {

	body := map[string]map[string]string{
		"credentials": {"iamRole": iamRole, "externalId": externalID},
	}

	req, err := as.httpClient.NewRequest(http.MethodPost, "/setup/credentials/aws", &body)

	q, _ := url.ParseQuery(req.URL.RawQuery)

	q.Add("accountId", accountID)

	req.URL.RawQuery = q.Encode()
	if err != nil {
		return err
	}

	log.Printf("%#v", req)

	var r common.Response

	_, err = as.httpClient.Do(req, &r)

	if len(r.Response.Errors) > 0 {
		return fmt.Errorf("failed setting up cloud credentials, %#v", r.Response.Errors)
	}

	log.Printf("%#v", r)

	return err
}

// Get returns account by id
func (as *Service) Get(id string) (*Account, error) {
	log.Printf("Getting account %v\n", id)

	req, err := as.httpClient.NewRequest(http.MethodGet, "/setup/account", nil)

	if err != nil {
		return nil, err
	}

	var r common.Response

	_, err = as.httpClient.Do(req, &r)
	if err != nil {
		return nil, err
	}

	accountList, err := accountsFromJSON(r)

	if err != nil {
		return nil, err
	}

	return filterAccountByID(id, accountList)
}

// Delete delets account by id
func (as *Service) Delete(id string) error {
	req, err := as.httpClient.NewRequest(http.MethodDelete, fmt.Sprintf("/setup/account/%s", id), nil)
	v := make(map[string]interface{})
	_, err = as.httpClient.Do(req, &v)
	return err
}

// IsAccountNotFoundErr checks whether errors is of type AccountNotFoundError
func IsAccountNotFoundErr(err error) bool {
	var found bool
	switch err.(type) {
	case *AccountNotFoundError:
		found = true
	default:
	}
	return found
}

func accountsFromJSON(r common.Response) ([]*Account, error) {
	accountList := make([]*Account, len(r.Response.Items))

	type accountJSON struct {
		Name               string `json:"name"`
		AccountID          string `json:"accountId"`
		OrganizationID     string `json:"organizationId"`
		ProviderExternalID string `json:"providerExternalId"`
	}

	for i, data := range r.Response.Items {

		var acc accountJSON
		err := json.Unmarshal(data, &acc)
		if err != nil {
			return accountList, err
		}

		accountList[i] = &Account{
			ID:                 acc.AccountID,
			Name:               acc.Name,
			OrganizationID:     acc.OrganizationID,
			ProviderExternalID: acc.ProviderExternalID,
		}

	}
	return accountList, nil
}

func filterAccountByID(id string, accountList []*Account) (*Account, error) {
	for _, a := range accountList {
		log.Printf("Checking %v with %v\n", a.ID, id)
		if a.ID == id {
			return a, nil
		}
	}
	return nil, &AccountNotFoundError{AccountID: id}
}
