# Crypto Trading Connector

暗号通貨(BTC, ETH)取引Webアプリケーション（フロントエンド + バックエンドAPI） 
1. 各通貨の現在価格を確認
2. 各通貨の買い注文を行う
3. 自動売買における利益を計算して表示

<img width="317" height="691" alt="crypto-sugar-ss-1" src="https://github.com/user-attachments/assets/50e6bc64-e156-4adb-8dce-ef2875006d23" />
<img width="317" height="691" alt="crypto-sugar-ss-2" src="https://github.com/user-attachments/assets/52fc0d45-a071-43b2-8542-6cdb85a69d07" />
<img width="314" height="690" alt="crypto-sugar-ss-3" src="https://github.com/user-attachments/assets/77a550f1-34d4-4434-859a-0eac61ab2a11" />

## プロジェクト構成

- **crypto-trading-connector-be**: バックエンドAPI（Go + Echo）
- **crypto-trading-connector-fe**: フロントエンド（Nuxt 4 + TypeScript）
- **openapi.yaml**: API仕様書（OpenAPI 3.0.3）
- Kiroを用いて仕様駆動開発を行なっている。

## OpenAPIコード生成

このプロジェクトでは、OpenAPI仕様書からコードを自動生成します。

### 全体生成（推奨）

バックエンドとフロントエンドの両方のコードを一度に生成：

```bash
make gen
```

### 個別生成

バックエンドのみ：
```bash
make gen-be
```

フロントエンドのみ：
```bash
make gen-fe
```

### 生成されるファイル

- **バックエンド**: `crypto-trading-connector-be/internal/generated/models.go`
- **フロントエンド**: `crypto-trading-connector-fe/types/api.ts`

## クイックスタート

### 1. バックエンドを起動

```bash
cd crypto-trading-connector-be
make run
```

### 2. フロントエンドを起動

```bash
cd crypto-trading-connector-fe
make run
```

### 3. ブラウザでアクセス

- **ローカル**: http://localhost:3000/market
- **ネットワーク**: http://YOUR_IP:3000/market

## 開発ワークフロー

### OpenAPI仕様書を更新した場合

1. `openapi.yaml`を編集
2. コードを再生成: `make gen`
3. 必要に応じて実装コードを更新
4. テストを実行

### 新しいエンドポイントを追加する場合

1. `openapi.yaml`に新しいエンドポイントを追加
2. `make gen`でコードを生成
3. バックエンドにハンドラー、サービス、リポジトリを実装
4. フロントエンドでcomposableを更新

## ドキュメント

- [バックエンドREADME](crypto-trading-connector-be/README.md)
- [フロントエンドREADME](crypto-trading-connector-fe/README.md)
- [OpenAPI仕様書](openapi.yaml)

## 備考
- 当該アプリは仕様駆動開発（Kiro）を用いて実装されています。
- 当該アプリはWebのUIを提供します。DBや自動売買のロジックは[こちら](https://github.com/Kohei-Sato-1221)のアプリで実装されています（当該アプリは、こちらのリポジトリで操作されているDBを参照するアプリです。）

## ライセンス

MIT License
