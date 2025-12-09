# ルート統合: /order → /trade

## 変更内容

### 削除したファイル
- `pages/order.vue` - 削除（`pages/trade.vue` に統合）

### 更新したファイル
- `test/order-page.test.ts` - `TradePage` を参照するように更新
  - `import OrderPage from '~/pages/order.vue'` → `import TradePage from '~/pages/trade.vue'`
  - すべての `OrderPage` → `TradePage` に置換
  - すべての `/order` → `/trade` に置換

### 既存のリンク（変更不要）
以下のファイルは既に `/trade` を使用しているため、変更不要：
- `components/NavigationBar.vue` - ナビゲーションバーの Trade リンク
- `components/CryptoCard.vue` - Trade ボタンのリンク

## 理由
- `/trade` と `/order` の両方が存在していたが、内容は同じ
- `/trade` は既存のナビゲーションで使用されているため、こちらを採用
- `/order` は開発中に作成されたテスト用ページだったため削除

## 影響
- ユーザーは `/trade` ルートで注文ページにアクセス
- すべてのテストは `/trade` ルートを使用
- ナビゲーションとリンクは変更不要（既に `/trade` を使用）

## 実装済み機能（/trade ページ）
- 通貨ペア選択（BTC/JPY, ETH/JPY）
- 現在価格と変動率表示
- 時間範囲フィルター（1H, 24H, 7D, 30D, 1Y）
- 価格チャート
- 指値注文フォーム
- 価格入力と割引ボタン（99%, 97%, 95%）
- 数量入力と増減ボタン
- 残高表示と残高チェック
- 推定合計の計算
- 注文バリデーション
- レスポンシブデザイン
- アクセシビリティ対応
- 通貨ペア切り替え時のデフォルト値リセット

## テスト
- 146テスト全て合格
- プロパティベーステスト実装済み
- ユニットテスト実装済み
