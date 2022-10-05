package main

import (
	"fmt"
	"log"
	"testing"

	linecon "github.com/lll-lll-lll-lll/lineconnpass/src"
)

func TestConnpassAPI(t *testing.T) {
	t.Run("connpass api叩く", func(t *testing.T) {
		conpass := linecon.NewConnpass()
		query := map[string]string{"keyword": "go"}

		err := conpass.Request(conpass, query)
		if err != nil {
			log.Println(fmt.Errorf("no client: %s", err))
			return
		}

		t.Log(conpass.ConnpassResponse.ResultsReturned)
		for _, v := range conpass.ConnpassResponse.Events {
			s := linecon.CreateConnpassEventFlexMessages(v)
			// fmt.Printf("%T", s)
			t.Log(s)
		}

	})

}
