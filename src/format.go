package v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 今月を含めた３月分のイベントを取得
func GetForThreeMonthsEvent() string {
	now := time.Now()
	yearmonthsja := strings.NewReplacer(
		"January", "01",
		"February", "02",
		"March", "03",
		"April", "04",
		"May", "05",
		"June", "06",
		"July", "07",
		"August", "08",
		"September", "09",
		"October", "10",
		"November", "11",
		"December", "12",
	)
	// 13月をなくすために12で割った余を入れる
	nm := now.Month()
	sm, tm := CheckMonthFormat(nm)
	// a := yearmonthsja.Replace(fmt.Sprintf("%s", tm.String()))

	f := yearmonthsja.Replace(fmt.Sprintf("%d%s", now.Year(), nm.String()))
	s := yearmonthsja.Replace(fmt.Sprintf("%d%s", now.Year(), sm))
	t := yearmonthsja.Replace(fmt.Sprintf("%d%s", now.Year(), tm))
	return f + "," + s + "," + t
}

// 12で割るときに12月だけ0が返ってくるので、その時だけ"12"文字列を返す
func CheckMonthFormat(nm time.Month) (string, string) {
	ze := "%!Month(0)"
	sm := (nm + 1) % 12
	tm := (nm + 2) % 12
	if sm.String() == ze {
		return "12", tm.String()
	} else if tm.String() == ze {
		return sm.String(), "12"
	}
	return sm.String(), tm.String()
}

// 時刻を見やすいように変更
func ConvertStartAtTime(startedAt string) (string, error) {
	weekdaymonthja := strings.NewReplacer(
		"Sunday", "日",
		"Monday", "月",
		"Tueday", "火",
		"Wednesday", "水",
		"Thursday", "木",
		"Friday", "金",
		"Saturday", "土",
		"January", "1月",
		"February", "2月",
		"March", "3月",
		"April", "4月",
		"May", "5月",
		"June", "6月",
		"July", "7月",
		"August", "8月",
		"September", "9月",
		"October", "10月",
		"November", "11月",
		"December", "12月",
	)
	p, err := time.Parse(time.RFC3339, startedAt)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	f := func(p time.Time) string {
		if p.Minute() == 0 {
			return "00"
		} else {
			return strconv.Itoa(p.Minute())
		}
	}
	str := fmt.Sprintf("%s%d日(%s) %d:%s ~", p.Month().String(), p.Day(), p.Weekday(), p.Hour(), f(p))
	return weekdaymonthja.Replace(str), nil
}
