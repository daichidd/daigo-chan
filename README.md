## 第五ちゃん
第五ちゃんは、discord【おせち様教団】のtool botです。
herokuのdyno上で動いてます。

## envについて
開発環境では.envで動いていますが、
本番環境はherokuのため、osの環境変数で動いてます。

heroku-configを利用して反映した後、herokuの管理画面で
【HEROKU_ENV】を【true】にする必要がります。
.env.exampleを参照してください。
