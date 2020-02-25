package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    *url.URL
	UserAgent  string
	token      string
	httpClient *http.Client
}

func New(baseURLString, token string) *Client {
	baseURL, err := url.Parse(baseURLString)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout:   time.Second * 5,
		Transport: http.DefaultTransport,
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: client,
		token:      token,
	}
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {

		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)

		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", c.getToken())
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) getToken() string {
	return "Bearer " + c.token
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	return resp, err
}
