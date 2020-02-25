package users

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	consoleSignInURL = "https://console.spotinst.com/signIn"
	emptyString      = ``
)

type getConsoleTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type getConsoleTokenResponse struct {
	AuthToken string `json:"authToken"`
}

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

	if err != nil {
		return emptyString, err
	}

	defer resp.Body.Close()

	var r *getConsoleTokenResponse

	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return emptyString, err
	}

	return r.AuthToken, nil
}
