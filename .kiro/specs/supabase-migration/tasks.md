# Implementation Plan

## タスクの進め方
1. 実装エージェント
2. レビュワーエージェント

の２者でタスクを進めてください。
1が実装を進めて、完了した実装を2がレビューしてください。

このレビュー内容に関しては私がチェックします。私の判断を踏まえて1に修正をお願いするようなフローですすめたいと思います。私のレビューが完了したら、コミットをして次のタスクに進みます。
2の指摘事項がなかったら、私への確認事項はなしでコミットして次のタスクに進んでください。

## タスクの一覧

- [ ] 1. DB 接続層の PostgreSQL 対応
  - `pkg/database/db.go` に `ConnectPostgres()` を追加
  - `ConnectByType()` 関数を追加（DB_TYPE で切り替え）
  - `Config` 構造体に `DBType`, `DSN` フィールドを追加
  - `go.mod` に `github.com/lib/pq` を追加
  - `.env.example` に `DB_TYPE`, `DB_DSN` の設定例を追加
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.1_

- [ ] 2. CryptoRepository の PostgreSQL 実装
  - `internal/repository/crypto_repository_postgres.go` を新規作成
  - `PostgresCryptoRepository` 構造体を実装
  - `DATE_SUB(CURDATE(), INTERVAL ? DAY)` → `CURRENT_DATE - $1 * INTERVAL '1 day'`
  - プレースホルダーを `$1, $2` 形式に変更
  - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [ ] 3. TradeHistoryRepository の PostgreSQL 実装
  - `internal/repository/trade_history_repository_postgres.go` を新規作成
  - `PostgresTradeHistoryRepository` 構造体を実装
  - `DATE_SUB(NOW(), INTERVAL 7 DAY)` → `NOW() - INTERVAL '7 days'`
  - 動的クエリビルドのプレースホルダーカウンター管理
  - _Requirements: 3.1, 3.2, 3.3, 3.5_

- [ ] 4. OrderRepository の PostgreSQL 実装
  - `internal/repository/order_repository_postgres.go` を新規作成
  - `PostgresOrderRepository` 構造体を実装
  - プレースホルダーを `$1, $2, ...` 形式に変更
  - _Requirements: 3.1, 3.3, 3.6_

- [ ] 5. main.go の切り替えロジック実装
  - `DB_TYPE` に基づいて MySQL/PostgreSQL のリポジトリを切り替え
  - `ConnectByType()` を使用して DB 接続
  - `.env` を Supabase 向けに更新
  - _Requirements: 2.4_

- [ ] 6. テスト更新・実行・ドキュメント更新
  - PostgreSQL リポジトリのユニットテスト追加
  - 既存 MySQL テストが壊れていないことを確認（`make test`）
  - `CLAUDE.md` に PostgreSQL 対応の記述を追加
  - _Requirements: 4.1, 4.2, 4.3_
