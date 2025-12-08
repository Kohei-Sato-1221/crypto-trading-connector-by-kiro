# Design Document

## Overview

暗号通貨取引Webアプリケーションのマーケットページは、Nuxt 4をベースとしたSPAとして実装します。モバイルファーストのアプローチで、ダークテーマのUIを採用し、リアルタイムの市場データを視覚的に表現します。初期段階ではモックデータを使用し、将来的なREST API統合を見据えた柔軟な設計を行います。

## Architecture

### Technology Stack

- **Framework**: Nuxt 4（Vue 3 + Composition API）
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Chart Library**: Chart.js with vue-chartjs（軽量で柔軟なチャート描画）
- **State Management**: Composables（Nuxtの推奨パターン）
- **Routing**: Nuxt File-based Routing

### Application Structure

```
crypto-trading-app/
├── pages/
│   ├── index.vue              # リダイレクト or ランディング
│   ├── market.vue             # マーケットページ（メイン実装）
│   ├── trade.vue              # 取引ページ（空ページ）
│   └── history.vue            # 履歴ページ（空ページ）
├── components/
│   ├── CryptoCard.vue         # 暗号通貨カードコンポーネント
│   ├── PriceChart.vue         # 価格チャートコンポーネント
│   ├── MarketHeader.vue       # ヘッダーコンポーネント
│   └── NavigationBar.vue      # ナビゲーションバーコンポーネント
├── composables/
│   ├── useCryptoData.ts       # 暗号通貨データ取得ロジック
│   ├── useTimeFilter.ts       # 時間フィルター状態管理
│   └── useAutoRefresh.ts      # 自動更新ロジック
├── types/
│   └── crypto.ts              # 型定義
├── utils/
│   └── mockData.ts            # モックデータ
└── assets/
    └── icons/                 # SVGアイコン
```

## Components and Interfaces

### 1. MarketHeader Component

ヘッダーコンポーネントは、アプリケーションタイトル、市場ステータス、時間フィルターを表示します。

**Props**: なし

**Emits**: 
- `filter-change`: 時間フィルターが変更されたときに発火（payload: TimeFilter）

**Features**:
- 固定ヘッダー（sticky positioning）
- 背景ブラー効果
- 時間フィルターボタン（24h、7d、30d、1y、All）
- 市場ステータスインジケーター

### 2. CryptoCard Component

個別の暗号通貨情報を表示するカードコンポーネントです。

**Props**:
```typescript
interface CryptoCardProps {
  crypto: CryptoData;
  timeFilter: TimeFilter;
}
```

**Features**:
- 通貨アイコン（円形背景）
- 通貨名とペア表示
- 現在価格（大きなフォント）
- 変動率バッジ（色分け：緑=上昇、赤=下降）
- 7日間チャート
- Tradeボタン

### 3. PriceChart Component

価格推移を表示するチャートコンポーネントです。

**Props**:
```typescript
interface PriceChartProps {
  data: ChartDataPoint[];
  isPositive: boolean;
  currency: string;
}
```

**Features**:
- エリアチャート（グラデーション塗りつぶし）
- 動的な色変更（緑/赤）
- X軸：曜日ラベル
- Y軸：価格範囲
- レスポンシブサイズ

### 4. NavigationBar Component

アプリケーション下部のナビゲーションバーです。

**Props**: なし

**Features**:
- 固定配置（bottom: 0）
- 4つのナビゲーションリンク
- アクティブ状態の視覚的フィードバック
- アイコン + ラベル

## Data Models

### CryptoData Type

```typescript
interface CryptoData {
  id: string;
  name: string;
  symbol: string;
  pair: string;
  icon: string;
  iconColor: string;
  currentPrice: number;
  changePercent: number;
  chartData: ChartDataPoint[];
}
```

### ChartDataPoint Type

```typescript
interface ChartDataPoint {
  day: string;
  price: number;
}
```

### TimeFilter Type

