# Design Document

## Overview

暗号通貨取引バックエンドAPIは、Golang、Echoフレームワーク、MVCアーキテクチャで実装します。OpenAPI仕様書に基づいてRESTful APIを提供し、フロントエンドに暗号通貨の市場データを配信します。初期段階ではモックリポジトリを使用し、将来的なデータベース統合を見据えた設計を行います。

## Architecture

### Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Echo v4
- **API Specification**: OpenAPI 3.0.3
- **Database (Future)**: PostgreSQL with SQLC
- **Testing**: Go standard testing package
- **Architecture Pattern**: MVC (Model-View-Controller)

### Project Structure

```
crypto-trading-connector-be/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── handler/
│   │   ├── crypto_handler.go      # HTTP handlers
│   │   └── crypto_handler_test.go # Handler tests
│   ├── service/
│   │   ├── crypto_service.go      # Business logic
│   │   └── crypto_service_test.go # Service tests
│   ├── repository/
│   │   ├── crypto_repository.go   # Data access (price_histories)
│   │   └── crypto_repository_test.go # Repository tests
│   ├── client/
│   │   ├── bitflyer_client.go     # bitFlyer Lightning API client
│   │   └── bitflyer_client_test.go # Client tests
│   └── model/
│       └── crypto.go               # Data models
├── pkg/
│   └── database/
│       └── db.go                   # MySQL database connection
├── test/
│   └── e2e/
│       └── crypto_e2e_test.go     # E2E tests (with build tag)
├── .env.example                    # Environment variables example
├── go.mod
├── go.sum
└── README.md
```

## Components and Interfaces

### 1. Model Layer

データ構造を定義します。OpenAPI仕様書のスキーマと一致させます。

```go
type ChartDataPoint struct {
    Day   string  `json:"day"`
    Price float64 `json:"price"`
}

type CryptoData struct {
    ID            string           `json:"id"`
    Name          string           `json:"name"`
    Symbol        string           `json:"symbol"`
    Pair          string           `json:"pair"`
    Icon          string           `json:"icon"`
    IconColor     string           `json:"iconColor"`
    CurrentPrice  float64          `json:"currentPrice"`
    ChangePercent float64          `json:"changePercent"`
    ChartData     []ChartDataPoint `json:"chartData"`
}

type MarketResponse struct {
    Data      []CryptoData `json:"data"`
    Timestamp int64        `json:"timestamp"`
}

type ChartResponse struct {
    Data   []ChartDataPoint `json:"data"`
    Period string           `json:"period"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}
```

### 2. Repository Layer

データアクセスを抽象化します。price_historiesテーブルから過去の価格データを取得します。

```go
type CryptoRepository interface {
    GetPriceHistories(productCode string, days int) ([]model.PriceHistory, error)
    GetDailyAveragePrices(productCode string, days int) ([]model.ChartDataPoint, error)
}

type MySQLCryptoRepository struct {
    db *sql.DB
}
```

**Data Source**:
- price_historiesテーブルから取得
- product_code: BTC_JPY, ETH_JPY
- 日毎の平均価格を計算（DATE(datetime)でグループ化、AVG(price)）
- 7日間のデータを取得

### 3. bitFlyer Client Layer

bitFlyer Lightning APIとの通信を担当します。

```go
type BitFlyerClient interface {
    GetTicker(productCode string) (*model.TickerResponse, error)
}

type BitFlyerClientImpl struct {
    baseURL string
    client  *http.Client
}
```

**API Endpoint**:
- GET https://api.bitflyer.com/v1/ticker?product_code={BTC_JPY|ETH_JPY}
- レスポンス: ltp（Last Traded Price）を使用

### 4. Service Layer

ビジネスロジックを実装します。bitFlyer APIとリポジトリを組み合わせてデータを構築します。

```go
type CryptoService interface {
    GetMarketData() (*model.MarketResponse, error)
    GetCryptoByID(id string) (*model.CryptoData, error)
    GetChartData(id string, period string) (*model.ChartResponse, error)
}

type CryptoServiceImpl struct {
    repo          repository.CryptoRepository
    bitflyerClient client.BitFlyerClient
}
```

**Responsibilities**:
- bitFlyer APIから現在価格を取得
- price_historiesから過去7日間の平均価格を取得
- 変動率の計算
- タイムスタンプの生成
- データの検証
- エラーハンドリング

### 5. Handler Layer

HTTPリクエスト/レスポンスを処理します。Echoのコンテキストを使用します。

```go
type CryptoHandler struct {
    service service.CryptoService
}

func (h *CryptoHandler) GetMarketData(c echo.Context) error
func (h *CryptoHandler) GetCryptoByID(c echo.Context) error
func (h *CryptoHandler) GetChartData(c echo.Context) error
```

**Responsibilities**:
- リクエストパラメータの抽出
- サービス層の呼び出し
- レスポンスの生成
- HTTPステータスコードの設定

## API Endpoints

### GET /api/v1/crypto/market

すべての暗号通貨の市場データを取得します。

**Response 200**:
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
        ...
      ]
    }
  ],
  "timestamp": 1704067200
}
```

### GET /api/v1/crypto/{id}

特定の暗号通貨の詳細データを取得します。

