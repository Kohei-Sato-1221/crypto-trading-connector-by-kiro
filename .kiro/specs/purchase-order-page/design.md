# 設計ドキュメント

## 概要

Purchase Order（注文）ページは、ユーザーが暗号通貨の指値注文を行うための専用ページです。このページは、リアルタイムの価格情報、過去の価格チャート、直感的な注文フォームを統合し、ユーザーが情報に基づいた取引判断を行えるようにします。

既存のNuxt 3 + Vue 3 + TypeScriptプロジェクトに統合され、既存のコンポーネントパターンとスタイリングシステムを活用します。

## アーキテクチャ

### 技術スタック

- **フレームワーク**: Nuxt 3 + Vue 3 (Composition API)
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS（既存のデザインシステムに準拠）
- **状態管理**: Vue Composition API (ref, computed, reactive)
- **チャートライブラリ**: Chart.js + vue-chartjs（既存のPriceChartコンポーネントと同様）
- **API通信**: Nuxt $fetch（既存のuseCryptoDataパターンに準拠）

### ページ構造

```
pages/
  order.vue (新規作成)
    ├─ OrderHeader (新規コンポーネント)
    ├─ PriceDisplay (新規コンポーネント)
    ├─ TimeFilterButtons (新規コンポーネント)
    ├─ PriceChart (既存コンポーネント再利用)
    ├─ OrderForm (新規コンポーネント)
    └─ NavigationBar (既存コンポーネント再利用)
```

## コンポーネントとインターフェース

### 1. OrderPage (pages/order.vue)

メインページコンポーネント。全体のレイアウトと状態管理を担当。

**責務:**
- 選択された通貨ペアの管理
- 選択された時間フィルターの管理
- 価格データとチャートデータの取得
- 子コンポーネント間のデータフロー調整

**Props:** なし（ページコンポーネント）

**State:**
```typescript
{
  selectedPair: 'BTC/JPY' | 'ETH/JPY'
  selectedTimeFilter: '1H' | '24H' | '7D' | '30D' | '1Y'
  currentPrice: number
  priceChange: number
  chartData: ChartDataPoint[]
  availableBalance: number
  loading: boolean
  error: Error | null
}
```

### 2. OrderHeader (components/OrderHeader.vue)

通貨ペア選択タブを表示するヘッダーコンポーネント。

**Props:**
```typescript
{
  selectedPair: 'BTC/JPY' | 'ETH/JPY'
}
```

**Events:**
```typescript
{
  'update:selectedPair': (pair: 'BTC/JPY' | 'ETH/JPY') => void
}
```

### 3. PriceDisplay (components/PriceDisplay.vue)

現在価格と24時間の変動率を表示するコンポーネント。

**Props:**
```typescript
{
  currentPrice: number
  priceChange: number
  currency: string
}
```

### 4. TimeFilterButtons (components/TimeFilterButtons.vue)

時間範囲フィルターボタンを表示するコンポーネント。

**Props:**
```typescript
{
  selectedFilter: '1H' | '24H' | '7D' | '30D' | '1Y'
}
```

**Events:**
```typescript
{
  'update:selectedFilter': (filter: string) => void
}
```

### 5. OrderForm (components/OrderForm.vue)

注文入力フォームの主要コンポーネント。

**Props:**
```typescript
{
  selectedPair: 'BTC/JPY' | 'ETH/JPY'
  currentPrice: number
  availableBalance: number
}
```

**Events:**
```typescript
{
  'submit-order': (order: OrderData) => void
}
```

**Internal State:**
```typescript
{
  orderType: 'limit' // 固定
  price: number
  amount: number
  estimatedTotal: number
}
```

### 6. PriceInput (components/PriceInput.vue)

価格入力フィールドと割引ボタンを含むコンポーネント。

**Props:**
```typescript
{
  modelValue: number
  currentPrice: number
  currency: string
}
```

**Events:**
```typescript
{
  'update:modelValue': (value: number) => void
}
```

### 7. AmountInput (components/AmountInput.vue)

数量入力フィールドと増減ボタンを含むコンポーネント。

**Props:**
```typescript
{
  modelValue: number
  minAmount: number
  step: number
  symbol: string
}
```

**Events:**
```typescript
{
  'update:modelValue': (value: number) => void
}
```

