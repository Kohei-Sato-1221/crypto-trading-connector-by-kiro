# Requirements Document

## Introduction

暗号通貨取引WebアプリケーションのバックエンドAPIを実装します。このAPIは、フロントエンド（crypto-trading-connector-fe）に暗号通貨の市場データを提供します。Golang、Echo フレームワーク、SQLC を使用し、MVCアーキテクチャで構築します。初期段階ではモックデータを使用し、将来的にデータベース統合を見据えた設計を行います。

## Glossary

- **System**: 暗号通貨取引バックエンドAPI
- **API**: RESTful APIエンドポイント
- **OpenAPI**: API仕様を定義する標準フォーマット
- **Echo**: Golangの高速Webフレームワーク
- **SQLC**: SQLからGolangコードを生成するツール
- **MVC**: Model-View-Controller アーキテクチャパターン
- **Mock Repository**: データベースの代わりにメモリ内データを使用するリポジトリ
- **CORS**: Cross-Origin Resource Sharing（クロスオリジンリソース共有）

## Requirements

### Requirement 1

**User Story:** 開発者として、OpenAPI仕様書を用意したい。そうすることで、APIの仕様が明確になり、フロントエンドとの連携がスムーズになる。

#### Acceptance Criteria

1. WHEN OpenAPI仕様書が作成される THEN the System SHALL OpenAPI 3.0.3形式で記述する
2. WHEN API仕様が定義される THEN the System SHALL すべてのエンドポイント、リクエスト、レスポンスを含める
3. WHEN スキーマが定義される THEN the System SHALL フロントエンドの型定義と一致する構造にする
4. WHEN 仕様書が配置される THEN the System SHALL ルートディレクトリ（openapi.yaml）に配置する

### Requirement 2

**User Story:** 開発者として、市場データを取得するAPIエンドポイントを実装したい。そうすることで、フロントエンドが暗号通貨の一覧を表示できる。

#### Acceptance Criteria

1. WHEN GET /api/v1/crypto/market がリクエストされる THEN the System SHALL すべての暗号通貨データを返す
2. WHEN 現在価格が取得される THEN the System SHALL bitFlyer Lightning APIから最新価格を取得する
3. WHEN 過去価格データが取得される THEN the System SHALL price_historiesテーブルから7日間のデータを取得する
4. WHEN 1日に複数の価格データが存在する THEN the System SHALL その日の平均価格を計算する
5. WHEN レスポンスが返される THEN the System SHALL JSON形式でデータとタイムスタンプを含める
6. WHEN データが取得される THEN the System SHALL BTC_JPY と ETH_JPY の情報を含める
7. WHEN エラーが発生する THEN the System SHALL 適切なHTTPステータスコードとエラーメッセージを返す

### Requirement 3

**User Story:** 開発者として、特定の暗号通貨の詳細を取得するAPIエンドポイントを実装したい。そうすることで、個別の通貨情報を表示できる。

#### Acceptance Criteria

1. WHEN GET /api/v1/crypto/{id} がリクエストされる THEN the System SHALL 指定されたIDの暗号通貨データを返す
2. WHEN 存在しないIDがリクエストされる THEN the System SHALL 404ステータスコードを返す
3. WHEN データが返される THEN the System SHALL すべての必須フィールドを含める

### Requirement 4

**User Story:** 開発者として、チャートデータを取得するAPIエンドポイントを実装したい。そうすることで、価格推移グラフを表示できる。

#### Acceptance Criteria

1. WHEN GET /api/v1/crypto/{id}/chart がリクエストされる THEN the System SHALL 指定された暗号通貨のチャートデータを返す
2. WHEN チャートデータが取得される THEN the System SHALL price_historiesテーブルから過去データを取得する
3. WHEN 1日に複数の価格データが存在する THEN the System SHALL その日の平均価格を計算する
4. WHEN periodパラメータが指定される THEN the System SHALL 指定された期間のデータを返す
5. WHEN periodパラメータが省略される THEN the System SHALL デフォルトで7日間のデータを返す
6. WHEN チャートデータが返される THEN the System SHALL 曜日ラベルと価格の配列を含める

### Requirement 5

**User Story:** 開発者として、MVCアーキテクチャで実装したい。そうすることで、コードの保守性と可読性が向上する。

#### Acceptance Criteria

1. WHEN プロジェクト構造が作成される THEN the System SHALL Model、Handler、Repository、Service層に分離する
2. WHEN Handlerが実装される THEN the System SHALL HTTPリクエスト/レスポンスの処理のみを行う
3. WHEN Serviceが実装される THEN the System SHALL ビジネスロジックを実装する
4. WHEN Repositoryが実装される THEN the System SHALL データアクセスロジックを実装する
5. WHEN Modelが実装される THEN the System SHALL データ構造を定義する

### Requirement 6

**User Story:** 開発者として、Echoフレームワークを使用したい。そうすることで、高速で軽量なAPIサーバーを構築できる。

