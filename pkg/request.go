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

func (c *Connpass) CreateURL(q url.Values) (string, error) {
	u, err := url.Parse(CONNPASSAPIV1)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	u.Scheme = "https"
	u.Host = "connpass.com"
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *Connpass) Request(conpass *Connpass, query map[string]string) error {
	q := CreateQuery(query)
	conpass.Query = q
	url, err := conpass.CreateURL(conpass.Query)
	if err != nil {
		return err
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("no success connpass api \n %s", err)
	}
	defer res.Body.Close()
	if err := conpass.SetResponse(res); err != nil {
		return fmt.Errorf("no client: %s", err)
	}
	return nil
}

func GetUserName() string {
	user := os.Getenv("USER")
	if user == "" {
		return ""
	}
	return user
}

func CreateQuery(values map[string]string) url.Values {
	q := url.Values{}
	for k, v := range values {
		q.Add(k, v)
	}
	return q
}
