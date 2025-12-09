# Crypto Trading Connector - Backend API

暗号通貨取引WebアプリケーションのバックエンドAPIです。Golang、Echo フレームワーク、MVCアーキテクチャで構築されています。

## 概要

このAPIは、bitFlyer Lightning APIから現在の暗号通貨価格を取得し、MySQL（RDS）から過去の価格履歴を取得して、フロントエンドに市場データを提供します。

### 主な機能

- 暗号通貨市場データの取得（BTC、ETH）
- 個別の暗号通貨詳細情報の取得
- チャートデータの取得（期間指定可能）
- bitFlyer Lightning APIとの統合
- MySQL（RDS）からの価格履歴取得
- 日毎の平均価格計算

## 技術スタック

- **言語**: Go 1.21+
- **Webフレームワーク**: Echo v4
- **データベース**: MySQL (AWS RDS)
- **外部API**: bitFlyer Lightning API
- **アーキテクチャ**: MVC (Model-View-Controller)

## プロジェクト構造

```
crypto-trading-connector-be/
├── cmd/
│   └── server/
│       └── main.go                 # アプリケーションエントリーポイント
├── internal/
│   ├── generated/
│   │   └── models.go               # OpenAPIから生成されたモデル
│   ├── handler/
│   │   └── crypto_handler.go      # HTTPハンドラー
│   ├── service/
│   │   └── crypto_service.go      # ビジネスロジック
│   ├── repository/
│   │   └── crypto_repository.go   # データアクセス層
│   ├── client/
│   │   └── bitflyer_client.go     # bitFlyer APIクライアント
│   └── model/
│       └── crypto.go               # データモデル
├── pkg/
│   └── database/
│       └── db.go                   # データベース接続
├── test/
│   └── e2e/
│       └── crypto_e2e_test.go     # E2Eテスト
├── .env.example                    # 環境変数サンプル
├── Makefile                        # ビルド・テストコマンド
├── go.mod
└── README.md
```

## セットアップ

### 前提条件

- Go 1.21以上
- MySQL（AWS RDS）へのアクセス
- bitFlyer Lightning APIへのアクセス

### インストール

1. リポジトリをクローン

```bash
git clone <repository-url>
cd crypto-trading-connector-be
```

2. 依存関係をインストール

```bash
go mod download
```

3. 環境変数を設定

`.env.example`を`.env`にコピーして、必要な値を設定します：

```bash
cp .env.example .env
```

`.env`ファイルの内容：

```env
# Database Configuration
DB_HOST=your-rds-endpoint.rds.amazonaws.com
DB_PORT=3306
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=your-database

# Server Configuration
SERVER_PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000

# bitFlyer API Configuration
BITFLYER_API_URL=https://api.bitflyer.com
```

## 使用方法

### サーバーの起動

```bash
make run
```

サーバーは`http://localhost:8080`で起動します。

### テストの実行

#### ユニットテストの実行

```bash
make test
```

#### E2Eテストの実行

E2Eテストは実際のデータベースとbitFlyer APIに接続します。サーバーを起動してから実行してください：

```bash
# ターミナル1: サーバーを起動
make run

# ターミナル2: E2Eテストを実行
make e2e-test
```

### コードのフォーマット

```bash
make fmt
```

### OpenAPIからコード生成

OpenAPI仕様書（`../openapi.yaml`）からGoのモデルを生成します：

```bash
make generate
```

生成されたコードは`internal/generated/models.go`に出力されます。

**使用例:**

```go
import "github.com/crypto-trading-connector/backend/internal/generated"

// 生成された型を使用
var cryptoData generated.CryptoData
var marketResponse generated.MarketResponse
```

### ビルド

```bash
make build
```

バイナリは`bin/server`に生成されます。

## API エンドポイント

### GET /api/v1/crypto/market

すべての暗号通貨の市場データを取得します。

**レスポンス例:**

```json
{
  "data": [
    {
      "id": "bitcoin",
      "name": "Bitcoin",
      "symbol": "BTC",
      "pair": "BTC/JPY",
      "icon": "₿",
      "iconColor": "#f7931a",
      "currentPrice": 9850000,
      "changePercent": 5.2,
      "chartData": [
        {"day": "Mon", "price": 9350000},
        {"day": "Tue", "price": 9450000}
      ]
    }
  ],
  "timestamp": 1704067200
}
```

### GET /api/v1/crypto/:id

特定の暗号通貨の詳細データを取得します。

**パラメータ:**
- `id` (path): 暗号通貨ID（`bitcoin`, `ethereum`）

**レスポンス例:**

```json
{
  "id": "bitcoin",
  "name": "Bitcoin",
  "symbol": "BTC",
  "pair": "BTC/JPY",
  "icon": "₿",
  "iconColor": "#f7931a",
  "currentPrice": 9850000,
  "changePercent": 5.2,
  "chartData": [...]
}
```

### GET /api/v1/crypto/:id/chart

チャートデータを取得します。

**パラメータ:**
- `id` (path): 暗号通貨ID
- `period` (query, optional): 期間（`24h`, `7d`, `30d`, `1y`, `all`）、デフォルト: `7d`

**レスポンス例:**

```json
{
  "data": [
    {"day": "Mon", "price": 9350000},
    {"day": "Tue", "price": 9450000}
  ],
  "period": "7d"
}
```

### エラーレスポンス

エラーが発生した場合、以下の形式でレスポンスが返されます：

```json
{
  "error": "NOT_FOUND",
  "message": "Cryptocurrency not found"
}
```

**エラータイプ:**
- `NOT_FOUND` (404): リソースが見つからない
- `BAD_REQUEST` (400): 不正なリクエスト
- `INTERNAL_SERVER_ERROR` (500): サーバー内部エラー

## データソース

### 現在価格

bitFlyer Lightning API (`GET /v1/ticker`) から取得：
- BTC_JPY: ビットコインの円建て価格
- ETH_JPY: イーサリアムの円建て価格

### 過去価格データ

MySQL（RDS）の`price_histories`テーブルから取得：
- 日毎の平均価格を計算（`DATE(datetime)`でグループ化、`AVG(price)`）
- 指定された期間のデータを取得

## 開発

### Makefileコマンド

```bash
make run        # サーバーを起動
make test       # ユニットテストを実行
make e2e-test   # E2Eテストを実行
make fmt        # コードをフォーマット
make vet        # go vetを実行
make lint       # linterを実行（golangci-lintが必要）
make tidy       # go mod tidyを実行
make build      # バイナリをビルド
make clean      # ビルド成果物を削除
make help       # ヘルプを表示
```

### テスト戦略

#### ユニットテスト

各層（Handler、Service、Repository）を独立してテストします。モックを使用してデータベースや外部APIへの依存を排除します。

#### E2Eテスト

実際のデータベースとbitFlyer APIに接続して、エンドツーエンドの動作を確認します。ビルドタグ`//go:build e2e`を使用して、ユニットテストと分離しています。

## トラブルシューティング

### データベース接続エラー

- `.env`ファイルのデータベース接続情報を確認してください
- RDSのセキュリティグループでアクセスが許可されているか確認してください
- データベースが起動しているか確認してください

### bitFlyer API エラー

- インターネット接続を確認してください
- bitFlyer APIのレート制限に達していないか確認してください
- `BITFLYER_API_URL`が正しく設定されているか確認してください

### CORS エラー

- `CORS_ALLOWED_ORIGINS`にフロントエンドのURLが含まれているか確認してください
- デフォルトは`http://localhost:3000`です

## ライセンス

MIT License

## 作者

Crypto Trading Connector Team
