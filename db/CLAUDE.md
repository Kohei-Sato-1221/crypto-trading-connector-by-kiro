# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## DB ディレクトリ

本プロジェクトではDBマイグレーションは実装しない。別アプリ(crypto-trading-golang)で管理されているDBをモジュラーモノリス的に参照する。

このディレクトリには参考用のAtlas HCLスキーマファイル(`schema.hcl`)のみを配置。

### 主要テーブル（MySQL: crypto_trading_db）
- `buy_orders` — 買い注文（order_id, product_code, price, size, status, strategy）
- `sell_orders` — 売り注文（parentidでbuy_ordersと紐付け）
- `price_histories` — 価格履歴（product_code, price, datetime）

スキーマ変更は別リポジトリで行うこと。