```typescript
type TimeFilter = '24h' | '7d' | '30d' | '1y' | 'all';
```

### NavigationItem Type

```typescript
interface NavigationItem {
  name: string;
  path: string;
  icon: string;
  label: string;
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Crypto card rendering completeness

*For any* CryptoData object with valid fields, rendering the CryptoCard component should display all required information (name, symbol, pair, price, changePercent, chart, Trade button)
**Validates: Requirements 1.2**

### Property 2: Price change color consistency

*For any* CryptoData object, if changePercent is positive, the change badge should be green with an up arrow; if negative, it should be red with a down arrow
**Validates: Requirements 1.3, 1.4**

### Property 3: Chart color consistency

*For any* CryptoData object, if changePercent is positive, the chart gradient should be green; if negative, it should be red
**Validates: Requirements 2.4, 2.5**

### Property 4: Time filter exclusivity

*For any* time filter selection, exactly one filter button should be in the active state (blue background), and all others should be inactive (gray background)
**Validates: Requirements 4.2, 4.3**

### Property 5: Navigation active state

*For any* current route, exactly one navigation item should be highlighted (blue color), matching the current page path
**Validates: Requirements 5.3**

### Property 6: Mock data structure validity

*For any* mock CryptoData object, it should conform to the CryptoData type definition with all required fields present and correctly typed
**Validates: Requirements 8.3, 8.4**

### Property 7: Price update maintains data integrity

*For any* CryptoData object after a price update, all required fields should remain valid (price > 0, changePercent is a number, chartData array is not empty)
**Validates: Requirements 7.2**

### Property 8: Auto-refresh timer cleanup

*For any* component using auto-refresh, when the component is unmounted, the timer should be properly cleared to prevent memory leaks
**Validates: Requirements 7.4**

## Error Handling

### Data Loading Errors

- モックデータ読み込み失敗時：空配列を返し、エラーメッセージを表示
- 将来のAPI統合時：try-catchでエラーをキャッチし、ユーザーフレンドリーなメッセージを表示

### Navigation Errors

- 存在しないルートへのアクセス：404ページまたは/marketへリダイレクト
- ナビゲーション失敗：エラーログを記録し、現在のページを維持

### Chart Rendering Errors

- チャートデータが不正な場合：エラーバウンダリでキャッチし、プレースホルダーを表示
- Chart.jsの初期化失敗：コンソールエラーを記録し、静的な代替表示

## Testing Strategy

### Unit Testing

- **Composables**: useCryptoData、useTimeFilterの単体テスト
- **Utility Functions**: モックデータ生成関数のテスト
- **Type Guards**: 型チェック関数のテスト

### Property-Based Testing

このプロジェクトでは、**fast-check**ライブラリを使用してプロパティベーステストを実装します。各プロパティテストは最低100回の反復実行を行い、ランダムな入力に対する正確性を検証します。

各プロパティベーステストには、設計書の対応するプロパティを明示するコメントを付与します：
```typescript
// Feature: crypto-market-page, Property 1: Crypto card rendering completeness
```

### Component Testing

- **CryptoCard**: 様々なpropsでのレンダリングテスト
- **PriceChart**: 正負のデータでの色変更テスト
- **MarketHeader**: フィルター変更イベントのテスト
- **NavigationBar**: ルート変更時のアクティブ状態テスト

### Integration Testing

- ページ全体のレンダリング
- ナビゲーション遷移
- フィルター変更時のUI更新

## Styling Guidelines

### Color Palette

```css
/* Background */
--bg-primary: #101922
--bg-secondary: #1c2630
--bg-tertiary: #2e3e50

/* Text */
--text-primary: #ffffff
--text-secondary: #94a3b8 (slate-400)

/* Accent */
--accent-primary: #137fec (blue)
--accent-success: #0bda5b (green)
--accent-danger: #fa6238 (red)

/* Bitcoin */
--bitcoin-orange: #f7931a

