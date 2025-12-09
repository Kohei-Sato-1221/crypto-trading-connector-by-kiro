# Requirements Document

## Introduction

OpenAPI仕様書（openapi.yaml）から、バックエンド（Go）とフロントエンド（TypeScript）のコードを自動生成する仕組みを実装します。これにより、API仕様とコードの一貫性を保ち、開発効率を向上させます。

## Glossary

- **OpenAPI**: RESTful APIを記述するための標準仕様
- **oapi-codegen**: OpenAPI仕様からGoコードを生成するツール
- **openapi-typescript**: OpenAPI仕様からTypeScript型定義を生成するツール
- **Code Generation**: 仕様書からソースコードを自動生成すること

## Requirements

### Requirement 1

**User Story:** 開発者として、OpenAPI仕様書からGoのモデルとインターフェースを生成したい。そうすることで、型安全性が保証され、手動でのコード記述ミスを防げる。

#### Acceptance Criteria

1. WHEN oapi-codegenが実行される THEN the System SHALL openapi.yamlからGoの型定義を生成する
2. WHEN 生成されたコードが配置される THEN the System SHALL internal/generated/models.goに出力する
3. WHEN 生成コマンドが実行される THEN the System SHALL Makefileから実行可能にする
4. WHEN 生成されたコードがビルドされる THEN the System SHALL エラーなくコンパイルできる

### Requirement 2

**User Story:** 開発者として、OpenAPI仕様書からTypeScriptの型定義を生成したい。そうすることで、フロントエンドでAPI型安全性が保証される。

#### Acceptance Criteria

1. WHEN openapi-typescriptが実行される THEN the System SHALL openapi.yamlからTypeScript型定義を生成する
2. WHEN 生成されたコードが配置される THEN the System SHALL types/api.tsに出力する
3. WHEN 生成コマンドが実行される THEN the System SHALL Makefileから実行可能にする
4. WHEN 生成された型が使用される THEN the System SHALL 既存のcomposableで利用できる

### Requirement 3

**User Story:** 開発者として、OpenAPI仕様書を更新したら自動的にコードを再生成したい。そうすることで、常に最新の仕様に基づいたコードを維持できる。

#### Acceptance Criteria

1. WHEN make generateが実行される THEN the System SHALL バックエンドとフロントエンドの両方のコードを生成する
2. WHEN 生成コマンドが失敗する THEN the System SHALL エラーメッセージを表示する
3. WHEN 生成されたファイルが存在する THEN the System SHALL 上書きする

### Requirement 4

**User Story:** 開発者として、生成されたコードと手書きコードを分離したい。そうすることで、生成コードを再生成しても手書きコードが失われない。

#### Acceptance Criteria

1. WHEN 生成コードが配置される THEN the System SHALL generatedディレクトリに配置する
2. WHEN 手書きコードが配置される THEN the System SHALL generatedディレクトリ外に配置する
3. WHEN .gitignoreが更新される THEN the System SHALL 生成ファイルを除外しない（バージョン管理する）

### Requirement 5

**User Story:** 開発者として、生成されたコードの使い方を理解したい。そうすることで、既存コードとの統合がスムーズに行える。

#### Acceptance Criteria

1. WHEN READMEが更新される THEN the System SHALL コード生成の手順を記載する
2. WHEN READMEが更新される THEN the System SHALL 生成コマンドの使い方を記載する
3. WHEN READMEが更新される THEN the System SHALL 既存コードとの統合方法を記載する