#### Acceptance Criteria

1. WHEN サーバーが起動される THEN the System SHALL Echoフレームワークを使用する
2. WHEN ルーティングが設定される THEN the System SHALL Echoのルーターを使用する
3. WHEN ミドルウェアが設定される THEN the System SHALL Logger、CORS、Recoverミドルウェアを使用する
4. WHEN サーバーが起動される THEN the System SHALL ポート8080でリッスンする

### Requirement 7

**User Story:** 開発者として、CORSを設定したい。そうすることで、フロントエンド（localhost:3000）からAPIにアクセスできる。

#### Acceptance Criteria

1. WHEN CORSが設定される THEN the System SHALL localhost:3000からのリクエストを許可する
2. WHEN プリフライトリクエストが送信される THEN the System SHALL 適切なCORSヘッダーを返す
3. WHEN 許可されたメソッドが定義される THEN the System SHALL GET、POST、PUT、DELETE、OPTIONSを含める

### Requirement 8

**User Story:** 開発者として、既存のデータベースに接続したい。そうすることで、実際の取引データを使用できる。

#### Acceptance Criteria

1. WHEN データベース接続が実装される THEN the System SHALL 既存のRDS（MySQL）に接続する
2. WHEN 接続情報が設定される THEN the System SHALL 環境変数から接続情報を読み込む
3. WHEN price_historiesテーブルが照会される THEN the System SHALL product_codeとdatetimeでフィルタリングする
4. WHEN 日毎の平均価格が計算される THEN the System SHALL DATE(datetime)でグループ化してAVG(price)を計算する

### Requirement 11

**User Story:** 開発者として、bitFlyer Lightning APIから現在価格を取得したい。そうすることで、リアルタイムの市場価格を表示できる。

#### Acceptance Criteria

1. WHEN 現在価格が必要である THEN the System SHALL bitFlyer Lightning API（GET /v1/ticker）を呼び出す
2. WHEN APIリクエストが送信される THEN the System SHALL product_codeパラメータを指定する（BTC_JPY, ETH_JPY）
3. WHEN APIレスポンスが返される THEN the System SHALL ltp（Last Traded Price）フィールドを現在価格として使用する
4. WHEN 24時間変動率が必要である THEN the System SHALL price_ratio_24hまたは計算により取得する
5. WHEN API呼び出しが失敗する THEN the System SHALL エラーをログに記録し、適切なエラーレスポンスを返す

### Requirement 9

**User Story:** 開発者として、ユニットテストを実装したい。そうすることで、コードの品質と信頼性が向上する。

#### Acceptance Criteria

1. WHEN テストが実装される THEN the System SHALL 各層（Handler、Service、Repository）のテストを含める
2. WHEN Repositoryテストが実装される THEN the System SHALL モックデータを使用する
3. WHEN Handlerテストが実装される THEN the System SHALL HTTPリクエスト/レスポンスをテストする
4. WHEN テストが実行される THEN the System SHALL go test コマンドで実行可能にする

### Requirement 12

**User Story:** 開発者として、E2Eテストを実装したい。そうすることで、実際のDB・APIとの統合が正しく動作することを確認できる。

#### Acceptance Criteria

1. WHEN E2Eテストが実装される THEN the System SHALL 実際のRDS MySQLに接続する
2. WHEN E2Eテストが実装される THEN the System SHALL 実際のbitFlyer Lightning APIに接続する
3. WHEN E2Eテストが実行される THEN the System SHALL 正常系のレスポンスが返ることを確認する
4. WHEN APIエンドポイントがテストされる THEN the System SHALL GET /api/v1/crypto/market のテストを含める
5. WHEN APIエンドポイントがテストされる THEN the System SHALL GET /api/v1/crypto/{id} のテストを含める
6. WHEN APIエンドポイントがテストされる THEN the System SHALL GET /api/v1/crypto/{id}/chart のテストを含める
7. WHEN レスポンスが検証される THEN the System SHALL HTTPステータスコード200を確認する
8. WHEN レスポンスが検証される THEN the System SHALL 必須フィールドの存在を確認する
9. WHEN テストが実行される THEN the System SHALL go test -tags=e2e コマンドで実行可能にする

### Requirement 10

**User Story:** 開発者として、エラーハンドリングを実装したい。そうすることで、適切なエラーメッセージをクライアントに返せる。

#### Acceptance Criteria

1. WHEN エラーが発生する THEN the System SHALL 適切なHTTPステータスコードを返す
2. WHEN エラーレスポンスが返される THEN the System SHALL エラータイプとメッセージを含める
3. WHEN 404エラーが発生する THEN the System SHALL "NOT_FOUND" エラータイプを返す
4. WHEN 500エラーが発生する THEN the System SHALL "INTERNAL_SERVER_ERROR" エラータイプを返す
5. WHEN エラーがログに記録される THEN the System SHALL Echoのロガーを使用する
