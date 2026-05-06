# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## フロントエンド概要

Nuxt 4 + Vue 3 + TypeScript + Tailwind CSS による暗号通貨取引Webアプリケーション。モックデータとバックエンドAPIの両方に対応。

## よく使うコマンド

```bash
make run          # 開発サーバー起動（API接続モード、--host 0.0.0.0）
make run-mock     # モックデータモードで起動
make test         # Vitest実行（変更後は必ず実行）
make build        # プロダクションビルド
make gen          # OpenAPIからTypeScript型定義生成 → types/api.ts
npm run test:watch  # テストウォッチモード
```

## アーキテクチャ

### ディレクトリ構成
```
pages/          → ルーティング（market, trade, portfolio, history）
composables/    → API呼び出し・状態管理ロジック
components/     → Vueコンポーネント
types/          → 型定義（api.ts はOpenAPI生成、手動編集禁止）
utils/          → モックデータ、エラーハンドラー
test/           → テストファイル
```

### Composables（主要ロジック）
- `useApi.ts` — APIクライアントラッパー
- `useCryptoData.ts` — 暗号通貨データ取得・管理
- `useTradeHistory.ts` / `useTradeHistoryApi.ts` — 取引履歴
- `useCurrentOrders.ts` — 未約定注文
- `useOrderForm.ts` / `useOrderData.ts` — 注文フォーム
- `useAutoRefresh.ts` — 自動更新（5秒間隔）
- `usePriceRounding.ts` — 価格フォーマット
- `useTimeFilter.ts` — 期間フィルター

### 環境変数
- `USE_MOCK_DATA` — `true`でモックデータ使用（デフォルト）、`false`でAPI接続
- `API_BASE_URL` — バックエンドURL（デフォルト: `http://localhost:8090`）

### テスト
- Vitest + `@vue/test-utils` + happy-dom
- fast-check によるプロパティベーステスト
- OpenAPI生成型の使用例: `import type { components } from '~/types/api'`

### 生成コード
`types/api.ts`はOpenAPIから自動生成。手動編集禁止。変更は`openapi.yaml`を編集して`make gen`を実行する。
