# Design Document

## Overview

OpenAPI仕様書（openapi.yaml）から、バックエンド（Go）とフロントエンド（TypeScript）のコードを自動生成します。oapi-codegenとopenapi-typescriptを使用し、Makefileから簡単に実行できるようにします。

## Architecture

### Technology Stack

- **Backend Code Generator**: oapi-codegen v2
- **Frontend Code Generator**: openapi-typescript
- **Build Tool**: Make
- **OpenAPI Version**: 3.0.3

### Directory Structure

```
crypto-trading-connector-be/
├── internal/
│   └── generated/
│       └── models.go          # 生成されたGoモデル
├── openapi.yaml               # OpenAPI仕様書（ルート）
└── Makefile

crypto-trading-connector-fe/
├── types/
│   └── api.ts                 # 生成されたTypeScript型定義
└── Makefile
```

## Components

### 1. Backend Code Generation (oapi-codegen)

**生成内容:**
- Go構造体（models）
- JSONタグ
- バリデーション用のインターフェース

**設定ファイル:** 不要（コマンドラインオプションで指定）

**生成コマンド:**
```bash
oapi-codegen -package generated -generate types ../openapi.yaml > internal/generated/models.go
```

**生成されるコード例:**
```go
package generated

type CryptoData struct {
    ID            string           `json:"id"`
    Name          string           `json:"name"`
    Symbol        string           `json:"symbol"`
    CurrentPrice  float64          `json:"currentPrice"`
    ChangePercent float64          `json:"changePercent"`
    ChartData     []ChartDataPoint `json:"chartData"`
}

type ChartDataPoint struct {
    Day   string  `json:"day"`
    Price float64 `json:"price"`
}
```

### 2. Frontend Code Generation (openapi-typescript)

**生成内容:**
- TypeScript型定義
- パス定義
- レスポンス型

**生成コマンド:**
```bash
npx openapi-typescript ../openapi.yaml -o types/api.ts
```

**生成されるコード例:**
```typescript
export interface paths {
  "/api/v1/crypto/market": {
    get: {
      responses: {
        200: {
          content: {
            "application/json": components["schemas"]["MarketResponse"];
          };
        };
      };
    };
  };
}

export interface components {
  schemas: {
    CryptoData: {
      id: string;
      name: string;
      symbol: string;
      currentPrice: number;
      changePercent: number;
      chartData: components["schemas"]["ChartDataPoint"][];
    };
    ChartDataPoint: {
      day: string;
      price: number;
    };
  };
}
```

## Makefile Integration

### Backend Makefile

```makefile
## generate: Generate Go code from OpenAPI spec
generate:
	@echo "Generating Go code from OpenAPI spec..."
	@mkdir -p internal/generated
	oapi-codegen -package generated -generate types ../openapi.yaml > internal/generated/models.go
	@echo "Generated: internal/generated/models.go"
	go fmt ./internal/generated/...
```

### Frontend Makefile

```makefile
## generate: Generate TypeScript types from OpenAPI spec
generate:
	@echo "Generating TypeScript types from OpenAPI spec..."
	@mkdir -p types
	npx openapi-typescript ../openapi.yaml -o types/api.ts
	@echo "Generated: types/api.ts"
```

## Integration with Existing Code

### Backend Integration

既存の`internal/model/crypto.go`は残し、生成されたコードは参照用として使用します。将来的に、手書きコードを生成コードに置き換えることができます。

**段階的な移行:**
1. 生成コードを確認
2. 既存コードと比較
3. 必要に応じて既存コードを生成コードに置き換え

### Frontend Integration

生成された型定義を既存のcomposableで使用します。

**使用例:**
```typescript
import type { components } from '~/types/api'

type CryptoData = components['schemas']['CryptoData']
type MarketResponse = components['schemas']['MarketResponse']

export const useCryptoData = () => {
  const cryptoData = ref<CryptoData[]>([])
  // ...
}
```

## Error Handling

### 生成エラー

- OpenAPI仕様書の構文エラー
- ツールのインストール不足
- ファイル書き込み権限エラー

**対処:**
- エラーメッセージを表示
- 生成前にopenapi.yamlをバリデーション
- 必要なツールのインストール確認

## Testing Strategy

### 生成コードのテスト

- 生成コードがコンパイルできることを確認
- 既存のテストが引き続き動作することを確認

### バリデーション

- OpenAPI仕様書のバリデーション
- 生成されたコードの構文チェック

## Future Enhancements

1. **CI/CD統合**: PRごとに自動生成して差分をチェック
2. **バリデーション**: 生成前にOpenAPI仕様書をバリデーション
3. **クライアント生成**: APIクライアントコードの自動生成
4. **ドキュメント生成**: Swagger UIの自動生成
