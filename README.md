# CUTAGENT

期限付きダイエット目標を支援する PWA。Gemini が献立とワークアウトをレコメンドし、目標達成を後押しする。

> Google Cloud DevOps AI Agent Hackathon 2026 出展用途

## 技術スタック

- **フロントエンド**: React + Vite + Tailwind CSS (`web/`)
- **バックエンド**: Go 1.24 + net/http + [Huma](https://huma.rocks/) (`api/`)
- **AI**: Gemini API
- **DB**: Firestore
- **ホスティング**: Cloud Run

詳細は [docs/requirements.md](docs/requirements.md) を参照。

## セットアップ

### 前提

- [Nix](https://nixos.org/download) (Determinate Nix を推奨)
- [direnv](https://direnv.net/)

### 手順

```sh
git clone <this-repo>
cd cutagent
direnv allow
```

`direnv allow` を実行すると、`.envrc` の `use flake` により [`flake.nix`](flake.nix) が起動し、以下が固定バージョンで揃う。

- Go 1.24 / `gopls` / `goimports` / `golangci-lint`
- Node.js 22 LTS / pnpm
- Biome
- Google Cloud SDK

直接シェルに入りたい場合:

```sh
nix develop
```

> 初回起動はパッケージのフェッチで数十分かかる場合あり。

## ディレクトリ構成

```
cutagent/
├── web/        # フロントエンド (React + Vite)
├── api/        # バックエンド (Go + Huma)
├── prompts/    # Gemini プロンプト (Markdown)
├── evals/      # Evals データセット (JSONL)
├── docs/       # 設計ドキュメント
├── flake.nix   # Nix dev shell 定義
└── .envrc      # direnv → flake 自動起動
```

> 現時点では `web/` / `api/` / `prompts/` / `evals/` は未作成。M1 以降で順次追加。

## ドキュメント

- [要件定義書](docs/requirements.md)
- [プロジェクト用 CLAUDE.md](CLAUDE.md) — 開発規約・コード規約
- [PR テンプレート](.github/PULL_REQUEST_TEMPLATE.md)
