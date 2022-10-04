package main

import (
	"fmt"
	"testing"

	linecon "github.com/lll-lll-lll-lll/lineconnpass/v1"
)

func TestConnpassAPI(t *testing.T) {
	con := linecon.NewConnpass()
	con.ConnpassUSER = "Shun_Pei"
	query := map[string]string{"nickname": con.ConnpassUSER}
	q := linecon.CreateQuery(query)
	con.Query = q
	u, err := con.CreateURL(con.Query)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := con.Request(u)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()

	if err := con.SetResponse(res); err != nil {
		t.Error(err)
		return
	}
	t.Run("connpass api叩く", func(t *testing.T) {
		t.Log(RecursiveCreateConnpassEventFlexMessages(con.ConnpassResponse.Events, len(con.ConnpassResponse.Events)-1))
		for _, e := range con.ConnpassResponse.Events {
			t.Log(e.Title)
		}
		fmt.Println(len(con.ConnpassResponse.Events))
	})

	// t.Run("event test", func(t *testing.T) {
	// 	events := con.ConnpassResponse.Events
	// 	fmt.Println(len(events))
	// 	for _, v := range events {
	// 		fmt.Print(v.Title, "---")
	// 		fmt.Print(v.Accepted, "---")
	// 		fmt.Print(v.EventUrl, "---")
	// 		fmt.Println("")
	// 	}
	// })
}
