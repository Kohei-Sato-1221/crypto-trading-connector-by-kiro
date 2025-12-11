# Raspberry Pi デプロイメントガイド

このガイドでは、Crypto Trading Connectorをラズパイ（IP: 10.221.1.162）にデプロイし、systemdでデーモン化する手順を説明します。

## 前提条件

- Raspberry Pi OS (Debian系)
- Go 1.19以上
- Node.js 18以上
- npm
- systemd

## 1. 依存関係のインストール

```bash
# Go のインストール (必要に応じて)
wget https://go.dev/dl/go1.21.5.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Node.js のインストール (必要に応じて)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# その他の依存関係
sudo apt-get update
sudo apt-get install -y git make curl jq lsof
```

## 2. プロジェクトのクローン

```bash
cd /home/pi
git clone <your-repo-url> crypto-trading-connector-by-kiro
cd crypto-trading-connector-by-kiro
```

## 3. 環境設定

プロジェクトは既にラズパイ用に設定済みです：

- **Backend**: ポート8090で動作
- **Frontend**: ポート3000で動作
- **IP**: 10.221.1.162

### 設定ファイルの確認

```bash
# バックエンド設定
cat crypto-trading-connector-be/.env

# フロントエンド設定
cat crypto-trading-connector-fe/.env
```

## 4. 依存関係のインストール

```bash
# フロントエンドの依存関係
cd crypto-trading-connector-fe
npm install
cd ..

# バックエンドの依存関係（Go modules）
cd crypto-trading-connector-be
go mod tidy
cd ..
```

## 5. systemdサービスのインストール

```bash
# systemdサービスファイルをインストール
sudo make install-services
```

これにより以下が実行されます：
- `systemd/crypto-trading-be.service` → `/etc/systemd/system/`
- `systemd/crypto-trading-fe.service` → `/etc/systemd/system/`
- サービスの有効化

## 6. サービスの開始

```bash
# サービス開始
make start

# サービス状態確認
make status
```

## 7. アクセス確認

### ローカルネットワークからのアクセス

- **Backend API**: http://10.221.1.162:8090
- **Frontend Web**: http://10.221.1.162:3000

### API テスト

```bash
# 市場データ取得
curl http://10.221.1.162:8090/api/v1/crypto/market

# Bitcoin データ取得
curl http://10.221.1.162:8090/api/v1/crypto/bitcoin
```

## 8. サービス管理コマンド

```bash
# サービス開始
make start

# サービス停止
make stop

# サービス状態確認
make status

# 開発モードで起動（デバッグ用）
make run
```

## 9. ログの確認

```bash
# systemd ログ
sudo journalctl -u crypto-trading-be -f
sudo journalctl -u crypto-trading-fe -f

# 開発モード時のログ
tail -f logs/backend.log
tail -f logs/frontend.log
```

## 10. トラブルシューティング

### サービスが起動しない場合

```bash
# サービス状態の詳細確認
sudo systemctl status crypto-trading-be
sudo systemctl status crypto-trading-fe

# ログの確認
sudo journalctl -u crypto-trading-be --no-pager -l
sudo journalctl -u crypto-trading-fe --no-pager -l
```

### ポートが使用中の場合

```bash
# ポート使用状況確認
sudo lsof -i :8090
sudo lsof -i :3000

# プロセス強制終了
sudo pkill -f "crypto-trading"
```

### 設定の再読み込み

```bash
# systemd設定の再読み込み
sudo systemctl daemon-reload

# サービスの再起動
sudo systemctl restart crypto-trading-be
sudo systemctl restart crypto-trading-fe
```

## 11. 自動起動の設定

サービスは既に自動起動に設定されていますが、手動で設定する場合：

```bash
sudo systemctl enable crypto-trading-be
sudo systemctl enable crypto-trading-fe
```

## 12. セキュリティ考慮事項

- サービスは `pi` ユーザーで実行されます
- 必要最小限の権限で動作するよう設定済み
- ファイアウォール設定が必要な場合は適切に設定してください

## 13. 更新手順

```bash
# コードの更新
git pull

# サービス停止
make stop

# 依存関係の更新（必要に応じて）
cd crypto-trading-connector-fe && npm install && cd ..
cd crypto-trading-connector-be && go mod tidy && cd ..

# サービス再開
make start
```

## 注意事項

- 本番環境では適切なセキュリティ設定を行ってください
- データベース接続情報などの機密情報は適切に管理してください
- 定期的なバックアップを推奨します