package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/kinvey/terraform-provider-spotinstadmin/client"
	"github.com/kinvey/terraform-provider-spotinstadmin/client/common"
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
	ExternalID         string `json:"externalId,omitempty"`
}

type ExternalID struct {
	ID         string `json:"externalId"`
	Expiration string `json:"maxValidUntil"`
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
func (as *Service) Create(name string) (*Account, error) {

	body := map[string]map[string]string{
		"account": {"name": name},
	}

	log.Printf("[TRACE] Making request %v\n", body)
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

	externalId, err := as.generateExternalId(account.ID)

	if err != nil {
		_ = as.Delete(account.ID)
		return nil, err
	}

	account.ExternalID = externalId.ID

	return &account, nil
}

func (as *Service) generateExternalId(accountID string)(*ExternalID, error) {

	log.Printf("[TRACE] Generating ExternalID\n")
	req, err := as.httpClient.NewRequest(http.MethodPost, "/setup/credentials/aws/externalId", nil)

	if err != nil {
		return nil, err
	}
	
	q, _ := url.ParseQuery(req.URL.RawQuery)

	q.Add("accountId", accountID)
	
	req.URL.RawQuery = q.Encode()
	
	var v common.Response
	_, err = as.httpClient.Do(req, &v)
	if err != nil {
		return nil, err
	}
	if len(v.Response.Items) == 0 {
		return nil, errors.New("Couldn't generate externalID")
	}

	var externalID ExternalID

	fmt.Println(string(v.Response.Items[0]))

	err = json.Unmarshal(v.Response.Items[0], &externalID)

	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 5)

	return &externalID, nil
}

func (as *Service) Link(accountID, iamRole string) error {

	body := map[string]map[string]string{
		"credentials": {"iamRole": iamRole},
	}

	req, err := as.httpClient.NewRequest(http.MethodPost, "/setup/credentials/aws", &body)

	q, _ := url.ParseQuery(req.URL.RawQuery)

	q.Add("accountId", accountID)

	req.URL.RawQuery = q.Encode()
	if err != nil {
		return err
	}

	log.Printf("[TRACE] %#v", req)

	var r common.Response

	_, err = as.httpClient.Do(req, &r)

	if len(r.Response.Errors) > 0 {
		return fmt.Errorf("failed linking account, %#v", r.Response.Errors)
	}

	if r.Response.Status.Code != 200 {
		return fmt.Errorf("Can't link accouts")
	}

	log.Printf("[TRACE] %#v", r)

	return err
}

// Get returns account by id
func (as *Service) Get(id string) (*Account, error) {
	log.Printf("[TRACE] Getting account %v\n", id)

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
		log.Printf("[TRACE] Checking %v with %v\n", a.ID, id)
		if a.ID == id {
			return a, nil
		}
	}
	return nil, &AccountNotFoundError{AccountID: id}
}