## データモデル

### OrderData

```typescript
interface OrderData {
  pair: 'BTC/JPY' | 'ETH/JPY'
  orderType: 'limit'
  price: number
  amount: number
  estimatedTotal: number
  timestamp: string
}
```

### PriceData

```typescript
interface PriceData {
  pair: string
  currentPrice: number
  priceChange: number
  changePercent: number
  timestamp: string
}
```

### ChartDataPoint

```typescript
interface ChartDataPoint {
  timestamp: string
  price: number
}
```

### OrderFormState

```typescript
interface OrderFormState {
  price: number
  amount: number
  estimatedTotal: number
}
```

## Composables

### useOrderData

注文ページのデータ取得と状態管理を担当するComposable。

```typescript
export const useOrderData = (pair: Ref<string>) => {
  const currentPrice = ref<number>(0)
  const priceChange = ref<number>(0)
  const chartData = ref<ChartDataPoint[]>([])
  const availableBalance = ref<number>(0)
  const loading = ref<boolean>(false)
  const error = ref<Error | null>(null)

  const fetchPriceData = async () => {
    // 価格データを取得
  }

  const fetchChartData = async (timeFilter: string) => {
    // チャートデータを取得
  }

  const fetchBalance = async () => {
    // 残高を取得
  }

  return {
    currentPrice,
    priceChange,
    chartData,
    availableBalance,
    loading,
    error,
    fetchPriceData,
    fetchChartData,
    fetchBalance
  }
}
```

### usePriceRounding

価格の丸め処理を担当するComposable。

```typescript
export const usePriceRounding = () => {
  const roundPrice = (price: number, pair: string): number => {
    // BTC: 百万円単位で丸める
    // ETH: 万円単位で丸める
    if (pair === 'BTC/JPY') {
      return Math.round(price / 1000000) * 1000000
    } else if (pair === 'ETH/JPY') {
      return Math.round(price / 10000) * 10000
    }
    return price
  }

  const calculateDiscountPrice = (price: number, percentage: number): number => {
    return Math.floor(price * percentage)
  }

  return {
    roundPrice,
    calculateDiscountPrice
  }
}
```

### useOrderForm

注文フォームのロジックを担当するComposable。

```typescript
export const useOrderForm = (
  pair: Ref<string>,
  currentPrice: Ref<number>,
  availableBalance: Ref<number>
) => {
  const price = ref<number>(0)
  const amount = ref<number>(0)
  const estimatedTotal = computed(() => price.value * amount.value)

  const initializeDefaults = () => {
    // デフォルト値を設定
    const { roundPrice } = usePriceRounding()
    price.value = roundPrice(currentPrice.value, pair.value)
    amount.value = pair.value === 'BTC/JPY' ? 0.001 : 0.01
  }

  const incrementAmount = () => {
    const step = pair.value === 'BTC/JPY' ? 0.001 : 0.01
    amount.value = parseFloat((amount.value + step).toFixed(3))
  }

  const decrementAmount = () => {
    const step = pair.value === 'BTC/JPY' ? 0.001 : 0.01
    const minAmount = step
    amount.value = Math.max(minAmount, parseFloat((amount.value - step).toFixed(3)))
  }

  const setDiscountPrice = (percentage: number) => {
    const { calculateDiscountPrice } = usePriceRounding()
    price.value = calculateDiscountPrice(currentPrice.value, percentage / 100)
  }

  const validateOrder = (): boolean => {
    if (price.value <= 0) return false
    if (amount.value <= 0) return false
    if (estimatedTotal.value > availableBalance.value) return false
    return true
  }

  return {
    price,
    amount,
    estimatedTotal,
    initializeDefaults,
    incrementAmount,
    decrementAmount,
    setDiscountPrice,
    validateOrder
  }
}
```


## 正確性プロパティ

*プロパティとは、システムのすべての有効な実行において真であるべき特性または動作です。本質的には、システムが何をすべきかについての形式的な記述です。プロパティは、人間が読める仕様と機械で検証可能な正確性保証との橋渡しとして機能します。*

### プロパティ1: 通貨ペア選択の一貫性

*任意の*通貨ペア選択において、選択されたペアは価格表示、チャート、注文フォームのすべてで一貫して反映されなければならない

