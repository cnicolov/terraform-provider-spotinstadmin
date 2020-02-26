package users

import "encoding/json"

type response struct {
	Kind  string            `json:"kind"`
	Items []json.RawMessage `json:"items"`
}