**Parameters**:
- `id` (path): 暗号通貨ID（例: bitcoin, ethereum）

**Response 200**:
```json
{
  "id": "bitcoin",
  "name": "Bitcoin",
  ...
}
```

**Response 404**:
```json
{
  "error": "NOT_FOUND",
  "message": "Cryptocurrency not found"
}
```

### GET /api/v1/crypto/{id}/chart

チャートデータを取得します。

**Parameters**:
- `id` (path): 暗号通貨ID
- `period` (query, optional): 期間（24h, 7d, 30d, 1y, all）、デフォルト: 7d

**Response 200**:
```json
{
  "data": [
    {"day": "Mon", "price": 9350000},
    ...
  ],
  "period": "7d"
}
```

## Error Handling

### Error Types

- `NOT_FOUND`: リソースが見つからない（404）
- `BAD_REQUEST`: 不正なリクエスト（400）
- `INTERNAL_SERVER_ERROR`: サーバー内部エラー（500）

### Error Response Format

```json
{
  "error": "ERROR_TYPE",
  "message": "Human-readable error message"
}
```

### Implementation

```go
func handleError(c echo.Context, statusCode int, errorType string, message string) error {
    return c.JSON(statusCode, model.ErrorResponse{
        Error:   errorType,
        Message: message,
    })
}
```

## Middleware Configuration

### CORS Middleware

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"http://localhost:3000"},
    AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
    AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
}))
```

### Logger Middleware

```go
e.Use(middleware.Logger())
```

### Recover Middleware

```go
e.Use(middleware.Recover())
```

## Testing Strategy

### Unit Testing

各層を独立してテストします。

**Repository Tests**:
- モックデータの取得
- 存在しないIDのエラーハンドリング
- データの整合性

**Service Tests**:
- リポジトリのモック化
- ビジネスロジックの検証
- エラーハンドリング

**Handler Tests**:
- HTTPリクエスト/レスポンスのテスト
- ステータスコードの検証
- JSONレスポンスの検証

### E2E Testing

実際のDB・APIとの統合をテストします。

**E2E Tests**:
- 実際のRDS MySQLに接続
- 実際のbitFlyer Lightning APIに接続
- サーバーを起動してHTTPリクエストを送信
- 正常系のレスポンスを検証

**Test Structure**:
```go
//go:build e2e

package e2e

import (
    "testing"
    "net/http"
    "encoding/json"
)

func TestGetMarketData(t *testing.T) {
    // 実際のサーバーに接続
    resp, err := http.Get("http://localhost:8080/api/v1/crypto/market")
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // レスポンスの検証
    var result model.MarketResponse
    json.NewDecoder(resp.Body).Decode(&result)
    assert.NotEmpty(t, result.Data)
    assert.NotZero(t, result.Timestamp)
}
```

**実行方法**:
```bash
# E2Eテストのみ実行
go test -tags=e2e ./test/e2e/...

# ユニットテストのみ実行（E2Eを除外）
go test ./...
```

### Test Coverage

- 目標: 80%以上のカバレッジ（ユニットテスト）
- 重要な機能は100%カバー
- E2Eテストは正常系のみ

### Testing Tools

```go
import (
    "testing"
    "net/http/httptest"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)
```

## Future Enhancements

### Database Integration

1. **PostgreSQL Setup**:
   - Docker Composeでローカル環境構築
   - マイグレーションツール（golang-migrate）

2. **SQLC Integration**:
   ```yaml
   # sqlc.yaml
   version: "2"
   sql:
     - schema: "scripts/schema.sql"
       queries: "scripts/queries.sql"
       engine: "postgresql"
       gen:
         go:
           package: "db"
           out: "internal/db"
   ```

3. **Database Schema**:
   ```sql
   CREATE TABLE cryptocurrencies (
       id VARCHAR(50) PRIMARY KEY,
       name VARCHAR(100) NOT NULL,
       symbol VARCHAR(10) NOT NULL,
       pair VARCHAR(20) NOT NULL,
       icon VARCHAR(10) NOT NULL,
       icon_color VARCHAR(7) NOT NULL,
       current_price DECIMAL(20, 2) NOT NULL,
       change_percent DECIMAL(10, 2) NOT NULL,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   CREATE TABLE chart_data (
       id SERIAL PRIMARY KEY,
       crypto_id VARCHAR(50) REFERENCES cryptocurrencies(id),
       day VARCHAR(10) NOT NULL,
       price DECIMAL(20, 2) NOT NULL,
       period VARCHAR(10) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

### Caching

- Redis統合
- 市場データのキャッシング（TTL: 5秒）
- チャートデータのキャッシング（TTL: 1分）

### Authentication

- JWT認証
- APIキー認証
- レート制限

### Monitoring

- Prometheus メトリクス
- ヘルスチェックエンドポイント
- ログ集約

## Performance Considerations

- **Response Time**: < 100ms (目標)
- **Concurrent Requests**: 1000+ req/s
- **Memory Usage**: < 100MB
- **Connection Pooling**: 将来のDB統合時に実装

## Security Considerations

- **CORS**: 許可されたオリジンのみ
- **Input Validation**: すべてのユーザー入力を検証
- **Error Messages**: 内部情報を漏らさない
- **HTTPS**: 本番環境では必須（将来）