**検証: 要件 1.2**

### プロパティ2: 価格変動の符号と表示の一致

*任意の*価格変動データにおいて、変動率がプラスの場合は緑色と上向き矢印が表示され、マイナスの場合は赤色と下向き矢印が表示されなければならない

**検証: 要件 2.2, 2.3**

### プロパティ3: 時間フィルター選択とチャートデータの対応

*任意の*時間フィルター選択において、選択された時間範囲に対応するチャートデータが表示されなければならない

**検証: 要件 3.2**

### プロパティ4: 価格丸め処理の正確性

*任意の*現在価格において、BTC/JPYの場合は百万円単位、ETH/JPYの場合は万円単位で丸められた値がデフォルト価格として設定されなければならない

**検証: 要件 5.1**

### プロパティ5: 数量増減の正確性

*任意の*数量値において、増加ボタンをクリックすると通貨ペアに応じた最小単位（BTC: 0.001、ETH: 0.01）だけ増加し、減少ボタンをクリックすると同じ値だけ減少しなければならない（ただし最小値を下回らない）

**検証: 要件 5.5, 5.6, 5.7**

### プロパティ6: 割引価格計算の正確性

*任意の*現在価格と割引率（99%, 97%, 95%）において、割引ボタンをクリックすると現在価格に割引率を乗じた値が価格フィールドに設定されなければならない

**検証: 要件 6.1, 6.2, 6.3**

### プロパティ7: 推定合計の計算正確性

*任意の*価格と数量の組み合わせにおいて、推定合計は価格×数量と等しくなければならない

**検証: 要件 7.2**

### プロパティ8: 残高超過時の注文防止

*任意の*注文データにおいて、推定合計が利用可能残高を超える場合、注文ボタンは無効化されなければならない

**検証: 要件 7.3**

### プロパティ9: 入力検証の完全性

*任意の*注文送信において、価格または数量が0以下の場合、または必須フィールドが空の場合、注文送信は防止されなければならない

**検証: 要件 8.1, 8.5**

## エラーハンドリング

### API通信エラー

- **価格データ取得失敗**: エラーメッセージを表示し、前回の価格データを保持
- **チャートデータ取得失敗**: エラーメッセージを表示し、空のチャートを表示
- **残高取得失敗**: エラーメッセージを表示し、注文ボタンを無効化
- **注文送信失敗**: 詳細なエラーメッセージを表示し、フォームの状態を保持

### 入力検証エラー

- **価格が0以下**: 「有効な価格を入力してください」と表示
- **数量が0以下**: 「有効な数量を入力してください」と表示
- **残高不足**: 「利用可能残高が不足しています」と表示し、注文ボタンを無効化
- **無効な数値入力**: 入力を拒否し、前の有効な値を保持

### ネットワークエラー

- **タイムアウト**: 「接続がタイムアウトしました。もう一度お試しください」と表示
- **接続エラー**: 「サーバーに接続できません。ネットワーク接続を確認してください」と表示

## テスト戦略

### ユニットテスト

各コンポーネントとComposableに対してユニットテストを実装します。

**テスト対象:**
- usePriceRounding: 価格の丸め処理と割引計算
- useOrderForm: 注文フォームのロジック（増減、検証）
- PriceInput: 価格入力と割引ボタンの動作
- AmountInput: 数量入力と増減ボタンの動作
- OrderForm: フォーム全体の統合動作

**テストフレームワーク:** Vitest（既存プロジェクトと同様）

### プロパティベーステスト

正確性プロパティを検証するためのプロパティベーステストを実装します。

**使用ライブラリ:** fast-check

**テスト設定:**
- 各プロパティテストは最低100回の反復を実行
- 各テストは設計ドキュメントの正確性プロパティを明示的に参照

**テスト例:**

```typescript
// Feature: purchase-order-page, Property 5: 数量増減の正確性
test('amount increment/decrement maintains correct step', () => {
  fc.assert(
    fc.property(
      fc.double({ min: 0.001, max: 10 }),
      fc.constantFrom('BTC/JPY', 'ETH/JPY'),
      (initialAmount, pair) => {
        const step = pair === 'BTC/JPY' ? 0.001 : 0.01
        const incremented = initialAmount + step
        const decremented = Math.max(step, initialAmount - step)
        
        expect(incremented).toBeCloseTo(initialAmount + step, 3)
        expect(decremented).toBeGreaterThanOrEqual(step)
      }
    ),
    { numRuns: 100 }
  )
})
```

