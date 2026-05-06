# Design Document

## Overview

MySQL (AWS RDS) で動作しているバックエンドを、Supabase (PostgreSQL) でも動作するように拡張する。既存のリポジトリインターフェースを活用し、PostgreSQL 用の実装を別ファイルで追加する。`DB_TYPE` 環境変数で接続先を切り替える。

## Architecture

### 切り替え方式

```
環境変数による自動判定:
  POSTGRES_URL が設定されている → PostgreSQL を使用
  MYSQL_URL のみ設定されている  → MySQL を使用
  両方設定されている            → PostgreSQL を優先
  両方空                       → panic（起動失敗）
```

### 変更対象レイヤー

```
cmd/server/main.go         → Config.DBType による分岐
pkg/database/db.go         → MYSQL_URL / POSTGRES_URL ベースの接続
internal/repository/       → PostgreSQL 実装ファイル追加（既存は維持）
```

Handler 層・Service 層は変更なし（インターフェースで抽象化済み）。

## Components and Interfaces

### 1. DB 接続層 (pkg/database/db.go)

```go
type Config struct {
    MysqlURL    string // MYSQL_URL 環境変数
    PostgresURL string // POSTGRES_URL 環境変数
    DBType      string // 自動判定結果: "mysql" or "postgres"
}

func LoadConfigFromEnv() *Config  // URL 有無で DBType を自動判定
func Connect(config *Config) (*sql.DB, error) // DBType に基づき接続
```

### 2. リポジトリ切り替え (main.go)

```go
dbConfig := database.LoadConfigFromEnv()
db, err := database.Connect(dbConfig)

var cryptoRepo repository.CryptoRepository
var orderRepo repository.OrderRepository
var tradeHistoryRepo repository.TradeHistoryRepository

if dbConfig.DBType == "postgres" {
    cryptoRepo = repository.NewPostgresCryptoRepository(db)
    orderRepo = repository.NewPostgresOrderRepository(db)
    tradeHistoryRepo = repository.NewPostgresTradeHistoryRepository(db)
} else {
    cryptoRepo = repository.NewMySQLCryptoRepository(db)
    orderRepo = repository.NewOrderRepository(db)
    tradeHistoryRepo = repository.NewMySQLTradeHistoryRepository(db)
}
```

### 4. MySQL → PostgreSQL SQL 変換

| MySQL | PostgreSQL |
|-------|-----------|
| `?` プレースホルダー | `$1, $2, ...` |
| `DATE_SUB(CURDATE(), INTERVAL ? DAY)` | `CURRENT_DATE - $1 * INTERVAL '1 day'` |
| `DATE_SUB(NOW(), INTERVAL 7 DAY)` | `NOW() - INTERVAL '7 days'` |
| `DATE(datetime)` | `DATE(datetime)`（互換） |
| `AVG(price)` | `AVG(price)`（互換） |
| `ROUND(expr, 2)` | `ROUND(expr::numeric, 2)` |
| `COALESCE(...)` | `COALESCE(...)`（互換） |
| `LIMIT ? OFFSET ?` | `LIMIT $N OFFSET $M` |

### 5. ファイル構成

```
internal/repository/
├── crypto_repository.go              # インターフェース + MySQL 実装（既存）
├── crypto_repository_postgres.go     # PostgreSQL 実装（新規）
├── trade_history_repository.go       # インターフェース + MySQL 実装（既存）
├── trade_history_repository_postgres.go # PostgreSQL 実装（新規）
├── order_repository.go               # インターフェース + MySQL 実装（既存）
└── order_repository_postgres.go      # PostgreSQL 実装（新規）
```

## PostgreSQL テーブル定義

```sql
-- buy_orders (bitflyer)
CREATE TABLE IF NOT EXISTS buy_orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(50) UNIQUE,
    product_code VARCHAR(50),
    side VARCHAR(20),
    price DOUBLE PRECISION,
    size DOUBLE PRECISION,
    exchange VARCHAR(50),
    status VARCHAR(100) DEFAULT 'UNFILLED',
    strategy SMALLINT NOT NULL DEFAULT 99,
    remarks TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- sell_orders (bitflyer)
CREATE TABLE IF NOT EXISTS sell_orders (
    id SERIAL PRIMARY KEY,
    parentid VARCHAR(50),
    order_id VARCHAR(50) UNIQUE,
    product_code VARCHAR(50),
    side VARCHAR(20),
    price DOUBLE PRECISION,
    size DOUBLE PRECISION,
    exchange VARCHAR(50),
    status VARCHAR(100) DEFAULT 'UNFILLED',
    remarks TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- price_histories
CREATE TABLE IF NOT EXISTS price_histories (
    id SERIAL PRIMARY KEY,
    datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    product_code VARCHAR(50) NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    price_ratio_24h DOUBLE PRECISION,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_product_code_datetime
    ON price_histories (product_code, datetime);

-- okj_buy_orders (OKJ)
CREATE TABLE IF NOT EXISTS okj_buy_orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(50) UNIQUE,
    pair VARCHAR(50),
    price DOUBLE PRECISION,
    size DOUBLE PRECISION,
    exchange VARCHAR(50),
    state INTEGER DEFAULT 0,
    sell_order_id VARCHAR(50),
    sell_order_state VARCHAR(20),
    sell_price DOUBLE PRECISION,
    side VARCHAR(20),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- updatetime 自動更新トリガー
CREATE OR REPLACE FUNCTION update_updatetime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updatetime = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_buy_orders_updatetime
    BEFORE UPDATE ON buy_orders
    FOR EACH ROW EXECUTE FUNCTION update_updatetime();
CREATE TRIGGER trg_sell_orders_updatetime
    BEFORE UPDATE ON sell_orders
    FOR EACH ROW EXECUTE FUNCTION update_updatetime();
CREATE TRIGGER trg_okj_buy_orders_updatetime
    BEFORE UPDATE ON okj_buy_orders
    FOR EACH ROW EXECUTE FUNCTION update_updatetime();
```

## Testing Strategy

- `go-sqlmock` は `database/sql` ベースなので MySQL/PostgreSQL 共通で使用可能
- 各 PostgreSQL リポジトリに対してユニットテストを追加
- 既存の MySQL テストは変更なし
- `make test` で全テストがパスすることを確認

## Verification

1. `make test` — 既存テスト + 新規テストがすべてパス
2. `.env` を Supabase 設定に変更し `make run` → 各 API エンドポイントの動作確認
3. `make curl a=market` 等で実際のレスポンス確認
