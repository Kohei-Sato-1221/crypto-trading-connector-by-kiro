# Implementation Plan

- [x] 1. バックエンド: oapi-codegenのセットアップ
  - oapi-codegenをインストール
  - internal/generatedディレクトリを作成
  - Makefileにgenerateコマンドを追加
  - _Requirements: 1.1, 1.2, 1.3_

- [x] 2. バックエンド: コード生成の実行とテスト
  - make generateを実行
  - 生成されたmodels.goを確認
  - go buildでコンパイルエラーがないことを確認
  - _Requirements: 1.4_

- [x] 3. フロントエンド: openapi-typescriptのセットアップ
  - openapi-typescriptをインストール（npm）
  - types/ディレクトリを作成
  - Makefileにgenerateコマンドを追加
  - _Requirements: 2.1, 2.2, 2.3_

- [x] 4. フロントエンド: 型定義の生成とテスト
  - make generateを実行
  - 生成されたapi.tsを確認
  - TypeScriptのコンパイルエラーがないことを確認
  - _Requirements: 2.4_

- [x] 5. 統合: ルートMakefileの作成
  - ルートディレクトリにMakefileを作成
  - generate-allコマンドを追加（バックエンド + フロントエンド）
  - _Requirements: 3.1, 3.2_

- [x] 6. ドキュメント: README更新
  - バックエンドREADMEにコード生成の説明を追加
  - フロントエンドREADMEにコード生成の説明を追加
  - 使用例を追加
  - _Requirements: 5.1, 5.2, 5.3_

- [x] 7. 検証: 既存コードとの互換性確認
  - 既存のテストが動作することを確認
  - 生成された型と既存の型を比較
  - 必要に応じて既存コードを調整
  - _Requirements: 4.1, 4.2_

- [x] 8. Final Checkpoint
  - すべてのコマンドが動作することを確認
  - ドキュメントが完全であることを確認
  - 質問があればユーザーに確認


## 統合タスク（生成されたコードを既存コードに適用）

- [x] 9. バックエンド: 生成された型を既存コードに統合
  - service, handler, repositoryで`internal/generated`の型を使用
  - `model.CryptoData` → `generated.CryptoData`に置き換え
  - internal/model/crypto.goから重複する型定義を削除
  - bitFlyer固有の型（PriceHistory, TickerResponse）は保持
  - _Files: internal/service/crypto_service.go, internal/handler/crypto_handler.go, internal/repository/crypto_repository.go, internal/model/crypto.go_

- [x] 10. フロントエンド: 生成された型を既存コードに統合
  - types/crypto.tsで生成された型を再エクスポート
  - `CryptoData`, `ChartDataPoint`を`components['schemas']['*']`から取得
  - フロントエンド固有の型（TimeFilter, NavigationItem）は保持
  - 既存コードは型の再エクスポートにより変更不要
  - _Files: types/crypto.ts_

- [x] 11. 統合テスト: ビルドとテストの確認
  - バックエンド: `make build`成功
  - フロントエンド: `npm run build`成功
  - バックエンドテスト: 全テスト通過
  - フロントエンドテスト: 全71テスト通過
  - _Verification: All builds and tests pass_

## 完了
すべてのタスクが完了しました。OpenAPIから生成されたコードが既存のコードベースに統合され、正常に動作しています。