### E2Eテスト

主要なユーザーフローをE2Eテストで検証します。

**テストシナリオ:**
1. ページ読み込みとデフォルト値の確認
2. 通貨ペア切り替えと価格更新の確認
3. 時間フィルター切り替えとチャート更新の確認
4. 注文フォームの入力と検証
5. 注文送信の成功フロー

**テストフレームワーク:** Playwright（推奨）

## スタイリングとデザイン

### カラーパレット

既存のデザインシステムに準拠：

- **背景色**: `#101922` (メイン), `#1c2936` (カード)
- **テキスト**: `#ffffff` (プライマリ), `#92adc9` (セカンダリ), `#64748b` (tertiary)
- **アクセント**: `#137fec` (プライマリボタン)
- **成功**: `#0bda5b` (プラス変動)
- **エラー**: `#fa6238` (マイナス変動、エラー)
- **ボーダー**: `#233648`, `#2e3e50`

### レスポンシブデザイン

- **モバイル**: 320px - 767px
- **タブレット**: 768px - 1023px
- **デスクトップ**: 1024px以上

### アクセシビリティ

- タッチターゲットは最低44x44px
- コントラスト比はWCAG AA基準を満たす
- キーボードナビゲーションをサポート
- スクリーンリーダー対応のARIAラベル

## パフォーマンス考慮事項

### データ取得の最適化

- 価格データは5秒ごとに自動更新（既存のautoRefreshパターンを使用）
- チャートデータは時間フィルター変更時のみ取得
- 残高データはページ読み込み時と注文送信後に取得

### レンダリング最適化

- computed プロパティを使用して不要な再計算を防ぐ
- v-memo ディレクティブを使用してチャートの再レンダリングを最適化
- 大きなリストには仮想スクロールを検討（将来的な拡張）

## セキュリティ考慮事項

### 入力検証

- すべての数値入力はクライアント側とサーバー側の両方で検証
- XSS攻撃を防ぐため、ユーザー入力をサニタイズ
- CSRFトークンを使用して注文送信を保護

### API通信

- HTTPS通信を強制
- 認証トークンをHTTP-onlyクッキーに保存
- レート制限を実装して過度なリクエストを防止

## バックエンド設計

### アーキテクチャ

既存のGo + Echo + PostgreSQLバックエンドに統合します。

**技術スタック:**
- **言語**: Go 1.21+
- **フレームワーク**: Echo v4
- **データベース**: PostgreSQL
- **ORM**: database/sql（標準ライブラリ）
- **外部API**: bitFlyer API（価格データ取得）

### APIエンドポイント

#### 1. 注文作成エンドポイント

```
POST /api/v1/orders
```

**リクエストボディ:**
```json
{
  "pair": "BTC/JPY",
  "orderType": "limit",
  "price": 14000000,
  "amount": 0.001
}
```

**レスポンス（成功）:**
```json
{
  "orderId": "uuid-string",
  "pair": "BTC/JPY",
  "orderType": "limit",
  "price": 14000000,
  "amount": 0.001,
  "estimatedTotal": 14000,
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**レスポンス（エラー）:**
```json
{
  "error": "INSUFFICIENT_BALANCE",
  "message": "Available balance is insufficient for this order"
}
```

#### 2. 残高取得エンドポイント

```
GET /api/v1/balance
```

**レスポンス:**
```json
{
  "availableBalance": 1540200,
  "currency": "JPY",
  "timestamp": 1704067200
}
```

#### 3. 注文履歴取得エンドポイント（将来的な拡張）

```
GET /api/v1/orders?limit=10&offset=0
```

### データモデル

#### Order（注文）テーブル

```sql
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  pair VARCHAR(20) NOT NULL,
  order_type VARCHAR(20) NOT NULL,
  price DECIMAL(20, 2) NOT NULL,
  amount DECIMAL(20, 8) NOT NULL,
  estimated_total DECIMAL(20, 2) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
