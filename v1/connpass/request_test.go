package v1

import (
	"testing"
)

func TestConnpassAPI(t *testing.T) {
	t.Run("connpass api叩く", func(t *testing.T) {
		con := NewConnpass()
		con.ConnpassUSER = "Shun_Pei"
		query := map[string]string{"nickname": con.ConnpassUSER}
		q := CreateQuery(query)
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
		groups := con.ConnpassResponse.GetGroupIds()
		groupsNum := len(groups)
		t.Log(groups)
		if groupsNum != 10 {
			t.Error("グループ数に間違いあり")
		}

	})
}
