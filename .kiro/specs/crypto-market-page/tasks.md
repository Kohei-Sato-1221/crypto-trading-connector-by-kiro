# Implementation Plan

- [x] 1. プロジェクトのセットアップとベースページ作成
  - Nuxt 4プロジェクトの初期化
  - Tailwind CSSの設定
  - TypeScript設定の確認
  - Google Fonts（Manrope）の追加
  - _Requirements: 9.1, 9.2, 9.3_

- [x] 2. 基本ページとルーティングの作成
  - pages/market.vue（空ページ）を作成
  - pages/trade.vue（空ページ）を作成
  - pages/history.vue（空ページ）を作成
  - pages/index.vueを作成し、/marketへリダイレクト
  - _Requirements: 5.4, 5.5_

- [x] 3. ナビゲーションバーコンポーネントの実装
  - components/NavigationBar.vueを作成
  - 4つのナビゲーション項目（Market、Trade、History、Portfolio）を実装
  - アクティブ状態のスタイリング（青色ハイライト）
  - 固定配置（bottom: 0）の実装
  - アイコンとラベルの表示
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 3.1 ナビゲーションバーのプロパティテストを作成
  - **Property 5: Navigation active state**
  - **Validates: Requirements 5.3**

- [x] 3.2 ナビゲーションバーのユニットテストを作成
  - 4つのリンクが表示されることを確認
  - クリック時のルート遷移を確認
  - _Requirements: 5.2, 5.4_

- [x] 4. Checkpoint - ページ遷移の確認
  - すべてのページ（market、trade、history）が作成されていることを確認
  - ナビゲーションバーからページ遷移が正しく動作することを確認
  - 各ページでナビゲーションバーのアクティブ状態が正しく表示されることを確認

- [ ] 5. 型定義の作成
  - types/crypto.tsを作成
  - CryptoData型を定義
  - ChartDataPoint型を定義
  - TimeFilter型を定義
  - NavigationItem型を定義
  - _Requirements: 8.3, 8.4_

- [ ] 5.1 型定義のプロパティテストを作成
  - **Property 6: Mock data structure validity**
  - **Validates: Requirements 8.3, 8.4**

- [ ] 6. モックデータの作成
  - utils/mockData.tsを作成
  - Bitcoin（BTC/JPY）のモックデータを作成
  - Ethereum（ETH/JPY）のモックデータを作成
  - 7日間のチャートデータを生成
  - 価格更新関数（updateMockPrice）を実装
  - _Requirements: 1.2, 7.3_


- [ ] 6.1 モックデータのプロパティテストを作成
  - **Property 7: Price update maintains data integrity**
  - **Validates: Requirements 7.2**

- [ ] 7. Composablesの実装
  - composables/useCryptoData.tsを作成
  - モックデータ取得関数を実装
  - 将来のAPI統合を見据えたインターフェース設計
  - _Requirements: 8.1, 8.2_

- [ ] 7.1 自動更新Composableの実装
  - composables/useAutoRefresh.tsを作成
  - setIntervalを使用した5秒ごとの更新ロジック
  - タイマーのクリーンアップ処理（onUnmounted）
  - start/stop関数の実装
  - _Requirements: 7.1, 7.4_

- [ ] 7.2 自動更新のプロパティテストを作成
  - **Property 8: Auto-refresh timer cleanup**
  - **Validates: Requirements 7.4**

- [ ] 7.3 時間フィルターComposableの実装
  - composables/useTimeFilter.tsを作成
  - 選択状態の管理（ref）
  - デフォルト値を'24h'に設定
  - _Requirements: 4.2, 4.3_

- [ ]* 7.4 時間フィルターのプロパティテストを作成
  - **Property 4: Time filter exclusivity**
  - **Validates: Requirements 4.2, 4.3**

- [ ] 8. ヘッダーコンポーネントの実装
  - components/MarketHeader.vueを作成
  - アプリケーションタイトルの表示
  - 市場ステータスインジケーター（緑色の点 + テキスト）
  - 時間フィルターボタン（24h、7d、30d、1y、All）
  - 固定ヘッダー（sticky positioning）
  - 背景ブラー効果
  - _Requirements: 1.5, 4.1, 4.2, 4.3_

