package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const usersServiceSignInURL = usersServiceBaseURL + "/api/auth/signIn"

type getConsoleTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type getConsoleTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

// GetConsoleToken issues console token for a given Spotinst user
func GetConsoleToken(email, password string) (string, error) {
	b := &getConsoleTokenRequest{
		Email:    email,
		Password: password,
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(b)

	if err != nil {
		return emptyString, err
	}
	resp, err := http.Post(usersServiceSignInURL, "application/json", &buf)

	if err != nil {
		fmt.Println(err)
		return emptyString, err
	}

	var user getConsoleTokenResponse

	r, err := readResponseBody(resp.Body)

	if err != nil {
		log.Println(err)
		return emptyString, err
	}

	if len(r.Items) == 0 {
		return emptyString, fmt.Errorf("Cannot get token for %s", email)
	}

	err = json.Unmarshal(r.Items[0], &user)

	return user.AccessToken, err
}

func readResponseBody(r io.Reader) (*response, error) {
	respBytes, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var resp response

	err = json.Unmarshal(respBytes, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
