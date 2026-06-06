# CUTAGENT — 作業引き継ぎメモ

> このファイルは新しい Claude セッションが直前の作業を把握するためのメモ。
> 作業するたびに「完了」「次やること」を更新すること。

---

## 現在の状態（2026-06-06 時点）

**ブランチ:** `main`
（新規作業時は `<type>/<kebab-case>` の feature ブランチを切る。`main` 直 push 禁止、PR 経由のみ）

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
- AI レコメンドは **Gemini 2.5 Flash 呼び出し**（`google.golang.org/genai`）。`GEMINI_API_KEY` 未設定や API エラー / プロンプト読み込み失敗 / JSON パース失敗時はスタブ候補にフォールバック
- プロンプトは `PROMPTS_DIR` 環境変数で指定（デフォルト `prompts`、Dockerfile では `/app/prompts`）
- `PORT` 環境変数で待受ポートを変更可能（Cloud Run 対応）
- `api/Dockerfile` で distroless イメージにビルド。context はリポルート（`docker build -f api/Dockerfile .`）で、`prompts/` も内包
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
- DB は **Firestore**、AI は **Gemini API**、バックエンドは **Cloud Run**
- フロント（React SPA）の配信は **Firebase Hosting**（GCP 系で統一、CDN 付き、無料枠あり）
- Cloud Run デプロイは Workload Identity Federation を使用（サービスアカウントキーを使わない）

---

## 次にやること（優先順）

### 🔲 1. Firestore 連携
`api/internal/repository/memory.go` を Firestore 実装に差し替える。
ローカル開発は Firestore Emulator を使う（`flake.nix` に追加必要）。

### 🔲 2. GCP セットアップ
Cloud Run サービス / Artifact Registry / Workload Identity Federation の初期設定。
`deploy.yml` が使う Secrets（`GCP_PROJECT_ID`, `WIF_PROVIDER`, `WIF_SERVICE_ACCOUNT`）を GitHub に登録。
Gemini API キーは Secret Manager の `gemini-api-key` に登録し、`deploy.yml` が `--set-secrets` で注入する設計（[docs/gcp-setup.md](docs/gcp-setup.md) 参照）。

### 🔲 3. OpenAPI → TS 型自動生成
Huma が出力する `/openapi.json` から `openapi-typescript` 等で型を生成し `web/src/api.ts` の手書き型を置き換える。

### 🔲 4. Evals CI ゲート
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
- [x] `claude/repository-next-steps-WJga1` を PR #2 として main にマージ（`f054151`）
- [x] ハーネス拡充ドキュメント一式（`docs/architecture.md` / `docs/adr/` (README + template + 0001〜0005) / `prompts/README.md` / `evals/README.md` / `docs/gcp-setup.md`）— PR #4 マージ（`2fc2a0c`）
- [x] AI レコメンドを Gemini 2.5 Flash 実装に差し替え（`google.golang.org/genai`、スタブフォールバック付き、`PROMPTS_DIR` env、Dockerfile context をリポルートに変更）— PR #5 マージ（`f01ed81`）
- [x] 開発時起動の手動 2 ターミナル方式を README に整備（mprocs を一度試して撤回、PR #6 / #7 / #8）
- [x] 体重ページが空状態で真っ白になるバグを修正（`ListWeights` が nil 起点で JSON `null` になっていたのを `[]` 起点に統一、回帰テスト追加）— PR #9 マージ（`3a7f9d2`）
- [x] PC 表示で md (≥ 768px) 以上はサイドナビ + 2 カラム、モバイルは従来のボトムナビという責務分担に切り替え（`web/src/App.tsx` 1 ファイル変更、各ページは触らず）— PR #11 マージ（`0815252`）
- [x] React Router v7 を導入し URL ベースのページ遷移に対応（`web/src/App.tsx` を `BrowserRouter` + `<Routes>` + `<NavLink>` ベースに書き換え、未知の path は `/` へフォールバック、`useState<Page>` を削除）— PR #13 マージ（`1fc9197`）
- [x] ホームの目標進捗表記を「{現在} kg → {目標} kg（あと N kg 減量/増量）」+ 「期日 {日付}」の 2 行に変更（旧「残り N kg → 目標 kg（日付）」が矢印の意味を取り違えやすかった点を解消、増量・目標達成も自然な日本語に）— PR #14 マージ（`088d41d`）
- [x] `api/README.md` を最新実装に揃え Jr 向けに構造説明を追加（PR #5 で Gemini 実装に切り替わったのに「現状スタブ」のままだった記述を修正、レイヤ依存方向の図、データフローのシーケンス、Huma の役割解説、AI レコメンド専用セクション、関連ドキュメントへのリンク、TODO 表の最新化）— PR #16 マージ（`db1f05e`）
