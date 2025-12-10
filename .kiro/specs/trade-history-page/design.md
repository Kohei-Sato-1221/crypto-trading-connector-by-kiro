# 設計書

## 概要

取引履歴ページは、ユーザーの暗号通貨取引パフォーマンスを包括的に表示するWebアプリケーション機能です。このページは、利益統計、約定数、詳細な取引ログを提供し、資産と時間によるフィルタリング機能を含みます。

## アーキテクチャ

### フロントエンド（Nuxt.js/Vue.js）
- **ページコンポーネント**: `/pages/history.vue` - メインの取引履歴ページ
- **コンポーザブル**: 取引データの取得と状態管理
- **コンポーネント**: 統計カード、取引ログ、フィルターUI

### バックエンド（Go/Echo）
- **ハンドラー**: HTTP リクエストの処理
- **サービス**: ビジネスロジックと利益計算
- **リポジトリ**: データベースアクセス層
- **モデル**: データ構造定義

### データベース
- **buy_orders**: 購入注文データ
- **sell_orders**: 売却注文データ（parentidで買い注文と関連付け）

## コンポーネントとインターフェース

### フロントエンドコンポーネント

#### 1. TradeHistoryPage (`/pages/history.vue`)
メインページコンポーネント
- 統計セクション（総利益、約定数）
- フィルターセクション（時間、資産）
- 取引ログセクション
- ページネーション

#### 2. TradeStatistics コンポーネント
取引統計の表示
```typescript
interface TradeStatistics {
  totalProfit: number
  profitPercentage: number
  executionCount: number
  period: string
}
```

#### 3. TransactionLog コンポーネント
取引履歴リストの表示
```typescript
interface Transaction {
  id: string
  cryptocurrency: string
  timestamp: string
  profit: number
  orderType: 'sell'
  orderId: string
  buyPrice: number
  sellPrice: number
  amount: number
}
```

#### 4. FilterControls コンポーネント
フィルター制御UI
```typescript
interface FilterState {
  timeFilter: 'all' | '7days'
  assetFilter: 'all' | 'BTC' | 'ETH'
}
```

### バックエンドインターフェース

#### 1. TradeHistoryHandler
HTTP リクエストハンドラー
```go
type TradeHistoryHandler struct {
    tradeHistoryService TradeHistoryService
}

func (h *TradeHistoryHandler) GetTradeStatistics(c echo.Context) error
func (h *TradeHistoryHandler) GetTransactionLog(c echo.Context) error
```

#### 2. TradeHistoryService
ビジネスロジック層
```go
type TradeHistoryService interface {
    GetTradeStatistics(filters TradeFilters) (*TradeStatistics, error)
    GetTransactionLog(filters TradeFilters, pagination Pagination) (*TransactionLogResponse, error)
    CalculateProfit(sellOrder SellOrder, buyOrder BuyOrder) float64
}
```

#### 3. TradeHistoryRepository
データアクセス層
```go
type TradeHistoryRepository interface {
    GetCompletedTransactions(filters TradeFilters, pagination Pagination) ([]Transaction, error)
    GetTotalProfit(filters TradeFilters) (float64, error)
    GetExecutionCount(filters TradeFilters) (int, error)
}
```

## データモデル

### API レスポンスモデル

#### TradeStatistics
```go
type TradeStatistics struct {
    TotalProfit      float64 `json:"totalProfit"`      // 小数点第1位まで、第2位以下四捨五入
    ProfitPercentage float64 `json:"profitPercentage"` // 小数点第1位まで、第2位以下四捨五入
    ExecutionCount   int     `json:"executionCount"`
    Period          string  `json:"period"`
}
```

#### Transaction
```go
type Transaction struct {
    ID             string    `json:"id"`
    Cryptocurrency string    `json:"cryptocurrency"`
    Timestamp      time.Time `json:"timestamp"`
    Profit         float64   `json:"profit"`         // 小数点第1位まで、第2位以下四捨五入
    OrderType      string    `json:"orderType"`
    OrderID        string    `json:"orderId"`
    BuyPrice       float64   `json:"buyPrice"`
    SellPrice      float64   `json:"sellPrice"`
    Amount         float64   `json:"amount"`
    BuyOrderID     string    `json:"buyOrderId"`
}
```

#### TransactionLogResponse
```go
type TransactionLogResponse struct {
    Transactions []Transaction `json:"transactions"`
    HasMore      bool         `json:"hasMore"`
    Total        int          `json:"total"`
}
```

### データベースクエリ

