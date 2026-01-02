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

## Seaweedfs 概要 (compose.yaml)

S3 互換のオブジェクトストレージとして SeaweedFS という分散オブジェクトストレージを使用する

### SeaweedFSアーキテクチャ概要

４つのコンポーネントが階層的に連携して動作

```text
      ┌─────────────┐
      │   S3 API    │ ← S3互換インターフェース
      └──────┬──────┘
             │
      ┌──────▼──────┐
      │   Filer     │ ← ファイルシステムインターフェース
      └──────┬──────┘
             │
    ┌────────┼─────────┐
    │        │         │
┌───▼───┐ ┌──▼────┐ ┌──▼────┐
│Master │ │Volume │ │Volume │ ← データ保存層
└───────┘ └───────┘ └───────┘
```

1. Master Service

- 役割: クラスタ全体の司令塔として機能
  - ボリュームサーバーの一覧を管理
  - ファイルをどのボリュームに配置するかを決定
  - 容量の監視とロードバランシング
  - メタデータの管理 (どのファイルIDがどこにあるか)

- 使用ポート
  - 9333: HTTP/REST API用ポート
    - ブラウザでアクセス可能: http://localhost:9333
    - ボリューム一覧の確認、クラスタの状態確認などに使用
  - 19333: gRPC API用ポート
    - 内部サービス間の高速通信に使用

2. Volume Service

- 役割: 実際のファイルデータを保存するストレージノード
  - ファイルの実体をディスク上に保存
  - ファイルの読み書きリクエストを処理
  - データの複製（レプリケーション）をサポート
  - ファイルIDに基づいた高速なファイル検索

- 使用ポート
  - 8080: ファイルのアップロード/ダウンロード用 HTTP API
    - `curl http://localhost:8080/{file_id}` でファイル取得可能
  - 18080: gRPC API (内部通信用)

3. Filer サービス (compose.yaml:20-32)

- 役割: ファイルシステムの抽象化レイヤーとして機能
  - ディレクトリ構造の管理 (パスベースのアクセス)
  - ファイル名 → ファイルIDの変換
  - メタデータの保存 (ファイル名、タイムスタンプ、権限など)
  - POSIX ライクなファイルシステムインターフェースの提供
  - SeaweedFSの低レベルAPI (ファイルID直指定) を、高レベルAPI (パスベース) に変換

- 使用ポート
  - 8888: Filer HTTP API
    - `curl http://localhost:8888/path/to/file` でファイルアクセス可能
    - WebDAV Protocol サポート
  - 18888: gRPC API

4. S3 Service

- 役割: AWS S3 互換の API ゲートウェイとして機能
  - S3プロトコル (GetObject, PutObject, ListBucketsなど) のサポート
  - 既存のS3クライアントツール (aws-cli, boto3など) との互換性
  - Filerを通じてファイルにアクセス

- 使用ポート
  - 8333: S3互換API endpoint
    - `http://localhost:8333` を S3 エンドポイントとして設定可能
    - バケット操作、オブジェクトの PUT/GET/DELETE が可能

### 起動方法

**コンテナ起動**

```bash
docker compose up -d
```

**ログ確認**

```bash
docker compose logs -f
```

**停止**

```bash
docker compose down
```

### アクセスポイント

- Master管理UI: http://localhost:9333
- Filer API: http://localhost:8888
- S3 API: http://localhost:8333
- Volume API: http://localhost:8080
