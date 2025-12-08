# Implementation Plan

- [x] 1. プロジェクトセットアップ
  - Go modulesの初期化
  - 必要な依存関係のインストール（Echo、lib/pq）
  - ディレクトリ構造の作成
  - .gitignoreの作成
  - _Requirements: 5.1, 6.1_

- [x] 2. OpenAPI仕様書の作成
  - openapi.yamlをルートディレクトリに作成
  - エンドポイント定義（/crypto/market, /crypto/{id}, /crypto/{id}/chart）
  - スキーマ定義（CryptoData, ChartDataPoint, MarketResponse, ChartResponse, ErrorResponse）
  - フロントエンドの型定義と一致させる
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [x] 3. モデル層の実装
  - internal/model/crypto.goを作成
  - ChartDataPoint構造体を定義
  - CryptoData構造体を定義
  - MarketResponse構造体を定義
  - ChartResponse構造体を定義
  - ErrorResponse構造体を定義
  - JSONタグを追加
  - _Requirements: 5.5_

- [x] 4. データベース接続の実装
  - pkg/database/db.goを作成
  - 環境変数から接続情報を読み込む
  - MySQL接続プールを作成
  - 接続テスト関数を実装
  - _Requirements: 8.1, 8.2_

- [x] 5. 外部API統合（bitFlyer Lightning API）
  - internal/client/bitflyer_client.goを作成
  - BitFlyerClientインターフェースを定義
  - GetTicker()メソッドを実装（GET /v1/ticker）
  - BTC_JPYとETH_JPYの価格取得
  - エラーハンドリングとリトライロジック
  - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5_

- [-] 6. リポジトリ層の実装
  - internal/repository/crypto_repository.goを作成
  - CryptoRepositoryインターフェースを定義
  - MySQLCryptoRepository構造体を実装
  - GetPriceHistories()メソッドを実装（price_historiesテーブルから取得）
  - 日毎の平均価格計算（DATE(datetime)でグループ化、AVG(price)）
  - 7日間のデータ取得
  - product_codeでフィルタリング（BTC_JPY, ETH_JPY）
  - _Requirements: 8.3, 8.4, 2.3, 2.4, 4.2, 4.3_

- [ ] 6.1 リポジトリ層のユニットテストを作成
  - internal/repository/crypto_repository_test.goを作成
  - GetPriceHistories()のテスト（モックDB使用）
  - 日毎平均価格計算のテスト
  - product_codeフィルタリングのテスト
  - _Requirements: 9.1, 9.2_

- [-] 7. サービス層の実装
  - internal/service/crypto_service.goを作成
  - CryptoServiceインターフェースを定義
  - CryptoServiceImpl構造体を実装
  - GetMarketData()メソッドを実装
    - bitFlyer APIから現在価格を取得
    - price_historiesから過去7日間の平均価格を取得
    - 変動率を計算
    - タイムスタンプ生成
  - GetCryptoByID()メソッドを実装
  - GetChartData()メソッドを実装（デフォルトperiod処理）
  - エラーハンドリングを実装
  - _Requirements: 5.3, 2.1, 2.2, 2.3, 2.4, 2.5, 3.1, 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 7.1 サービス層のユニットテストを作成
  - internal/service/crypto_service_test.goを作成
  - モックリポジトリとモックbitFlyerクライアントを使用
  - GetMarketData()のテスト
  - GetCryptoByID()のテスト（正常系）
  - GetCryptoByID()のテスト（エラー系）
  - GetChartData()のテスト
  - _Requirements: 9.1, 9.3_

- [-] 8. ハンドラー層の実装
  - internal/handler/crypto_handler.goを作成
  - CryptoHandler構造体を実装
  - GetMarketData()ハンドラーを実装
  - GetCryptoByID()ハンドラーを実装（パスパラメータ抽出）
  - GetChartData()ハンドラーを実装（クエリパラメータ抽出）
  - エラーレスポンス処理を実装
  - _Requirements: 5.2, 2.1, 2.7, 3.1, 3.2, 4.1, 4.6, 10.1, 10.2, 10.3, 10.4_

- [ ] 8.1 ハンドラー層のユニットテストを作成
  - internal/handler/crypto_handler_test.goを作成
  - モックサービスを使用
  - GetMarketData()のHTTPテスト
  - GetCryptoByID()のHTTPテスト（正常系）
  - GetCryptoByID()のHTTPテスト（404エラー）
  - GetChartData()のHTTPテスト
  - JSONレスポンスの検証
  - _Requirements: 9.1, 9.3_

- [x] 9. メインサーバーの実装
  - cmd/server/main.goを作成
  - Echoインスタンスの初期化
  - ミドルウェアの設定（Logger、CORS、Recover）
  - ルーティングの設定
  - 依存性注入（Repository → Service → Handler）
  - サーバー起動処理（ポート8080）
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 7.1, 7.2, 7.3_

- [x] 10. CORS設定の実装
  - CORSミドルウェアの設定
  - localhost:3000からのアクセスを許可
  - 許可するHTTPメソッドの設定（GET, POST, PUT, DELETE, OPTIONS）
  - 許可するヘッダーの設定
  - _Requirements: 7.1, 7.2, 7.3_

- [x] 11. エラーハンドリングの実装
  - カスタムエラーハンドラーの実装
  - 404エラーのハンドリング
  - 500エラーのハンドリング
  - エラーログの記録
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [x] 12. README.mdの作成
  - プロジェクト概要
  - セットアップ手順
  - API使用方法
  - テスト実行方法
  - ディレクトリ構造の説明

- [x] 13. E2Eテストの実装
  - test/e2e/crypto_e2e_test.goを作成
  - ビルドタグ `//go:build e2e` を追加
  - 実際のRDS MySQLに接続
  - 実際のbitFlyer Lightning APIに接続
  - サーバーを起動してテスト
  - GET /api/v1/crypto/marketのE2Eテスト
    - HTTPステータスコード200を確認
    - dataフィールドが存在することを確認
    - timestampフィールドが存在することを確認
  - GET /api/v1/crypto/bitcoinのE2Eテスト
    - HTTPステータスコード200を確認
    - 必須フィールド（id, name, symbol, currentPrice等）の存在を確認
  - GET /api/v1/crypto/bitcoin/chartのE2Eテスト
    - HTTPステータスコード200を確認
    - dataフィールドが配列であることを確認
    - periodフィールドが存在することを確認
  - 環境変数から接続情報を読み込む
  - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5, 12.6, 12.7, 12.8, 12.9_

- [x] 14. フロントエンドとの統合
  - フロントエンドのuseCryptoData.tsを更新
  - モックデータからAPI呼び出しに切り替え
  - 環境変数でモック/API切り替えを実装
  - エラーハンドリングの実装
  - ローディング状態の実装

- [x] 15. Final Checkpoint - すべてのテストが通ることを確認
  - すべてのユニットテストが通ることを確認
  - 統合テストが成功することを確認
  - フロントエンドとバックエンドの連携を確認
  - 質問があればユーザーに確認
