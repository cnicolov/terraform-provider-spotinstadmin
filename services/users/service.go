package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	consoleSignInURL = "https://console.spotinst.com/auth/signIn"
	emptyString      = ``
)

type getConsoleTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type getConsoleTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type response struct {
		Kind string `json:"kind"`
		Items []json.RawMessage `json:"items"`
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
	resp, err := http.Post(consoleSignInURL, "application/json", &buf)

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return emptyString, err
	}

	defer resp.Body.Close()

	var r response

	err = json.Unmarshal(respBytes, &r)

	if err != nil {
		return emptyString, err
	}

	if len(r.Items) == 0 {
		return emptyString, fmt.Errorf(`Cannot issue token for user "%s"`, email)
	}

	var user getConsoleTokenResponse

	err = json.Unmarshal(r.Items[0], &user)

	return user.AccessToken, err
}
