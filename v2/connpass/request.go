package v2

import (
	"log"
	"net/http"
	"net/url"
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

func (c *Connpass) CreateUrl(q url.Values) string {
	u, err := url.Parse(CONNPASSAPIV1)
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "connpass.com"
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Connpass) Request(url string) *http.Response {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	return res
}
