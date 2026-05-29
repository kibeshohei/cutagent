# CUTAGENT API

Go 製バックエンド。`net/http` + [Huma v2](https://github.com/danielgtaylor/huma) で実装し、OpenAPI 3.1 仕様書を自動生成する。

## 起動

```bash
# デフォルト :8080
go run ./cmd/server

# ポートを変える（Cloud Run はこれで受け取る）
PORT=9090 go run ./cmd/server
```

起動後、以下で仕様書を確認できる。

| URL | 内容 |
|-----|------|
| `http://localhost:8080/docs` | Swagger UI |
| `http://localhost:8080/openapi.json` | OpenAPI JSON（TS 型生成に使う） |

## エンドポイント一覧

### システム
| メソッド | パス | 概要 |
|----------|------|------|
| GET | `/api/health` | ヘルスチェック |

### 体重
| メソッド | パス | 概要 |
|----------|------|------|
| GET | `/api/weight` | 体重記録一覧 |
| POST | `/api/weight` | 体重記録追加 |

### 食事
| メソッド | パス | 概要 |
|----------|------|------|
| GET | `/api/meals?date=` | 指定日の食事記録（未指定なら全件） |
| POST | `/api/meals` | 食事記録追加 |
| DELETE | `/api/meals/{id}` | 食事記録削除 |
| GET | `/api/meal-master` | プリセット食品一覧 |

### 運動
| メソッド | パス | 概要 |
|----------|------|------|
| GET | `/api/workouts?date=` | 指定日の運動記録（未指定なら全件） |
| POST | `/api/workouts` | 運動記録追加 |
| DELETE | `/api/workouts/{id}` | 運動記録削除 |
| GET | `/api/workout-master` | プリセットワークアウト一覧 |

### 日次サマリー
| メソッド | パス | 概要 |
|----------|------|------|
| GET | `/api/summary?date=` | 摂取・消費・収支・残 kcal |

### 目標
| メソッド | パス | 概要 |
|----------|------|------|
| GET | `/api/goal` | 目標取得 |
| PUT | `/api/goal` | 目標更新 |

### AI レコメンド
| メソッド | パス | 概要 |
|----------|------|------|
| POST | `/api/ai/recommend-meal` | 献立レコメンド |
| POST | `/api/ai/recommend-workout` | ワークアウトレコメンド |

## ディレクトリ構成

```
api/
├── cmd/server/main.go          # エントリポイント・サーバー起動
├── internal/
│   ├── handler/                # HTTP ハンドラ（エンドポイントごとにファイル分割）
│   │   ├── handler.go          # RegisterAll（全ハンドラをまとめて登録）
│   │   ├── health.go
│   │   ├── weight.go
│   │   ├── meals.go
│   │   ├── workouts.go
│   │   ├── masters.go
│   │   ├── summary.go
│   │   ├── goal.go
│   │   └── ai.go               # Gemini レコメンド（現状スタブ）
│   ├── repository/
│   │   └── memory.go           # インメモリストア（将来 Firestore に差し替え）
│   └── schema/
│       └── schema.go           # 全データモデル定義（Huma が OpenAPI スキーマを自動生成）
├── Dockerfile                  # マルチステージビルド → distroless イメージ
├── go.mod
└── go.sum
```

## 現状の制約・TODO

| 項目 | 現状 | TODO |
|------|------|------|
| データ永続化 | **インメモリ**（再起動でリセット） | Firestore に差し替え |
| AI レコメンド | **スタブ**（固定ダミーデータを返す） | Gemini API (`google.golang.org/genai`) を呼び出す |
| 日次摂取目標 | **固定値** 1800 kcal | 目標体重・目標日から動的に算出（要件 §5.6） |
| TS 型 | web 側で手書き | `/openapi.json` から `openapi-typescript` で自動生成に切り替え |

## テスト

```bash
go test ./...
```

ハンドラのテストは `humatest` を使いサーバーを立てずに HTTP レベルで検証する（`internal/handler/handler_test.go`）。

## ビルド / デプロイ

```bash
# ローカルで Docker ビルド
docker build -t cutagent-api .

# Cloud Run へは GitHub Actions (deploy.yml) が main push 時に自動デプロイ
```
