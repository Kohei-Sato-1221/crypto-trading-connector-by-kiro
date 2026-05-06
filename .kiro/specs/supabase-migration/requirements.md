# Requirements Document

## Introduction

暗号通貨取引バックエンドのデータベース接続を MySQL (AWS RDS) から Supabase (PostgreSQL) に移行します。MySQL に戻す可能性があるため、DB 接続を interface 化して切り替え可能な設計とします。修正対象はバックエンドのみで、フロントエンド・API 仕様は変更しません。

## Glossary

- **Supabase**: PostgreSQL ベースの BaaS（Backend as a Service）
- **DB_TYPE**: MySQL / PostgreSQL を切り替える環境変数
- **DSN**: Data Source Name（データベース接続文字列）
- **Repository**: データアクセス層の抽象化インターフェース
- **lib/pq**: Go の PostgreSQL ドライバ

## Requirements

### Requirement 1

**User Story:** 開発者として、環境変数で MySQL と PostgreSQL を切り替えたい。そうすることで、デプロイ先に応じて DB を変更できる。

#### Acceptance Criteria

1. WHEN `DB_TYPE=postgres` が設定される THEN the System SHALL PostgreSQL に接続する
2. WHEN `DB_TYPE=mysql` が設定される THEN the System SHALL MySQL に接続する（現行動作）
3. WHEN `DB_TYPE` が未設定である THEN the System SHALL デフォルトで MySQL に接続する
4. WHEN `DB_DSN` が設定される THEN the System SHALL 直接接続文字列として使用する（PostgreSQL 用）
5. WHEN 接続に失敗する THEN the System SHALL 適切なエラーメッセージを返す

### Requirement 2

**User Story:** 開発者として、DB 接続クライアントを interface 化したい。そうすることで、MySQL と PostgreSQL を交換可能にできる。

#### Acceptance Criteria

1. WHEN リポジトリが定義される THEN the System SHALL インターフェースを通じてデータアクセスを抽象化する
2. WHEN MySQL 実装が存在する THEN the System SHALL 既存コードをそのまま残す
3. WHEN PostgreSQL 実装が追加される THEN the System SHALL 同一インターフェースを実装する
4. WHEN DB_TYPE が切り替わる THEN the System SHALL main.go で適切なリポジトリ実装を注入する

### Requirement 3

**User Story:** 開発者として、PostgreSQL 用のリポジトリ実装を追加したい。そうすることで、Supabase 上でアプリが動作する。

#### Acceptance Criteria

1. WHEN PostgreSQL クエリが実行される THEN the System SHALL `$1, $2, ...` 形式のプレースホルダーを使用する
2. WHEN 日付計算が行われる THEN the System SHALL `DATE_SUB()` の代わりに `INTERVAL` 構文を使用する
3. WHEN 既存の SQL が移植される THEN the System SHALL 同一の結果セットを返す
4. WHEN CryptoRepository が使用される THEN the System SHALL PostgresCryptoRepository で同一インターフェースを実装する
5. WHEN TradeHistoryRepository が使用される THEN the System SHALL PostgresTradeHistoryRepository で同一インターフェースを実装する
6. WHEN OrderRepository が使用される THEN the System SHALL PostgresOrderRepository で同一インターフェースを実装する

### Requirement 4

**User Story:** 開発者として、API のレスポンスが変わらないことを保証したい。そうすることで、フロントエンドに影響を与えない。

#### Acceptance Criteria

1. WHEN PostgreSQL に切り替えた場合 THEN the System SHALL 全 API エンドポイントで同一レスポンスを返す
2. WHEN 既存のユニットテストを実行する THEN the System SHALL すべてパスする
3. WHEN PostgreSQL 用のテストを実行する THEN the System SHALL すべてパスする

## Connection Information

### AWS (MySQL) — 現行環境

```
Host: crypto-trading-db.cva64ye44jkh.ap-northeast-1.rds.amazonaws.com
Port: 1221
User: crypto_trading_root
Password: CryptoSl0S&mDGdY098
Database: crypto_trading
```

### Supabase (PostgreSQL) — 移行先

```
Project URL:      https://hastyhkkwedwqchtlsiv.supabase.co
DB Password:      AbcfbG294QB2hR
Publishable Key:  sb_publishable_ODStx-BX8rqpKFGPmlEAGw_jV4A8sAA
Direct Connection:
  postgresql://postgres.hastyhkkwedwqchtlsiv:AbcfbG294QB2hR@
    aws-1-ap-northeast-1.pooler.supabase.com:6543/postgres
```
