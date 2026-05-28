# CUTAGENT — 作業引き継ぎメモ

> このファイルは新しい Claude セッションが直前の作業を把握するためのメモ。
> 作業するたびに「完了」「次やること」を更新すること。

---

## 現在の状態（2026-05-28 時点）

**ブランチ:** `claude/repository-next-steps-WJga1`
（`main` への変更は必ず PR 経由。直 push 禁止）

### 何がある？

| 場所 | 状態 |
|------|------|
| `docs/requirements.md` | 要件定義 v0.6 完成 |
| `api/` | 全 API エンドポイント実装済み（インメモリ） |
| `web/` | React + Vite + Tailwind 骨格完成、全ページあり |
| `prompts/` | 献立・ワークアウトレコメンドのプロンプト v0.1 あり |
| `evals/` | 献立・ワークアウトの評価ケース JSONL あり |
| `.github/workflows/ci.yml` | Go + React の CI あり |
| `.github/workflows/deploy.yml` | Cloud Run デプロイ（main push 時）あり |
| `.claude/skills/plan/` | PLAN.md 運用手順 Skill |

### `api/` の中身

- **Huma v2.31** + `net/http` でサーバー起動
- エンドポイント: `/api/health`, `/api/weight`, `/api/meals`, `/api/workouts`, `/api/meal-master`, `/api/workout-master`, `/api/summary`, `/api/goal`, `/api/ai/recommend-meal`, `/api/ai/recommend-workout`
- データ永続化は**インメモリ**（`api/internal/repository/memory.go`）
- AI レコメンドは**スタブ実装**（固定のダミーデータを返す）
- `PORT` 環境変数で待受ポートを変更可能（Cloud Run 対応）
- `api/Dockerfile` で distroless イメージにビルド
- `go 1.24.0` / module: `github.com/kibeshohei/cutagent/api`

### `web/` の中身

- React 19 + Vite 8 + Tailwind CSS v4
- ページ: ダッシュボード / 食事記録 / 運動記録 / 体重記録 / AI レコメンド / 設定
- ボトムナビゲーション、モバイルファースト
- `vite.config.ts` の proxy で `/api` → `localhost:8080` に転送（開発時）

### 技術的な決定事項（変えないこと）

- Go は **1.24**（Huma v2.38 は Go 1.25 必須のため v2.31 を使用）
- フロントは **React + Vite + Tailwind**、パッケージマネージャは **pnpm**
- OpenAPI は Huma が自動生成 → TS 型はそこから生成予定（現状は手書き）
- DB は **Firestore**、AI は **Gemini API**、ホストは **Cloud Run**
- Cloud Run デプロイは Workload Identity Federation を使用（サービスアカウントキーを使わない）

---

## 次にやること（優先順）

### 🔲 1. この PR を main にマージ
`claude/repository-next-steps-WJga1` の PR を作って CI が通るか確認してからマージ。

### 🔲 2. Firestore 連携
`api/internal/repository/memory.go` を Firestore 実装に差し替える。
ローカル開発は Firestore Emulator を使う（`flake.nix` に追加必要）。

### 🔲 3. Gemini API 連携
`api/internal/handler/ai.go` のスタブを実際の Gemini 呼び出しに差し替える。
`google.golang.org/genai` を使用。API キーは Secret Manager 経由で注入（`GEMINI_API_KEY` 環境変数）。
プロンプトは `prompts/recommend-meal.md` / `prompts/recommend-workout.md` を参照。

### 🔲 4. GCP セットアップ
Cloud Run サービス / Artifact Registry / Workload Identity Federation の初期設定。
`deploy.yml` が使う Secrets（`GCP_PROJECT_ID`, `WIF_PROVIDER`, `WIF_SERVICE_ACCOUNT`）を GitHub に登録。

### 🔲 5. OpenAPI → TS 型自動生成
Huma が出力する `/openapi.json` から `openapi-typescript` 等で型を生成し `web/src/api.ts` の手書き型を置き換える。

### 🔲 6. Evals CI ゲート
`evals/` の JSONL ケースを実行して LLM-as-a-Judge でスコアリングする GitHub Actions を追加。スコア低下時は CI fail。

---

## 完了済み

- [x] 要件定義 v0.6
- [x] Nix flake / direnv セットアップ
- [x] CLAUDE.md / settings.json / PR テンプレ / .editorconfig
- [x] Go バックエンド骨格（Huma + `/api/health`）
- [x] api/Dockerfile
- [x] GitHub Actions CI（Go build / test / lint）
- [x] PLAN.md（セッション引き継ぎメモ）と plan Skill（その運用手順）
- [x] Go 全 API エンドポイント（インメモリ実装・AI スタブ含む）
- [x] React フロント骨格（Vite + Tailwind、全ページ）
- [x] prompts/ と evals/ の初版
- [x] GitHub Actions CI に web ジョブ追加
- [x] Cloud Run デプロイワークフロー（deploy.yml）
