# Requirements Document

## Introduction

暗号通貨取引Webアプリケーションのマーケットページを実装します。このページでは、ユーザーが暗号通貨の現在価格、価格変動、チャートを確認できます。モバイルファーストのSPAとして、Nuxt 4、TypeScript、Tailwind CSSを使用して構築します。初期段階ではモックデータを使用し、将来的にREST APIからデータを取得できるよう設計します。

## Glossary

- **System**: 暗号通貨取引Webアプリケーション
- **Market Page**: 暗号通貨の市場情報を表示するページ（/market）
- **Crypto Card**: 個別の暗号通貨情報を表示するカードコンポーネント
- **Price Chart**: 価格推移を視覚化するチャートコンポーネント
- **Time Filter**: チャート表示期間を選択するフィルター（24h、7d、30d、1y、All）
- **Navigation Bar**: アプリケーション下部のナビゲーションバー
- **Mock Data**: バックエンド実装前に使用する静的なテストデータ
- **SPA**: Single Page Application（シングルページアプリケーション）

## Requirements

### Requirement 1

**User Story:** ユーザーとして、マーケットページで暗号通貨の一覧を確認したい。そうすることで、各通貨の現在価格と変動率を把握できる。

#### Acceptance Criteria

1. WHEN ユーザーが/marketページにアクセスする THEN the System SHALL 暗号通貨カードのリストを表示する
2. WHEN 暗号通貨カードが表示される THEN the System SHALL 通貨名、通貨ペア（例：BTC/JPY）、アイコン、現在価格、価格変動率、7日間チャートを含める
3. WHEN 価格変動率が正の値である THEN the System SHALL 緑色で表示し、上向き矢印アイコンを表示する
4. WHEN 価格変動率が負の値である THEN the System SHALL 赤色で表示し、下向き矢印アイコンを表示する
5. WHEN マーケットステータスが表示される THEN the System SHALL "Market is Open • Real-time Data"というテキストと緑色のインジケーターを表示する

### Requirement 2

**User Story:** ユーザーとして、各暗号通貨の価格チャートを確認したい。そうすることで、価格トレンドを視覚的に理解できる。

#### Acceptance Criteria

1. WHEN 暗号通貨カードが表示される THEN the System SHALL 7日間の価格推移チャートを表示する
2. WHEN チャートが表示される THEN the System SHALL X軸に曜日（Mon、Tue、Wed、Thu、Fri、Sat、Sun）を表示する
3. WHEN チャートが表示される THEN the System SHALL Y軸に価格範囲を表示する
4. WHEN 価格が上昇トレンドである THEN the System SHALL チャートを緑色のグラデーションで表示する
5. WHEN 価格が下降トレンドである THEN the System SHALL チャートを赤色のグラデーションで表示する

### Requirement 3

**User Story:** ユーザーとして、各暗号通貨から取引ページに移動したい。そうすることで、すぐに取引を開始できる。

#### Acceptance Criteria

1. WHEN 暗号通貨カードが表示される THEN the System SHALL "Trade"ボタンを表示する
2. WHEN ユーザーが"Trade"ボタンをクリックする THEN the System SHALL 取引ページ（/trade）に遷移する

### Requirement 4

**User Story:** ユーザーとして、ヘッダーで時間フィルターを選択したい。そうすることで、異なる期間のデータを確認できる。

#### Acceptance Criteria

1. WHEN ヘッダーが表示される THEN the System SHALL 時間フィルターボタン（24h、7d、30d、1y、All）を表示する
2. WHEN デフォルト状態である THEN the System SHALL "24h"ボタンを選択状態（青色）で表示する
3. WHEN ユーザーが時間フィルターボタンをクリックする THEN the System SHALL そのボタンを選択状態にし、他のボタンを非選択状態にする
4. WHEN 時間フィルターが変更される THEN the System SHALL 将来的にチャートデータを更新できるよう設計する

### Requirement 5

**User Story:** ユーザーとして、ナビゲーションバーで他のページに移動したい。そうすることで、アプリケーション内を自由に移動できる。

#### Acceptance Criteria

1. WHEN ページが表示される THEN the System SHALL 画面下部にナビゲーションバーを固定表示する
2. WHEN ナビゲーションバーが表示される THEN the System SHALL Market、Trade、History、Portfolioの4つのリンクを表示する
3. WHEN 現在のページがMarketである THEN the System SHALL Marketアイコンとラベルを青色で表示する
4. WHEN ユーザーがナビゲーションリンクをクリックする THEN the System SHALL 対応するページに遷移する
5. WHEN ユーザーがTrade、History、Portfolioリンクをクリックする THEN the System SHALL 空のページを表示する

### Requirement 6

**User Story:** ユーザーとして、モバイルデバイスで快適に閲覧したい。そうすることで、スマートフォンから簡単にアクセスできる。

#### Acceptance Criteria

1. WHEN ページが表示される THEN the System SHALL モバイルビューポート（375px〜428px）に最適化されたレイアウトを表示する
2. WHEN ページがスクロールされる THEN the System SHALL ヘッダーを固定表示し、コンテンツのみスクロールする
3. WHEN ページがスクロールされる THEN the System SHALL ナビゲーションバーを固定表示する
4. WHEN タッチ操作が行われる THEN the System SHALL 適切なタッチターゲットサイズ（最小44px）を提供する

### Requirement 7

**User Story:** ユーザーとして、マーケットページで最新の価格情報を自動的に確認したい。そうすることで、手動でリフレッシュすることなく市場の変動を追跡できる。

#### Acceptance Criteria

1. WHEN マーケットページが表示される THEN the System SHALL 5秒ごとに価格データを自動更新する
2. WHEN 価格データが更新される THEN the System SHALL 現在価格、価格変動率、チャートデータを最新の値に更新する
3. WHEN モックデータが使用される THEN the System SHALL ランダムな価格変動をシミュレートする
4. WHEN ページから離れる THEN the System SHALL 自動更新タイマーをクリーンアップする
5. WHEN 価格が更新される THEN the System SHALL スムーズなトランジション効果を適用する

### Requirement 8

**User Story:** 開発者として、将来的にREST APIからデータを取得できるようにしたい。そうすることで、バックエンドが実装された際にスムーズに統合できる。

#### Acceptance Criteria

1. WHEN データ取得ロジックが実装される THEN the System SHALL モックデータとAPI取得を切り替え可能な設計にする
2. WHEN コンポーネントが実装される THEN the System SHALL データソースに依存しないインターフェースを使用する
3. WHEN 型定義が作成される THEN the System SHALL TypeScriptの型でデータ構造を明確に定義する
4. WHEN APIレスポンスの型が定義される THEN the System SHALL 暗号通貨データの必須フィールド（id、name、symbol、price、changePercent、chartData）を含める

### Requirement 9

**User Story:** 開発者として、Nuxt 4のベストプラクティスに従いたい。そうすることで、保守性と拡張性の高いコードベースを構築できる。

#### Acceptance Criteria

1. WHEN プロジェクト構造が作成される THEN the System SHALL Nuxt 4の推奨ディレクトリ構造（pages、components、composables、types）を使用する
2. WHEN スタイリングが実装される THEN the System SHALL Tailwind CSSを使用する
3. WHEN コンポーネントが作成される THEN the System SHALL Vue 3のComposition APIを使用する
4. WHEN ルーティングが実装される THEN the System SHALL Nuxtのファイルベースルーティングを使用する
5. WHEN 状態管理が必要である THEN the System SHALL composablesパターンを使用する
