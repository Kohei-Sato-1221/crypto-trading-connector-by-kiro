# Crypto Trading Connector - Frontend

暗号通貨取引Webアプリケーションのフロントエンドです。Nuxt 4、TypeScript、Tailwind CSSで構築されています。

## 概要

このアプリケーションは、暗号通貨の市場データを表示するWebインターフェースを提供します。モックデータとバックエンドAPIの両方に対応しています。

### 主な機能

- 暗号通貨市場データの表示（BTC、ETH）
- リアルタイム価格更新（5秒ごと）
- 価格チャートの表示
- レスポンシブデザイン（モバイル対応）
- モックデータとAPI切り替え可能

## 技術スタック

- **フレームワーク**: Nuxt 4
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS
- **テスト**: Vitest + fast-check（プロパティベーステスト）

## セットアップ

### 前提条件

- Node.js 18以上
- npm、pnpm、yarn、またはbun

### インストール

1. 依存関係をインストール

```bash
npm install
# または
pnpm install
# または
yarn install
# または
bun install
```

2. 環境変数を設定（オプション）

`.env.example`を`.env`にコピーして、必要に応じて値を設定します：

```bash
cp .env.example .env
```

`.env`ファイルの内容：

```env
# Use mock data (true) or connect to backend API (false)
USE_MOCK_DATA=true

# Backend API base URL
API_BASE_URL=http://localhost:8080
```

## 使用方法

### 開発サーバーの起動

開発サーバーを`http://localhost:3000`で起動します：

```bash
npm run dev
# または
pnpm dev
# または
yarn dev
# または
bun run dev
```

### モックデータとAPIの切り替え

#### モックデータを使用する場合（デフォルト）

`.env`ファイルで以下のように設定：

```env
USE_MOCK_DATA=true
```

または環境変数なしで起動（デフォルトでモックデータを使用）

#### バックエンドAPIを使用する場合

1. バックエンドサーバーを起動（`crypto-trading-connector-be`）

```bash
cd ../crypto-trading-connector-be
make run
```

2. `.env`ファイルで以下のように設定：

```env
USE_MOCK_DATA=false
API_BASE_URL=http://localhost:8080
```

3. フロントエンドを起動

```bash
npm run dev
```

### テストの実行

```bash
npm run test
# または
pnpm test
# または
yarn test
# または
bun run test
```

### プロダクションビルド

```bash
npm run build
# または
pnpm build
# または
yarn build
# または
bun run build
```

### プロダクションプレビュー

```bash
npm run preview
# または
pnpm preview
# または
yarn preview
# または
bun run preview
```

## プロジェクト構造

```
crypto-trading-connector-fe/
├── app/                        # Nuxt 4 app directory
├── assets/
│   └── css/
│       └── main.css           # グローバルスタイル
├── components/
│   ├── CryptoCard.vue         # 暗号通貨カードコンポーネント
│   ├── MarketHeader.vue       # マーケットヘッダー
│   ├── NavigationBar.vue      # ナビゲーションバー
│   └── PriceChart.vue         # 価格チャート
├── composables/
│   └── useCryptoData.ts       # 暗号通貨データ取得composable
├── pages/
│   └── market.vue             # マーケットページ
├── types/
│   └── crypto.ts              # TypeScript型定義
├── utils/
│   └── mockData.ts            # モックデータ
├── test/                      # テストファイル
├── nuxt.config.ts             # Nuxt設定
├── tailwind.config.js         # Tailwind CSS設定
├── vitest.config.ts           # Vitest設定
└── README.md
```

## 主要コンポーネント

### MarketHeader

マーケットページのヘッダーを表示します。タイトルと説明を含みます。

### CryptoCard

個別の暗号通貨情報を表示するカードコンポーネントです。価格、変動率、チャートを含みます。

### PriceChart

価格推移を表示するチャートコンポーネントです。SVGで描画されます。

### NavigationBar

アプリケーションのナビゲーションバーです。

## API統合

### useCryptoData Composable

暗号通貨データを取得するためのcomposableです。モックデータとAPIの両方に対応しています。

```typescript
const { cryptoData, loading, error, fetchCryptoData, refresh, useMockData } = useCryptoData()

// データを取得
await fetchCryptoData()

// データを更新
await refresh()
```

### 環境変数

- `USE_MOCK_DATA`: モックデータを使用するかどうか（`true`/`false`）
- `API_BASE_URL`: バックエンドAPIのベースURL（デフォルト: `http://localhost:8080`）

## テスト

### テスト戦略

- **ユニットテスト**: 各コンポーネントとcomposableをテスト
- **プロパティベーステスト**: fast-checkを使用してランダムなデータでテスト
- **モバイル最適化テスト**: レスポンシブデザインのテスト

### テストの実行

```bash
npm run test
```

### テストカバレッジ

```bash
npm run test:coverage
```

## トラブルシューティング

### APIに接続できない

- バックエンドサーバーが起動しているか確認してください
- `.env`ファイルの`API_BASE_URL`が正しいか確認してください
- CORSエラーが発生する場合は、バックエンドのCORS設定を確認してください

### モックデータが表示されない

- `.env`ファイルの`USE_MOCK_DATA`が`true`に設定されているか確認してください
- ブラウザのキャッシュをクリアしてください

## ライセンス

MIT License

## 作者

Crypto Trading Connector Team
