package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const CONNPASSAPIV1 = "https://connpass.com/api/v1/event/?"

type Connpass struct {
	ConnpassUSER     string            `json:"user"`
	ConnpassResponse *ConnpassResponse `json:"connpass"`
	Query            url.Values        `json:"query"`
}

func NewConnpass() *Connpass {
	return &Connpass{}
}

func (c *Connpass) CreateUrl(q url.Values) (string, error) {
	u, err := url.Parse(CONNPASSAPIV1)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	u.Scheme = "https"
	u.Host = "connpass.com"
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *Connpass) Request(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return res, nil
}

func GetUserName() string {
	user := os.Getenv("USER")
	if user == "" {
		return ""
	}
	return user
}
