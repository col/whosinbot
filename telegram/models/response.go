package models

import (
	"encoding/json"
)

type Response struct {
	Message string `json:"message"`
}

func (r Response) String() (string) {
	body, _ := json.Marshal(r)
	return string(body)
}