```

#### Balance（残高）テーブル

```sql
CREATE TABLE balances (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL UNIQUE,
  available_balance DECIMAL(20, 2) NOT NULL DEFAULT 0,
  currency VARCHAR(10) NOT NULL DEFAULT 'JPY',
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_balances_user_id ON balances(user_id);
```

### レイヤー構造

#### 1. Handler層（internal/handler/order_handler.go）

HTTPリクエストの処理とレスポンスの返却を担当。

```go
type OrderHandler struct {
    service service.OrderService
}

func (h *OrderHandler) CreateOrder(c echo.Context) error
func (h *OrderHandler) GetBalance(c echo.Context) error
```

#### 2. Service層（internal/service/order_service.go）

ビジネスロジックの実装を担当。

```go
type OrderService interface {
    CreateOrder(userID string, req CreateOrderRequest) (*Order, error)
    GetBalance(userID string) (*Balance, error)
    ValidateOrder(req CreateOrderRequest, balance float64) error
}

type OrderServiceImpl struct {
    orderRepo   repository.OrderRepository
    balanceRepo repository.BalanceRepository
}
```

**ビジネスロジック:**
- 注文の検証（価格、数量、残高チェック）
- 推定合計の計算
- 残高の更新（トランザクション処理）
- 注文ステータスの管理

#### 3. Repository層（internal/repository/order_repository.go）

データベースアクセスを担当。

```go
type OrderRepository interface {
    Create(order *Order) error
    GetByID(id string) (*Order, error)
    GetByUserID(userID string, limit, offset int) ([]Order, error)
}

type BalanceRepository interface {
    GetByUserID(userID string) (*Balance, error)
    UpdateBalance(userID string, amount float64) error
}
```

### バリデーション

#### 注文バリデーションルール

1. **価格検証**
   - 価格 > 0
   - 価格は数値型
   - 小数点以下の桁数制限

2. **数量検証**
   - 数量 > 0
   - BTC: 最小 0.001、最大 100
   - ETH: 最小 0.01、最大 1000
   - 小数点以下の桁数制限

3. **残高検証**
   - 推定合計 ≤ 利用可能残高
   - 残高は正の値

4. **通貨ペア検証**
   - サポートされている通貨ペア（BTC/JPY, ETH/JPY）のみ

### エラーハンドリング

#### エラーコード

- `BAD_REQUEST`: 無効なリクエストパラメータ
- `INSUFFICIENT_BALANCE`: 残高不足
- `INVALID_AMOUNT`: 無効な数量
- `INVALID_PRICE`: 無効な価格
- `UNSUPPORTED_PAIR`: サポートされていない通貨ペア
- `INTERNAL_SERVER_ERROR`: サーバー内部エラー

### トランザクション管理

注文作成時は以下の操作をトランザクション内で実行：

1. 残高の確認
2. 注文レコードの作成
3. 残高の更新（減算）

いずれかの操作が失敗した場合、すべてロールバック。

### セキュリティ

#### 認証・認可

- JWTトークンによる認証
- ユーザーIDはトークンから取得
- 他のユーザーの残高や注文にアクセス不可

#### レート制限

- 注文作成: 1ユーザーあたり10リクエスト/分
- 残高取得: 1ユーザーあたり30リクエスト/分

#### 入力サニタイゼーション

- すべての入力値をバリデーション
- SQLインジェクション対策（プリペアドステートメント使用）
- XSS対策（出力時のエスケープ）

## 実装フェーズ

### フェーズ1: フロントエンド（モックデータ）

1. ページとコンポーネントの作成
2. モックデータを使用した動作確認
3. スタイリングとレスポンシブ対応
4. ユニットテストとプロパティベーステストの実装

### フェーズ2: API設計

1. OpenAPI仕様書の作成
2. 必要なエンドポイントの定義
3. リクエスト/レスポンスモデルの定義
4. エラーレスポンスの定義

### フェーズ3: バックエンド実装

1. データベーススキーマの作成（マイグレーション）
2. Repository層の実装
3. Service層の実装（ビジネスロジック）
4. Handler層の実装（APIエンドポイント）
5. ユニットテストの実装

### フェーズ4: 統合

1. フロントエンドとバックエンドの接続
2. E2Eテストの実装
3. エラーハンドリングの改善
4. パフォーマンステストと最適化
