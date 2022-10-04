package main

import (
	"fmt"
	"log"
	"testing"

	linecon "github.com/lll-lll-lll-lll/lineconnpass/v1"
)

func TestConnpassAPI(t *testing.T) {
	t.Run("connpass api叩く", func(t *testing.T) {
		conpass := linecon.NewConnpass()
		query := map[string]string{"keyword": "go"}
		q := linecon.CreateQuery(query)
		conpass.Query = q
		u, err := conpass.CreateURL(conpass.Query)
		if err != nil {
			log.Println(fmt.Errorf("no client: %s", err))
			return
		}
		res, err := conpass.Request(u)
		if err != nil {
			log.Println(fmt.Errorf("no client: %s", err))
			return
		}
		defer res.Body.Close()
		if err := conpass.SetResponse(res); err != nil {
			log.Println(fmt.Errorf("no client: %s", err))
			return
		}
		t.Log(conpass.ConnpassResponse.ResultsReturned)
		for _, v := range conpass.ConnpassResponse.Events {
			t.Log(v.Title)
			t.Log("series title", v.Series.Title)
		}
		t.Log(linecon.CreateConnpassEventFlexMessages(conpass.ConnpassResponse.Events))
	})

}