/* Ethereum */
--ethereum-purple: #627eea
```

### Typography

- **Font Family**: Manrope（Google Fonts）
- **Heading**: Bold, 18px-30px
- **Body**: Regular/Medium, 12px-14px
- **Small**: Medium, 10px

### Spacing

- **Card Padding**: 20px
- **Card Gap**: 16px
- **Component Margin**: 16px
- **Border Radius**: 8px-16px

## Auto-Refresh Mechanism

### Overview

マーケットページでは、5秒ごとに価格データを自動更新し、リアルタイムな市場の動きをシミュレートします。

### Implementation Strategy

**useAutoRefresh Composable**:
```typescript
export const useAutoRefresh = (callback: () => void, interval: number = 5000) => {
  const timerId = ref<NodeJS.Timeout | null>(null);

  const start = () => {
    if (timerId.value) return;
    timerId.value = setInterval(callback, interval);
  };

  const stop = () => {
    if (timerId.value) {
      clearInterval(timerId.value);
      timerId.value = null;
    }
  };

  onUnmounted(() => {
    stop();
  });

  return { start, stop };
};
```

### Mock Data Update Strategy

モックデータの価格更新は以下のロジックで実装します：

1. **価格変動**: 現在価格の±0.5%〜2%のランダムな変動
2. **変動率の再計算**: 新しい価格に基づいて変動率を更新
3. **チャートデータの更新**: 最新の価格ポイントを追加し、古いデータをシフト
4. **トランジション**: CSS transitionで価格変更を視覚的に表現

```typescript
const updateMockPrice = (crypto: CryptoData): CryptoData => {
  const changeRate = (Math.random() - 0.5) * 0.04; // -2% to +2%
  const newPrice = crypto.currentPrice * (1 + changeRate);
  const newChangePercent = crypto.changePercent + changeRate * 100;

  return {
    ...crypto,
    currentPrice: newPrice,
    changePercent: newChangePercent,
    chartData: updateChartData(crypto.chartData, newPrice)
  };
};
```

### UI Feedback

- **価格表示**: transition-all duration-300でスムーズに更新
- **変動率バッジ**: 色の変更もアニメーション付き
- **チャート**: 新しいデータポイントがスムーズに追加される

## Performance Considerations

- **Lazy Loading**: Chart.jsは動的インポートで必要時のみ読み込み
- **Virtual Scrolling**: 暗号通貨リストが増えた場合に備えて将来的に実装
- **Memoization**: チャートデータの計算結果をキャッシュ
- **Image Optimization**: SVGアイコンをインライン化してHTTPリクエストを削減
- **Timer Cleanup**: コンポーネントのアンマウント時に必ずタイマーをクリア
- **Throttling**: 短時間での複数更新を防ぐため、更新処理を適切に制御

## Future API Integration

### API Endpoint Design

```typescript
// GET /api/crypto/market
interface MarketResponse {
  data: CryptoData[];
  timestamp: number;
}

// GET /api/crypto/:id/chart?period=7d
interface ChartResponse {
  data: ChartDataPoint[];
  period: TimeFilter;
}
```

### Composable Design for API Integration

```typescript
// useCryptoData.ts
export const useCryptoData = () => {
  const USE_MOCK = true; // 環境変数で切り替え可能

  const fetchCryptoData = async (): Promise<CryptoData[]> => {
    if (USE_MOCK) {
      return getMockCryptoData();
    }
    
    // 将来のAPI実装
    const response = await $fetch('/api/crypto/market');
    return response.data;
  };

  return {
    fetchCryptoData
  };
};
```

## Accessibility

- **Semantic HTML**: 適切なHTML要素の使用（nav、button、article）
- **ARIA Labels**: アイコンボタンにaria-labelを追加
- **Keyboard Navigation**: すべてのインタラクティブ要素にキーボードアクセス
- **Color Contrast**: WCAG AA基準を満たすコントラスト比
- **Touch Targets**: 最小44x44pxのタッチターゲットサイズ
