# OpenAPI Generated Code Integration Tasks

## 概要
生成されたコードを既存のコードベースに統合し、手動で定義された型定義を生成されたものに置き換えます。

## バックエンド統合タスク

- [x] 1. バックエンド: 生成された型を既存コードにインポート
  - internal/service/crypto_service.goで`internal/generated`をインポート
  - internal/handler/crypto_handler.goで`internal/generated`をインポート
  - internal/repository/crypto_repository.goで`internal/generated`をインポート
  - 既存の`model.CryptoData`等を`generated.CryptoData`に置き換え
  - 注意: 生成された型は`Id`（小文字始まり）を使用
  - _Files: service, handler, repository_

- [x] 2. バックエンド: internal/model/crypto.goの整理
  - 生成された型と重複する定義を削除（CryptoData, ChartDataPoint, MarketResponse, ChartResponse, ErrorResponse）
  - bitFlyer固有の型（PriceHistory, TickerResponse）は残す
  - _Files: internal/model/crypto.go_

- [x] 3. バックエンド: テストとビルド確認
  - `make build`でコンパイルエラーがないことを確認
  - `make test`で全テストが通ることを確認
  - 必要に応じて型の不一致を修正
  - _Verification: make build && make test_

## フロントエンド統合タスク

- [x] 4. フロントエンド: 生成された型を既存コードにインポート
  - types/crypto.tsで`types/api`から型をインポート
  - 既存の`CryptoData`を`components['schemas']['CryptoData']`に置き換え
  - 既存の`ChartDataPoint`を`components['schemas']['ChartDataPoint']`に置き換え
  - 生成された型を再エクスポート（利便性のため）
  - _Files: types/crypto.ts_

- [x] 5. フロントエンド: コンポーネントの型更新
  - types/crypto.tsで型を再エクスポートしたため、既存コードは変更不要
  - composables/useCryptoData.ts、components/CryptoCard.vue、pages/market.vue、utils/mockData.tsは既存のimportをそのまま使用
  - _Files: 変更不要（型の再エクスポートにより互換性維持）_

- [x] 6. フロントエンド: types/crypto.tsの整理
  - 生成された型と重複する定義を削除（CryptoData, ChartDataPoint）
  - フロントエンド固有の型（TimeFilter, NavigationItem）は残す
  - 生成された型を再エクスポート（利便性のため）
  - _Files: types/crypto.ts_

- [x] 7. フロントエンド: テスト確認
  - `make test`で全71テストが通ることを確認
  - 型エラーがないことを確認
  - 必要に応じてテストファイルの型参照を更新
  - _Verification: make test_

## 最終確認

- [x] 8. 統合テスト: 全体動作確認
  - バックエンド: `make build`成功 ✓
  - フロントエンド: `npm run build`成功 ✓
  - バックエンドテスト: 全テスト通過 ✓
  - フロントエンドテスト: 全71テスト通過 ✓
  - _Verification: All builds and tests pass_

## 注意事項
- 各タスク完了後に必ずテストを実行
- 型の不一致がある場合は、生成された型に合わせて既存コードを修正
- フィールド名の違い（例: `Id` vs `ID`）に注意
- 生成された型は`components['schemas']['TypeName']`の形式でアクセス
