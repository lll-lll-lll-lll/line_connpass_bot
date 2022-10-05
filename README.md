# line_connpass_bot

# ダイアグラム

![](line_conpass_bot.drawio.svg)


## Google Cloudの環境変数
<!-- - `USER` -->
- `CHANNEL_TOKEN`
- `CHANNEL_SECRET`


# 何を作ろうとしていたのか
目的もたず作っていて途中で断念。本当に作りたいものができたら再開しようと思う

# できること
- Connpass構造体のRequestメソッドに`key`がconnpassAPIのパラメータで`value`がqueryの値の`map`を渡せばapiを叩ける<br>
- context.Contextに初期化した`linebot.Client`を入れて値の伝播してる
