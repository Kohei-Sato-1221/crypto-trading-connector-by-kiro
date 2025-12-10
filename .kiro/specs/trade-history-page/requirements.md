# 要件定義書

## 概要

この機能は、暗号通貨取引アプリケーション用の包括的な取引履歴ページを実装します。ユーザーの取引統計、利益計算、詳細な取引ログを表示します。ユーザーは異なる期間での取引パフォーマンスを確認し、暗号通貨資産（BTCとETH）で取引をフィルタリングできます。

## 用語集

- **Trade_History_System**: 取引履歴と統計を表示するWebアプリケーションコンポーネント
- **Transaction_Log**: 利益計算を含む完了した買い注文と売り注文の時系列リスト
- **Total_Profit**: 完了したすべての売り注文から対応する買い注文を差し引いて計算された累積利益
- **Execution_Count**: 完了した売り注文（約定）の総数
- **Buy_Order**: 暗号通貨の購入を表すbuy_ordersテーブルのレコード
- **Sell_Order**: parentidを介して買い注文にリンクされた暗号通貨の売却を表すsell_ordersテーブルのレコード
- **Asset_Filter**: ユーザーが暗号通貨タイプで取引をフィルタリングできるUIコンポーネント
- **Time_Filter**: ユーザーが期間（全期間、過去7日間）でデータをフィルタリングできるUIコンポーネント
- **Profit_Calculation**: 数学的操作：(売却価格 * 売却数量) - (購入価格 * 購入数量)

## 要件

### 要件 1

**ユーザーストーリー:** トレーダーとして、全体的な取引パフォーマンス統計を確認したいので、収益性と取引活動を理解できます。

#### 受入基準

1. WHEN ユーザーが取引履歴ページを訪問する時、THE Trade_History_System SHALL 完了したすべての売り注文から計算された総利益を表示する
2. WHEN 総利益を計算する時、THE Trade_History_System SHALL 各取引ペアについて対応する買い注文の価値を売り注文の価値から差し引く
3. WHEN 利益を表示する時、THE Trade_History_System SHALL 適切な色分けでパーセンテージ変化インジケーターを表示する
4. WHEN ユーザーが統計を確認する時、THE Trade_History_System SHALL 約定総数（完了した売り注文）を表示する
5. WHERE 利益データが利用可能な場合、THE Trade_History_System SHALL 適切な千の位区切り文字で日本円の金額値をフォーマットする

### 要件 2

**ユーザーストーリー:** トレーダーとして、取引の詳細な取引ログを確認したいので、個別の取引詳細とパフォーマンスをレビューできます。

#### 受入基準

1. WHEN ユーザーが取引ログを確認する時、THE Trade_History_System SHALL 各売り注文を対応する買い注文の詳細と共に表示する
2. WHEN 取引詳細を表示する時、THE Trade_History_System SHALL 暗号通貨名、タイムスタンプ、利益額、注文タイプ、注文IDを表示する
3. WHEN 注文詳細を表示する時、THE Trade_History_System SHALL 購入価格と売却価格をそれぞれの数量と共に表示する
4. WHEN 個別取引利益を計算する時、THE Trade_History_System SHALL 次の式を使用する：(売却価格 * 売却数量) - (購入価格 * 購入数量)
5. WHEN タイムスタンプを表示する時、THE Trade_History_System SHALL ユーザーフレンドリーな形式でフォーマットする（例：「今日、14:30」、「昨日、09:15」）

### 要件 3

**ユーザーストーリー:** トレーダーとして、暗号通貨タイプで取引履歴をフィルタリングしたいので、特定の資産のパフォーマンスに焦点を当てることができます。

#### 受入基準

1. WHEN ユーザーが取引履歴ページを確認する時、THE Trade_History_System SHALL 「全資産」、「BTC」、「ETH」のフィルターボタンを提供する
2. WHEN ユーザーが資産フィルターを選択する時、THE Trade_History_System SHALL その特定の暗号通貨の取引のみを表示する
3. WHEN 資産でフィルタリングする時、THE Trade_History_System SHALL 統計と取引ログの両方を適切に更新する
4. WHEN 「全資産」が選択される時、THE Trade_History_System SHALL サポートされているすべての暗号通貨の取引を表示する
5. WHERE 選択されたフィルターに対して取引が存在しない場合、THE Trade_History_System SHALL 適切な空の状態を表示する

### 要件 4

**ユーザーストーリー:** トレーダーとして、期間で取引データをフィルタリングしたいので、特定の時間枠でのパフォーマンスを分析できます。

#### 受入基準

1. WHEN ユーザーが取引履歴ページを確認する時、THE Trade_History_System SHALL 「全期間」と「過去7日間」の時間フィルターオプションを提供する
2. WHEN ユーザーが時間フィルターを選択する時、THE Trade_History_System SHALL その時間範囲内のデータのみを表示する
3. WHEN 時間でフィルタリングする時、THE Trade_History_System SHALL 選択された期間の利益計算と約定数を更新する
4. WHEN 「全期間」が選択される時、THE Trade_History_System SHALL 利用可能なすべての取引データを表示する
5. WHEN 「過去7日間」が選択される時、THE Trade_History_System SHALL 過去7日間の取引のみを表示する

### 要件 5

**ユーザーストーリー:** トレーダーとして、レスポンシブデザインでページが素早く読み込まれることを望むので、どのデバイスでも効率的に取引履歴にアクセスできます。

#### 受入基準

1. WHEN ユーザーが取引履歴ページを読み込む時、THE Trade_History_System SHALL 通常のネットワーク条件下で2秒以内にデータを取得して表示する
2. WHEN ページがモバイルデバイスで表示される時、THE Trade_History_System SHALL レスポンシブレイアウトで読みやすさと使いやすさを維持する
3. WHEN データが読み込み中の時、THE Trade_History_System SHALL 適切な読み込みインジケーターを表示する
4. WHERE ネットワークエラーが発生する場合、THE Trade_History_System SHALL 再試行オプション付きのユーザーフレンドリーなエラーメッセージを表示する
5. WHEN ページがレンダリングされる時、THE Trade_History_System SHALL 既存のアプリケーションデザインシステムと一貫したスタイリングを維持する

### 要件 6

**ユーザーストーリー:** トレーダーとして、古い取引を閲覧したいので、完全な取引履歴をレビューできます。

#### 受入基準

1. WHEN 取引ログが初期の取引セットを表示する時、THE Trade_History_System SHALL 「古い取引を表示」ボタンを提供する
2. WHEN ユーザーが「古い取引を表示」ボタンをクリックする時、THE Trade_History_System SHALL 追加の履歴取引を読み込んで表示する
3. WHEN 古い取引を読み込む時、THE Trade_History_System SHALL 現在のフィルター設定を維持する
4. WHERE これ以上の取引が利用できない場合、THE Trade_History_System SHALL 「古い取引を表示」ボタンを非表示または無効にする
5. WHEN ページネーションが発生する時、THE Trade_History_System SHALL ユーザーの現在のスクロール位置とコンテキストを保持する