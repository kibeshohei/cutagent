# CUTAGENT API

Go 製バックエンド。`net/http` + [Huma v2](https://github.com/danielgtaylor/huma) で実装し、構造体タグから OpenAPI 3.1 仕様書を自動生成する。

## 起動

`api/` ディレクトリで:

```bash
go run ./cmd/server
```

- デフォルト待受: `:8080`
- `PORT` で変更可 (Cloud Run はこれで受け取る)
- `GEMINI_API_KEY` を渡すと AI レコメンドが実 API で動作。未設定時はスタブ候補にフォールバック
- プロンプトディレクトリは `PROMPTS_DIR` で指定（デフォルト `prompts`）

起動後、以下で仕様書を確認できる:

| URL | 内容 |
|---|---|
| `http://localhost:8080/docs` | Swagger UI |
| `http://localhost:8080/openapi.json` | OpenAPI 3.1 JSON（TS 型生成に使う） |

## ディレクトリ構成

```
api/
├── cmd/server/main.go          # サーバー起動エントリ
├── internal/
│   ├── handler/                # HTTP の入口 (1 機能 = 1 ファイル)
│   │   ├── handler.go          # RegisterAll で全ハンドラを集約
│   │   ├── health.go
│   │   ├── weight.go
│   │   ├── meals.go
│   │   ├── workouts.go
│   │   ├── masters.go
│   │   ├── summary.go
│   │   ├── goal.go
│   │   └── ai.go               # Gemini レコメンド (フォールバック付き)
│   ├── repository/
│   │   └── memory.go           # インメモリストア (将来 Firestore に差し替え)
│   └── schema/
│       └── schema.go           # 共通データモデル (Goal, WeightLog, MealLog 等)
├── Dockerfile                  # context = リポルート / prompts/ 内包
├── go.mod / go.sum
└── README.md
```

## レイヤと依存方向

```
main.go
   │ 起動して
   ▼
handler/*.go ────→ repository/memory.go
                ↓ どちらも使う
            schema/schema.go
```

| レイヤ | 責務 |
|---|---|
| **main.go** | HTTP サーバ起動、Huma 初期化、graceful shutdown |
| **handler/** | エンドポイント登録、入力パース、repository 呼び出し、レスポンス整形 |
| **repository/** | データ操作 (現状はインメモリ slice) |
| **schema/** | 共通の型定義 (Huma が OpenAPI スキーマに変換) |

DDD 的な厚い service 層は意図的に持たない。MVP なのでビジネスロジックが薄く、handler → repository の 2 層で十分。Firestore 連携時に `Repository` interface を切る予定。

## データの流れ (例: 食事記録 POST)

```
ブラウザ (React)
   │ POST /api/meals { date, mealType, name, calories }
   ▼
[main.go] net/http ServeMux
   ▼
[humago adapter]
   │ JSON パース → 構造体タグでバリデーション (mealType の enum 等)
   ▼
[handler/meals.go]
   │ store.AddMeal(meal)
   ▼
[repository/memory.go]
   │ mu.Lock() → slice に append → mu.Unlock() / ID 自動生成
   ▼
レスポンス: 201 Created + MealLog (ID 付き)
```

## エンドポイント一覧

### システム
| メソッド | パス | 概要 |
|---|---|---|
| GET | `/api/health` | ヘルスチェック |

### 体重
| メソッド | パス | 概要 |
|---|---|---|
| GET | `/api/weight` | 体重記録一覧 |
| POST | `/api/weight` | 体重記録追加 |

### 食事
| メソッド | パス | 概要 |
|---|---|---|
| GET | `/api/meals?date=` | 指定日の食事記録（未指定なら全件） |
| POST | `/api/meals` | 食事記録追加 |
| DELETE | `/api/meals/{id}` | 食事記録削除 |
| GET | `/api/meal-master` | プリセット食品一覧 |

### 運動
| メソッド | パス | 概要 |
|---|---|---|
| GET | `/api/workouts?date=` | 指定日の運動記録（未指定なら全件） |
| POST | `/api/workouts` | 運動記録追加 |
| DELETE | `/api/workouts/{id}` | 運動記録削除 |
| GET | `/api/workout-master` | プリセットワークアウト一覧 |

### 日次サマリー / 目標
| メソッド | パス | 概要 |
|---|---|---|
| GET | `/api/summary?date=` | 摂取・消費・収支・残 kcal |
| GET | `/api/goal` | 目標取得 |
| PUT | `/api/goal` | 目標更新 |

### AI レコメンド
| メソッド | パス | 概要 |
|---|---|---|
| POST | `/api/ai/recommend-meal` | 献立レコメンド (Gemini 2.5 Flash) |
| POST | `/api/ai/recommend-workout` | ワークアウトレコメンド (Gemini 2.5 Flash) |

## Huma で得ているもの

Huma は「構造体タグから OpenAPI を自動生成するフレームワーク」。例:

```go
type MealLog struct {
    ID       string `json:"id,omitempty" doc:"ID（サーバー生成）"`
    Date     string `json:"date" doc:"記録日 (ISO 8601)" example:"2026-05-28"`
    MealType string `json:"mealType" enum:"breakfast,lunch,dinner,snack"`
    Name     string `json:"name" doc:"メニュー名" example:"鶏むね肉のサラダ"`
    Calories int    `json:"calories" doc:"摂取カロリー (kcal)" example:"350"`
}
```

このタグだけで:

- **バリデーション**（`enum` で値を縛る、`format` `minimum` 等も同様）
- **OpenAPI 3.1 ドキュメント**を `/openapi.json` で配信
- **TS 型生成の素材**（将来 `openapi-typescript` で `web/src/api.ts` を自動生成予定）

採用理由は [docs/adr/0002-use-huma-v2-31.md](../docs/adr/0002-use-huma-v2-31.md) を参照。

## AI レコメンド (ai.go)

| 項目 | 内容 |
|---|---|
| モデル | `gemini-2.5-flash` (`google.golang.org/genai`) |
| 認証 | `GEMINI_API_KEY` 環境変数 (本番は Secret Manager 経由で注入) |
| プロンプト | `prompts/recommend-meal.md` / `prompts/recommend-workout.md` を `PROMPTS_DIR` で位置指定 |
| 入力コンテキスト | `store.Goal()` + `store.Summary(date)` を `{{var}}` で埋め込み |
| 出力パース | ` ```json ... ``` ` ブロックを剥がして `json.Unmarshal` |
| **フォールバック** | キー未設定 / プロンプト読込失敗 / API エラー / JSON パース失敗 → スタブ候補を返す |

プロンプトの書き方は [prompts/README.md](../prompts/README.md)、評価の運用は [evals/README.md](../evals/README.md) を参照。

## テスト

```bash
go test ./...
```

- ハンドラのテストは [`humatest`](https://github.com/danielgtaylor/huma) を使い、サーバーを立てずに HTTP レベルで検証 (`internal/handler/handler_test.go`)
- AI のフォールバックパス専用に `ai_test.go` を分離
- 体重ページ空状態の回帰テスト (`TestWeightListEmptyReturnsArray`) で「`/api/weight` が空でも `[]` を返す」ことを担保

## ビルド / デプロイ

```bash
# ローカル Docker ビルド (context は必ずリポルート)
docker build -f api/Dockerfile -t cutagent-api .

# Cloud Run へは GitHub Actions (deploy.yml) が main push 時に自動デプロイ
```

`api/Dockerfile` は **context = リポルート**で動く設計。`prompts/` を `/app/prompts/` にコピーしてイメージへ内包し、`ENV PROMPTS_DIR=/app/prompts` を設定する。本番の段取りは [docs/gcp-setup.md](../docs/gcp-setup.md)。

## 現状の制約・TODO

| 項目 | 現状 | TODO |
|------|------|------|
| データ永続化 | **インメモリ** (再起動でリセット) | Firestore に差し替え + Repository interface 化 |
| 日次摂取目標 | **固定値 1800 kcal** | 目標体重・目標日から動的算出 (要件 §5.6) |
| TS 型 | web 側で手書き (`web/src/api.ts`) | `/openapi.json` から `openapi-typescript` で自動生成 |
| 認証 | なし | Firebase Auth (将来) |
| 自律エージェント | オンデマンドレコメンドのみ | Cloud Scheduler + Pub/Sub で日次評価 (将来、要件 §13) |
