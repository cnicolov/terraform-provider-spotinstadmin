package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cnicolov/terraform-provider-spotinstadmin/client"
)

const (
	usersServiceBaseURL = "https://console.spotinst.com"
	emptyString         = ``
)

// Service ...
type Service struct {
	httpClient *client.Client
}

// New ..
func New(token string) *Service {
	return &Service{
		httpClient: client.New(usersServiceBaseURL, token),
	}
}

// User ...
// {
// 	"createdAt": "2020-02-26T12:23:57.000Z",
// 	"updatedAt": "2020-02-26T12:23:57.000Z",
// 	"deletedAt": null,
// 	"id": 181108,
// 	"roleBitMask": 2,
// 	"permissionStrategy": "ROLE_BASED",
// 	"alias": "Progress Software",
// 	"coreUser": {
// 		"id": 25826,
// 		"email": "b827f1e9-42ca-4bc2-a386-1d437e855253@ProgressSoftware.com",
// 		"firstName": "Bla",
// 		"lastName": null,
// 		"type": "programmatic"
// 	},
// 	"userId": 25826,
// 	"accountId": "act-6db0dbf7",
// 	"organizationId": 606079874257
// },
type User struct {
	ID          int    `json:"userId"`
	AccessToken string `json:"accessToken"`
	Name        string
	Description string `json:"description"`
	CoreUser    struct {
		ID        int
		FirstName string `json:"firstName"`
		Type      string
	} `json:"coreUser"`
	AccountID string `json:"accountId"`
}

type createProgrammaticUserRequest struct {
	AccountRole        int      `json:"accountRole"`
	Accounts           []string `json:"accounts"`
	Description        string   `json:"description"`
	Name               string   `json:"name"`
	PermissionStrategy string   `json:"permissionStrategy"`
	PolicyIds          []int    `json:"policyIds"`
}

type createProgrammaticUserResponse struct {
	Token string
	Type  string
	Name  string
}

// Create ..
func (us *Service) Create(username, description, accountID string) (*User, error) {

	b := &createProgrammaticUserRequest{
		AccountRole:        2,
		Accounts:           []string{accountID},
		Description:        description,
		Name:               username,
		PermissionStrategy: "ROLE_BASED",
		PolicyIds:          []int{},
	}

	req, err := us.httpClient.NewRequest(http.MethodPost, "/setup/shared/ums/programmaticUser", b)
	if err != nil {
		return nil, err
	}

	var responseBody response

	resp, err := us.httpClient.Do(req, &responseBody)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if len(responseBody.Items) == 0 {
		return nil, errors.New("Cannot provision user")
	}

	var result createProgrammaticUserRequest

	err = json.Unmarshal(responseBody.Items[0], &result)
	if err != nil {
		return nil, err
	}

	return us.Get(username, accountID)
}

// Get ...
func (us *Service) Get(username, accountID string) (*User, error) {

	b := make(map[string]interface{})
	req, err := us.httpClient.NewRequest(http.MethodGet, "/setup/shared/accountUserMapping", &b)
	if err != nil {
		return nil, err
	}

	u, _ := url.ParseQuery(req.URL.RawQuery)
	u.Add("spotinstAccountId", accountID)
	u.Add("shouldIncludeUser", "true")
	req.URL.RawQuery = u.Encode()

	var r response

	resp, err := us.httpClient.Do(req, &r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if len(r.Items) == 0 {
		return nil, errors.New("Cannot get user")
	}

	var ul []*User

	for _, jsonData := range r.Items {
		var obj User
		err = json.Unmarshal(jsonData, &obj)
		if err != nil {
			return nil, err
		}
		ul = append(ul, &obj)
	}

	return filterUserByName(username, ul)
}

func filterUserByName(username string, ul []*User) (*User, error) {
	for _, u := range ul {
		if strings.ToLower(u.CoreUser.FirstName) == username {
			fmt.Printf("%v\n", u)
			return u, nil
		}
	}
	return nil, fmt.Errorf("User %s not found", username)
}

// Update ...
func (us *Service) Update(u *User) (*User, error) {
	return u, nil
}

// Delete ...
func (us *Service) Delete(userID string) error {
	return nil
}
