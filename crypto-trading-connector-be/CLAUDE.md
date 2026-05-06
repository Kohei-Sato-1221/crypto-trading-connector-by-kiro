# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## バックエンド概要

Go + Echo v4 によるREST APIサーバー。bitFlyer Lightning APIから暗号通貨価格を取得し、MySQL(RDS)またはPostgreSQL(Supabase)から履歴データを取得してフロントエンドに提供する。`DB_TYPE`環境変数でDB種別を切り替え可能。

## よく使うコマンド

```bash
make run          # サーバー起動
make test         # ユニットテスト（変更後は必ず実行）
make e2e-test     # E2Eテスト（サーバー起動必須、//go:build e2eタグ）
make build        # バイナリビルド → bin/server
make fmt          # go fmt
make lint         # golangci-lint
make gen          # OpenAPIからモデル生成 → internal/generated/models.go
make curl a=market  # APIテスト（market, bitcoin, ethereum, bitcoin-chart, current-orders）
make get-balance  # bitFlyerから残高取得
make buy-order    # BTC/ETHの買い注文（現在価格の97%で実発注）
```

## アーキテクチャ

### レイヤー構成
```
cmd/server/main.go     → エントリーポイント、DI
internal/handler/      → HTTPハンドラー（Echo）
internal/service/      → ビジネスロジック
internal/repository/   → データアクセス（MySQL/PostgreSQL両対応）
                         MySQL版: *_repository.go / PostgreSQL版: *_repository_postgres.go
internal/client/       → 外部APIクライアント
internal/model/        → ドメインモデル
internal/generated/    → OpenAPI生成コード（手動編集禁止）
pkg/database/          → DB接続
```

### Exchange Client (Strategy Pattern)
`CryptoExchangeClient`インターフェースを通じて取引所APIを抽象化。
- `bitflyer_client.go` — HMAC-SHA256認証付きbitFlyer実装
- `bitflyer_client_mock.go` — テスト用モック
- 新取引所追加時: インターフェースを実装し`main.go`で切り替え

### CLIツール（cmd/配下）
- `buy-order/` — BTC/ETH買い注文CLI（実際に発注するので注意）
- `get-balance/` — 残高確認CLI
- `get-current-orders/` — 未約定注文一覧CLI

## テスト

- ユニットテスト: `go-sqlmock`でDB、モッククライアントで外部API依存を排除
- プロパティベーステスト: `gopter`使用
- E2Eテスト: `test/e2e/`、ビルドタグ`e2e`で分離、実DB/API接続が必要

## 環境変数（.envファイル）

### DB接続（URL方式）
- `MYSQL_URL` — MySQL接続URL（例: `user:pass@tcp(host:port)/dbname?parseTime=true&loc=Local`）
- `POSTGRES_URL` — PostgreSQL接続URL（例: `postgresql://user:pass@host:port/dbname?sslmode=require`）
- 両方設定した場合はPostgreSQLが優先。両方空の場合はpanic。

### その他
`SERVER_PORT`、`CORS_ALLOWED_ORIGINS`、`BITFLYER_API_URL`、bitFlyer認証キー
