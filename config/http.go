package config

import (
	"fmt"
	"io"
	"net/http"
)

func NewHttpRequest(method string, url string, body io.Reader, accessToken string) *http.Request {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return req
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", accessToken)

	return req
}
