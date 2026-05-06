# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

暗号通貨(BTC, ETH)取引Webアプリケーション。bitFlyer Lightning APIと連携し、価格確認・買い注文・自動売買利益表示を行う。Kiroによる仕様駆動開発。

DBや自動売買ロジックは別リポジトリ(https://github.com/Kohei-Sato-1221/crypto-trading-golang)で管理されており、本アプリはそのDBを参照する。DB マイグレーションは本プロジェクトでは行わない。

## モノレポ構成

- `crypto-trading-connector-be/` — バックエンド (Go + Echo v4, MVC)
- `crypto-trading-connector-fe/` — フロントエンド (Nuxt 4 + Vue 3 + TypeScript + Tailwind CSS)
- `openapi.yaml` — API仕様書 (OpenAPI 3.0.3)、バックエンド・フロントエンド共通の型定義ソース
- `db/` — 参照用Atlasスキーマのみ（マイグレーション実装なし）

## よく使うコマンド

### OpenAPIコード生成（API変更時は必ず実行）
```bash
make gen        # BE + FE 両方生成
make gen-be     # BE のみ（internal/generated/models.go）
make gen-fe     # FE のみ（types/api.ts）
```

### バックエンド（crypto-trading-connector-be/）
```bash
make run        # サーバー起動 (localhost:8080)
make test       # ユニットテスト
make e2e-test   # E2Eテスト（サーバー起動が必要、タグ: //go:build e2e）
make build      # バイナリビルド (bin/server)
make fmt        # go fmt
make lint       # golangci-lint
make curl a=market  # APIエンドポイントテスト
```

### フロントエンド（crypto-trading-connector-fe/）
```bash
make run        # 開発サーバー起動（API接続モード、localhost:3000）
make run-mock   # モックデータモードで起動
make test       # Vitest実行
make build      # プロダクションビルド
make gen        # OpenAPIから型定義生成
```

### 両方同時起動
```bash
make run        # ルートのMakefile（scripts/start-services.sh経由）
```

## リグレッションテスト（必須）

コード変更後は**必ず**該当ディレクトリで`make test`を実行すること。例外なし。
- バックエンド変更 → `cd crypto-trading-connector-be && make test`
- フロントエンド変更 → `cd crypto-trading-connector-fe && make test`
- 両方変更 → 両方実行

## アーキテクチャ

### バックエンド: レイヤードアーキテクチャ
```
handler → service → repository（MySQL/RDS or PostgreSQL/Supabase）
                  → client（bitFlyer API）
```
- **Exchange Client**: Strategy Patternで`CryptoExchangeClient`インターフェースを定義。現在はbitFlyerのみ実装。新取引所追加時はインターフェースを実装する。
- **OpenAPI生成コード**: `internal/generated/models.go` — 手動編集禁止
- **CLIツール**: `cmd/`配下に`buy-order`, `get-balance`, `get-current-orders`

### フロントエンド: Nuxt 4 Composition API
- **composables/**: API呼び出し・状態管理のロジック（`useApi`, `useCryptoData`, `useTradeHistory`等）
- **環境変数**: `USE_MOCK_DATA`でモック/API切替、`API_BASE_URL`でバックエンドURL指定
- **テスト**: Vitest + fast-check（プロパティベーステスト）
- **OpenAPI生成コード**: `types/api.ts` — 手動編集禁止

### API開発ワークフロー
1. `openapi.yaml`を編集
2. `make gen`でBE/FE両方のコード生成
3. BEにhandler/service/repositoryを実装
4. FEでcomposableを更新
