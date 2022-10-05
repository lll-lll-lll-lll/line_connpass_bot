# line_connpass_bot

# ダイアグラム

![](line_conpass_bot.drawio.svg)


## Google Cloudの環境変数
<!-- - `USER` -->
- `CHANNEL_TOKEN`
- `CHANNEL_SECRET`



# 何を作ろうとしていたのか
目的もたず作っていて途中で断念。本当に作りたいものができたら再開しようと思う
他のアプリ開発に時間を使おうと思う

# できること
- `go`をキーワードに５件取得検索
- Connpass構造体のRequestメソッドに`key`がconnpassAPIのパラメータで`value`がqueryの値の`map`を渡せばapiを叩ける<br>
- `LINEClientMiddleware()`でcontext.Contextに初期化した`linebot.Client`を入れて値の伝播してる
<br>
<br>

![RPReplay_Final1664973665_MP4_AdobeExpress](https://user-images.githubusercontent.com/63499912/194063023-14248956-8fcc-478e-b873-c0ba52474389.gif)
