package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kinvey/terraform-provider-spotinstadmin/client"
	"github.com/kinvey/terraform-provider-spotinstadmin/client/common"
)

const (
	usersServiceBaseURL = "https://api.spotinst.io"
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

type User struct {
	ID          string `json:"userId"`
	UserName    string `json:"username"`
	Type        string `json:"type"`
}

type UserDetails struct {
	ID          string `json:"userId"`
	AccessToken string `json:"accessToken"`
	UserName    string `json:"username"`
	Description string `json:"description"`
}

type createProgrammaticUserAccount struct{
	ID   string   `json:"id"`
	Role string   `json:"role"`
} 

type createProgrammaticUserRequest struct {
	Accounts           []createProgrammaticUserAccount `json:"accounts"`
	Description        string   `json:"description"`
	Name               string   `json:"name"`
}

type createProgrammaticUserResponse struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	ID     string `json:"id"`
}

// Create ..
func (us *Service) Create(username, description, accountID string) (*UserDetails, error) {

	b := &createProgrammaticUserRequest{
		Accounts:           []createProgrammaticUserAccount{{ID: accountID, Role: "editor"}},
		Description:        description,
		Name:               username,
	}
    
	req, err := us.httpClient.NewRequest(http.MethodPost, "/setup/user/programmatic", b)
	if err != nil {
		return nil, err
	}

	var v common.Response

	_, err = us.httpClient.Do(req, &v)
	if err != nil {
		return nil, err
	}

	if len(v.Response.Errors) > 0 {
		return nil, errors.New(v.Response.Errors[0].Message)
	}
	if len(v.Response.Items) == 0 {
		r, e := json.Marshal(b)
		rr, er:= json.Marshal(v)
		if e!= nil || er != nil { return nil, errors.New("Cannot provision user")} 
		return nil, errors.New(fmt.Sprintf("Cannot provision user: %#v %#v", string(r), string(rr)))
	}

	var result createProgrammaticUserResponse

	err = json.Unmarshal(v.Response.Items[0], &result)
	if err != nil {
		return nil, err
	}

	log.Printf("[TRACE] IN CREATE: %v\n", accountID)
	user, err := us.Get(username, accountID)
	if err != nil {
		return nil, err
	}

	user.AccessToken = result.Token
	return user, nil
}

// Get ...
func (us *Service) Get(username, accountID string) (*UserDetails, error) {

	req, err := us.httpClient.NewRequest(http.MethodGet, "/setup/organization/user", nil)
	if err != nil {
		return nil, err
	}

	log.Printf("[TRACE] IN GET: %v\n", accountID)

	var r common.Response

	_, err = us.httpClient.Do(req, &r)

	if err != nil {
		return nil, err
	}

	userList, err := usersFromJSON(r)

	if err != nil {
		return nil, err
	}

	if len(userList) == 0 {
		return nil, errors.New("Cannot get users")
	}

	user, err := filterUserByName(username, userList)

	if err != nil {
		return nil, err
	}

	userDetails, err := us.GetDetails(user.ID, accountID)

	return userDetails, nil
}

func usersFromJSON(r common.Response) ([]*User, error) {
	userList := make([]*User, len(r.Response.Items))

	for i, jsonData := range r.Response.Items {
		var obj User
		err := json.Unmarshal(jsonData, &obj)
		if err != nil {
			return nil, err
		}
		userList[i] = &obj
	}

	return userList, nil

}

func filterUserByName(username string, ul []*User) (*User, error) {
	for _, u := range ul {

		log.Printf("[TRACE] %v\n", u)
		log.Printf("[TRACE] Checking %v with %v\n", u.UserName, username)
		if strings.ToLower(u.UserName) == username {
			return u, nil
		}
	}
	return nil, fmt.Errorf("User %s not found", username)
}

func (us *Service) GetDetails(userId, accountID string) (*UserDetails, error) {
	reqString := fmt.Sprintf("/setup/user/%s", userId)
	req, err := us.httpClient.NewRequest(http.MethodGet, reqString, nil)

	if err != nil {
		return nil, err
	}

	var r common.Response

	_, err = us.httpClient.Do(req, &r)

	if err != nil {
		return nil, err
	}

	if len(r.Response.Items) == 0 {
		return nil, errors.New("Cannot get user")
	}
	
	var user *UserDetails
	err = json.Unmarshal(r.Response.Items[0], &user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update ...
func (us *Service) Update(u *User) (*User, error) {
	return u, nil
}

// Delete ...
func (us *Service) Delete(username, accountID string) error {
	user, err := us.Get(username, accountID)
	if err != nil {
		return err
	}

	log.Println(user.ID)
	req, err := us.httpClient.NewRequest(http.MethodDelete, fmt.Sprintf("/setup/user/%v", user.ID), nil)

	resp, err := us.httpClient.Do(req, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("[TRACE] RESPONSE !!!!!!!!!!!! %#v", resp)

	if resp.StatusCode > 399 {
		return errors.New("Cannot delete user: " + user.UserName)
	}
	return nil
}
