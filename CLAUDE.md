# CLAUDE.md (プロジェクト用)

このドキュメントは、グローバル `~/.claude/CLAUDE.md` の上に、CUTAGENT 固有の方針を追加するもの。

> **セッション開始時は必ず [PLAN.md](PLAN.md) を読んで現在の状態と次のタスクを把握すること。** 読み書きの手順は `plan` Skill ([.claude/skills/plan/SKILL.md](.claude/skills/plan/SKILL.md)) を参照。

## プロダクト概要

期限付きダイエット目標を支援する PWA。Google Cloud DevOps AI Agent Hackathon 2026 出展用途。詳細は [docs/requirements.md](docs/requirements.md) を参照。

## 技術スタック

- フロントエンド: React + Vite + Tailwind CSS / Biome / pnpm
- バックエンド: Go 1.24 + net/http + Huma / gofmt + goimports / golangci-lint
- AI: Gemini API (`google.golang.org/genai`)
- DB: Firestore
- ホスティング: Cloud Run
- 開発環境: Nix flake + direnv

## ディレクトリ規約

```
cutagent/
├── web/                      # React SPA (Vite)
├── api/                      # Go バックエンド
│   ├── cmd/server/
│   └── internal/
│       ├── handler/
│       ├── service/
│       ├── repository/
│       └── schema/
├── prompts/                  # Gemini プロンプト (Markdown)
├── evals/                    # Evals データセット (JSONL)
├── docs/                     # 設計ドキュメント
├── .github/workflows/
├── flake.nix / .envrc
└── .editorconfig
```

## Git / コミット規約

- ブランチ命名: `<type>/<kebab-case>` (type: `feature` / `fix` / `chore` / `docs` / `refactor` / `test` / `style`)
- コミットメッセージ: Conventional Commits 日本語。例: `feat: 体重ログ API を追加`
- 末尾に `Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>` を付与
- マージ方式: Squash merge
- `main` への変更は必ず PR 経由
- PR テンプレ: [.github/PULL_REQUEST_TEMPLATE.md](.github/PULL_REQUEST_TEMPLATE.md)

## コード規約

### Go (`api/`)
- フォーマット: `gofmt` / `goimports`
- リンタ: `golangci-lint` (default ルールセット)
- テスト: `go test` (stdlib)
- 依存追加 (`go get`, `go mod tidy`) は事前にユーザー確認

### TypeScript (`web/`)
- Biome (lint + format)
- `tsc --noEmit` で型チェック
- パッケージマネージャ: **pnpm**。依存追加 (`pnpm add`) は事前にユーザー確認

### 共通
- ファイル末尾改行、LF (`.editorconfig` で統一)
- コメント方針: WHY が非自明な時のみ (グローバル CLAUDE.md 準拠)

## ツール許可ポリシー

[.claude/settings.json](.claude/settings.json) で `allow` / `deny` を管理。
新規にツールを使いたい時は事前にユーザー確認。

## 開発環境

- `direnv allow` 後、`nix develop` が自動起動する想定
- 直接 `brew install`, `npm install -g`, `go install`（システム書込）で環境を汚さない。**Nix flake を真実とする**
- Python は使用しない (BE は Go)

## やらないこと

- 機能を勝手に追加しない。要件定義書 (v0.6) に書かれていないものは将来対応 (§13)
- スコープが見えない時は要件定義書の更新を提案する
- `main` へ直 push しない
- `--no-verify`, `git push --force`, `git reset --hard` は禁止
