package v2

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func NewConnpassResponse() *ConnpassResponse {
	return &ConnpassResponse{}
}

// SetResponse Requestメソッド後のレスポンスをConnpassResponseプロパティにセットする
func (c *Connpass) SetResponse(res *http.Response) error {
	body, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(body, &c.ConnpassResponse)
	if err != nil {
		return err
	}
	return nil
}

// JoinGroupIdsByComma groupidを「,」で繋げる。connpassapiで複数指定は「,」で可能だから
func (c *Connpass) JoinGroupIdsByComma() string {
	var seriesId string
	gs := c.ConnpassResponse.GetGroupIds()
	for _, v := range gs {
		v := strconv.Itoa(v)
		seriesId += v + ","
	}
	return seriesId
}

// GetGroups 所属してるグループIDを取得
func (c *ConnpassResponse) GetGroupIds() []int {
	var g []int
	for _, v := range c.Events {
		g = append(g, v.Series.Id)
	}
	return g
}

type Event struct {
	Series struct {
		// グループID
		Id int `json:"id"`
		// グループのタイトル
		Title string `json:"title"`
		// グループのconnpass.com 上のURL
		Url string `json:"url"`
	} `json:"series"`
	// 管理者のニックネーム
	OwnerNickname string `json:"owner_nickname"`
	// キャッチコピー
	Catch string `json:"catch"`
	// 概要(HTML形式)
	Description string `json:"description"`
	// connpass.com 上のURL
	EventUrl string `json:"event_url"`
	// Twitterのハッシュタグ
	HashTag string `json:"hash_tag"`
	// イベント開催日時 (ISO-8601形式)
	StartedAt string `json:"started_at"`
	// イベント終了日時 (ISO-8601形式)
	EndedAt string `json:"ended_at"`
	// 管理者の表示名
	OwnerDisplayName string `json:"owner_display_name"`
	// イベント参加タイプ
	EventType string `json:"event_type"`
	// タイトル
	Title string `json:"title"`
	// 開催場所
	Address string `json:"address"`
	// 開催会場
	Place string `json:"place"`
	// 更新日時 (ISO-8601形式)
	UpdatedAt string `json:"updated_at"`
	// イベントID
	EventId int `json:"event_id"`
	// 管理者のID
	OwnerId int `json:"owner_id"`
	// 定員
	Limit int `json:"limit"`
	// 参加者数
	Accepted int `json:"accepted"`
	// 補欠者数
	Waiting int `json:"waiting"`
	// 開催会場の緯度
	Lat string `json:"lat"`
	// 開催会場の経度
	Lon string `json:"lon"`
}

// ConnpassResponse コンパスapiのレスを持つ
type ConnpassResponse struct {
	// レスポンスに含まれる検索結果の件数
	ResultsReturned int `json:"results_returned"`
	// 検索結果の総件数
	ResultsAvailable int `json:"results_available"`
	// 検索の開始位置
	ResultsStart int `json:"results_start"`
	// 検索結果のイベントリスト
	Events []Event `json:"events"`
}
