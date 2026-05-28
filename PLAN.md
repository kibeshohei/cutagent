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
| `api/` | Go バックエンド骨格あり（後述） |
| `web/` | **未着手** |
| `prompts/` | **未着手** |
| `evals/` | **未着手** |
| `.github/workflows/ci.yml` | Go の CI あり（build / test / lint） |
| `.claude/skills/plan/` | PLAN.md の運用手順を定義する plan Skill |

### `api/` の中身

- **Huma v2.31** + `net/http` でサーバー起動
- エンドポイントは `GET /api/health` だけ（`{"status":"ok"}` を返す）
- `PORT` 環境変数で待受ポートを変更可能（Cloud Run 対応）
- `api/Dockerfile` で distroless イメージにビルド済み
- `go 1.24.0` / module: `github.com/kibeshohei/cutagent/api`

### 技術的な決定事項（変えないこと）

- Go は **1.24**（Huma v2.38 は Go 1.25 必須のため v2.31 を使用）
- フロントは **React + Vite + Tailwind**、パッケージマネージャは **pnpm**
- OpenAPI は Huma が自動生成 → TS 型はそこから生成（型の二重管理なし）
- DB は **Firestore**、AI は **Gemini API**、ホストは **Cloud Run**

---

## 次にやること（優先順）

### 🔲 1. この PR を main にマージ
`claude/repository-next-steps-WJga1` の PR を作って CI が通るか確認してからマージ。

### 🔲 2. 体重ログ API
`GET /api/weight` / `POST /api/weight` を実装。
Firestore 接続が先に必要（ローカル開発は Firestore Emulator を使う）。

### 🔲 3. 食事・運動記録 API
`GET/POST/DELETE /api/meals` と `/api/workouts`。
プリセットマスター (`/api/meal-master`, `/api/workout-master`) も含む。

### 🔲 4. 日次サマリー API
`GET /api/summary?date=` で摂取/消費/収支をまとめて返す。

### 🔲 5. 目標設定 API
`GET /api/goal` / `PUT /api/goal`。

### 🔲 6. Gemini レコメンド API
`POST /api/ai/recommend-meal` / `/api/ai/recommend-workout`。
Secret Manager から API キーを取得する実装が必要。

### 🔲 7. React フロント骨格 (`web/`)
Vite + React + Tailwind をセットアップ。
OpenAPI から TS 型を生成する仕組みも含む。

### 🔲 8. Cloud Run へのデプロイ CI/CD
`main` push 時に Cloud Run へ自動デプロイする GitHub Actions を追加。

### 🔲 9. Evals / プロンプト管理
`prompts/` に Gemini プロンプトを Markdown で配置。
`evals/` に評価ケース JSONL を作成。

---

## 完了済み

- [x] 要件定義 v0.6
- [x] Nix flake / direnv セットアップ
- [x] CLAUDE.md / settings.json / PR テンプレ / .editorconfig
- [x] Go バックエンド骨格（Huma + `/api/health`）
- [x] api/Dockerfile
- [x] GitHub Actions CI（Go build / test / lint）
- [x] PLAN.md（セッション引き継ぎメモ）と plan Skill（その運用手順）