- [ ] 8.1 ヘッダーのユニットテストを作成
  - 市場ステータスの表示を確認
  - 5つの時間フィルターボタンが表示されることを確認
  - デフォルトで24hが選択されていることを確認
  - _Requirements: 1.5, 4.1, 4.2_

- [ ] 9. チャートコンポーネントの実装
  - components/PriceChart.vueを作成
  - Chart.jsとvue-chartjsのインストール
  - エリアチャートの実装
  - 動的な色変更（緑/赤のグラデーション）
  - X軸：曜日ラベル（Mon〜Sun）
  - Y軸：価格範囲
  - レスポンシブサイズ対応
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [ ] 9.1 チャートのプロパティテストを作成
  - **Property 3: Chart color consistency**
  - **Validates: Requirements 2.4, 2.5**

- [ ] 9.2 チャートのユニットテストを作成
  - 曜日ラベルが表示されることを確認
  - Y軸に価格が表示されることを確認
  - _Requirements: 2.2, 2.3_

- [ ] 10. 暗号通貨カードコンポーネントの実装
  - components/CryptoCard.vueを作成
  - 通貨アイコン（円形背景）の表示
  - 通貨名とペア（BTC/JPY）の表示
  - 現在価格の表示（大きなフォント）
  - 変動率バッジ（色分け：緑=上昇、赤=下降、矢印アイコン付き）
  - PriceChartコンポーネントの統合
  - Tradeボタンの実装
  - トランジション効果（価格更新時）
  - _Requirements: 1.2, 1.3, 1.4, 2.1, 3.1, 7.5_

- [ ]* 10.1 暗号通貨カードのプロパティテストを作成
  - **Property 1: Crypto card rendering completeness**
  - **Validates: Requirements 1.2**

- [ ]* 10.2 暗号通貨カードのプロパティテストを作成
  - **Property 2: Price change color consistency**
  - **Validates: Requirements 1.3, 1.4**

- [ ]* 10.3 暗号通貨カードのユニットテストを作成
  - Tradeボタンが表示されることを確認
  - Tradeボタンクリック時の遷移を確認
  - _Requirements: 3.1, 3.2_

- [ ] 11. マーケットページの完全実装
  - pages/market.vueを完全に実装
  - MarketHeaderコンポーネントの統合
  - CryptoCardコンポーネントのリスト表示
  - useCryptoDataでモックデータを取得
  - useAutoRefreshで5秒ごとの自動更新を実装
  - スクロール可能なコンテンツエリア
  - NavigationBarコンポーネントの統合
  - _Requirements: 1.1, 6.2, 6.3, 7.1, 7.2_

- [ ]* 11.1 マーケットページのユニットテストを作成
  - 暗号通貨カードのリストが表示されることを確認
  - ヘッダーとナビゲーションバーが固定表示されることを確認
  - _Requirements: 1.1, 6.2, 6.3_

- [ ] 12. モバイル最適化とレスポンシブ対応
  - モバイルビューポート（375px〜428px）の調整
  - タッチターゲットサイズの確認（最小44px）
  - スクロール動作の最適化
  - _Requirements: 6.1, 6.4_

- [ ]* 12.1 タッチターゲットサイズのプロパティテストを作成
  - **Property 9: Touch target size compliance**
  - **Validates: Requirements 6.4**

- [ ] 13. スタイリングの仕上げ
  - カラーパレットの適用（ダークテーマ）
  - タイポグラフィの調整（Manropeフォント）
  - スペーシングとボーダー半径の統一
  - トランジション効果の追加
  - _Requirements: 7.5_

- [ ] 14. Final Checkpoint - すべてのテストが通ることを確認
  - すべてのテストが通ることを確認
  - ページ遷移が正しく動作することを確認
  - 自動更新が正しく動作することを確認
  - 質問があればユーザーに確認
