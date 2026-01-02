# 並行処理の学習

## 📚 Learn

1. Goの並行処理の実装パターン
2. [A Tour of Go](https://go-tour-jp.appspot.com/list) での実践

- `patterns/`: 並行処理の定番パターン
- `exercise/`: A Tour of Go の演習

## ⛹️‍♂️ 実践: ウェブサイトのグレースケール変換 API を実装

以下のようなシステムを作る

- `worker/` に処理実装
- 処理の流れ
  - ウェブサイトの URL を指定し、スクリーンショットを撮る
  - スクリーンショットをグレースケールに変換する
  - スクリーンショットと変換後の画像をストレージに保存する

### 🦜 Concurrency

以下の3つのステージに分けて並列処理を実装する (Pipeline Pattern)

- Stage 1: 指定した URL のサイトを ChromeDP でスクリーンショットする
- Stage 2: 取得した画像をグレースケールに変換する
- Stage 3: 取得した画像と、変換後の画像を外部ストレージに保存する
