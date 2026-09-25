#!/bin/bash

# Crypto Trading Backend Startup Script for systemd
#
# 環境変数について:
#   .env は systemd の EnvironmentFile と、アプリ側の godotenv.Load() が
#   それぞれ読み込むため、このスクリプトでは読まない。
#   .env の値には & や () が含まれており `source .env` は構文エラーになる。

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# systemd 配下では PATH が最小限で go が見つからないため、
# 一般的なインストール先を補う。
if ! command -v go >/dev/null 2>&1; then
    for candidate in /usr/local/go/bin /usr/lib/go/bin /opt/go/bin "${HOME}/go/bin" /snap/bin; do
        if [ -x "${candidate}/go" ]; then
            export PATH="${candidate}:${PATH}"
            break
        fi
    done
fi

if ! command -v go >/dev/null 2>&1; then
    echo "ERROR: 'go' command not found." >&2
    echo "  PATH=${PATH}" >&2
    echo "  Goのインストール先を systemd unit の Environment=PATH に追加してください。" >&2
    exit 1
fi

# Always rebuild so the running binary matches the checked-out source.
# (bin/server の存在チェックでビルドをスキップすると、古いコードが動き続ける)
echo "Building server binary with $(go version)..."
go build -o bin/server ./cmd/server

# Start the server
echo "Starting crypto trading backend server..."
exec ./bin/server