#### 利益計算クエリ
```sql
SELECT 
    s.id,
    s.order_id as sell_order_id,
    s.product_code,
    s.price as sell_price,
    s.size as sell_size,
    s.timestamp as sell_timestamp,
    b.order_id as buy_order_id,
    b.price as buy_price,
    b.size as buy_size,
    b.timestamp as buy_timestamp,
    (s.price * s.size) - (b.price * b.size) as profit
FROM sell_orders s
JOIN buy_orders b ON s.parentid = b.order_id
WHERE s.status = 'FILLED'
ORDER BY s.timestamp DESC
```

## Correctness Properties

*プロパティとは、システムのすべての有効な実行において真であるべき特性や動作のことです。本質的に、システムが何をすべきかについての形式的な記述です。プロパティは、人間が読める仕様と機械で検証可能な正確性保証の橋渡しとなります。*

### プロパティ 1: 利益計算の正確性
*任意の* 売り注文と対応する買い注文のペアについて、計算された利益は (売却価格 * 売却数量) - (購入価格 * 購入数量) と等しくなければならない
**検証対象: 要件 1.2, 2.4**

### プロパティ 2: 統計の一貫性
*任意の* フィルター条件について、表示される総利益は個別取引利益の合計と等しく、約定数は完了した売り注文の数と等しくなければならない
**検証対象: 要件 1.1, 1.4**

### プロパティ 3: 資産フィルタリングの正確性
*任意の* 資産フィルター（BTC、ETH、全資産）について、返される取引はすべて選択された暗号通貨に一致し、統計も対応して更新されなければならない
**検証対象: 要件 3.2, 3.3, 3.4**

### プロパティ 4: 時間フィルタリングの正確性
*任意の* 時間フィルター（全期間、過去7日間）について、返される取引はすべて指定された時間範囲内にあり、統計も対応して更新されなければならない
**検証対象: 要件 4.2, 4.3, 4.4, 4.5**

### プロパティ 5: 取引情報の完全性
*任意の* 表示される取引について、暗号通貨名、タイムスタンプ、利益額、注文ID、購入価格、売却価格、数量のすべての必要な情報が含まれなければならない
**検証対象: 要件 2.2, 2.3**

### プロパティ 6: データ関連付けの整合性
*任意の* 表示される取引について、売り注文に対応する買い注文がparentidを通じて正しく関連付けられなければならない
**検証対象: 要件 2.1**

### プロパティ 7: ページネーション時のフィルター維持
*任意の* ページネーション操作について、現在のフィルター設定（資産と時間）が維持され、追加データが正しい条件で取得されなければならない
**検証対象: 要件 6.2, 6.3**

### プロパティ 8: 数値フォーマットの正確性
*任意の* 金額値について、日本円として適切な千の位区切り文字でフォーマットされ、利益額は小数点第1位まで表示し小数点第2位以下は四捨五入されなければならない
**検証対象: 要件 1.5**

### プロパティ 9: タイムスタンプフォーマットの正確性
*任意の* タイムスタンプについて、ユーザーフレンドリーな形式（「今日、14:30」、「昨日、09:15」など）でフォーマットされなければならない
**検証対象: 要件 2.5**

## エラーハンドリング

### フロントエンドエラーハンドリング
1. **ネットワークエラー**: 再試行オプション付きのユーザーフレンドリーなメッセージ
2. **データ読み込みエラー**: 部分的なデータ表示と明確なエラー表示
3. **フィルターエラー**: デフォルト状態への復帰

### バックエンドエラーハンドリング
1. **データベース接続エラー**: 適切なHTTPステータスコードとエラーレスポンス
2. **データ不整合エラー**: ログ記録と代替データ提供
3. **計算エラー**: 安全なデフォルト値の提供

## テスト戦略

### 単体テスト
- 利益計算ロジックの検証
- フィルタリング機能の検証
- データ変換ロジックの検証
- エラーハンドリングの検証

### プロパティベーステスト
- **テストライブラリ**: Go用のgopter、TypeScript用のfast-check
- **実行回数**: 各プロパティテストは最低100回実行
- **テストタグ**: 各プロパティベーステストは対応する設計書のプロパティを明示的に参照

#### プロパティベーステストの実装要件
1. 各正確性プロパティは単一のプロパティベーステストで実装される
2. テストは以下の形式でタグ付けされる: '**Feature: trade-history-page, Property {number}: {property_text}**'
3. ランダムなテストデータ生成器を使用して幅広い入力をカバー
4. 実際の機能を検証し、モックやフェイクデータは使用しない

### 統合テスト
- API エンドポイントのテスト
- データベースクエリのテスト
- フロントエンド・バックエンド統合のテスト

### E2Eテスト
- ユーザーフローの完全なテスト
- フィルタリングとページネーションのテスト
- レスポンシブデザインのテスト