# slash-command

公開サーバを立てずにSlash Commandを動かしたい感じのやつ。

## しくみ

1. Slash CommandsがGASの公開URLにPOSTしてくる
2. GASが受けた内容をスプレッドシートにキューイングする
3. `main.go`がGASの公開URLにGETでポーリングする
4. キューの内容を消費して処理を実行する
5. これを繰り返す

## あそびかた

1. スプレッドシートを作る
2. Google App Scriptで`receiver.gs`を追加する
3. ウェブアプリとしてデプロイする
4. 発行されたURLをSlash Commandsに設定する
5. `config.json`の`gas-url`にも設定する
6. `go run .`する
7. Slack上で`/sushi 好きな寿司`
