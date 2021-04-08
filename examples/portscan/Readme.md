# 使い方

## 事前準備 (ローカル環境)

`go run` を実行できる以外にローカルに以下のインストールが必要です。
- [nmap](https://nmap.org/man/ja/index.html)のローカルへのインストール
- `go run` をsudo,特権ユーザーで実行できる環境
  - nmapのスキャンオプションの影響で特権ユーザーでないと失敗するため

## スキャン実行

- portscan.goのscanTargetを確認する
  - このファイル作成時点では、google.comへのnmapが有効になっています
- 必要に応じてローカル、自身の持つ環境を追記してください
- `sudo go run portscan.go` にて実行
- コンソールにnmapの結果、findingが出力されます

## 追加ターゲット
現状コメントアウトされているのは以下のようなターゲットです
- OpenProxy(http)が有効になっているとFindingが追加されるターゲット
  - zap,burp,squidなど通信が届く部分に構築してください
- OpenRelay(smtp)が有効になっているとFindingが追加されるターゲット
  - postfixなど通信が届く部分に構築してください
- PasswordAuth(ssh)が有効になっているとFindingが追加されるターゲット
  - sshなど通信が届く部分に構築してください
