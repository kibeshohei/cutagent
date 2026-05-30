# ADR-0005: 開発環境に Nix flake を採用

- **Status**: Accepted
- **Date**: 2026-05-28

## 文脈 (Context)

複数人 (自分 / 相方 / CI) が同じ Go・Node・pnpm・golangci-lint・biome・gcloud のバージョンで開発できる必要がある。マシンを汚さず、半年後も同じ環境を再現できる仕組みが欲しい。開発者は既に nix-darwin + Home Manager を日常運用している。

## 決定 (Decision)

**Nix flake (`flake.nix`) + direnv (`.envrc` に `use flake`)** を採用する。flake.lock で全パッケージのバージョンを固定する。

## 根拠 (Rationale)

- 言語ランタイム (Go / Node) と CLI ツール (pnpm / biome / golangci-lint / gcloud) を**同じ宣言で固定**できる (mise や asdf は CLI まで届かない)
- ホストマシンを汚さない (`/nix/store/` 内に閉じる、`brew install` 不要)
- macOS / Linux で同じ flake が動くため CI でも再利用可能
- `direnv allow` 後はターミナルを開くだけで dev shell が起動する
- 開発者の dotfiles エコシステム (nix-darwin) と統合済み

## 検討した代替案 (Alternatives)

| 候補 | 採用しなかった理由 |
|---|---|
| Docker + devcontainer | macOS でファイル I/O が遅い (ホットリロード劣化)、エディタ依存 (VS Code 中心) |
| Docker compose のみ | 言語ツールがホスト側で別管理になり、バージョン揃わない |
| mise (旧 rtx) / asdf | 言語ランタイム管理に特化。gcloud / biome 等の CLI は別途必要 |
| Homebrew + Brewfile | macOS 限定、バージョン固定が弱い (formula の最新を引く) |
| 完全 Nix 化 (node2nix 等) | `node_modules` まで Nix で作るのは生成物が膨大でメンテ困難 |

## トレードオフ・残課題 (Consequences)

- **得たもの**: 完全再現可能な開発環境、マシン汚染なし、CI でも再利用可能
- **諦めたもの**: Nix の学習コスト (開発者は既習)、初回 google-cloud-sdk のフェッチが重い (数分〜数十分)
- **本番イメージとの関係**: Nix は開発用、本番 Cloud Run イメージは別に `api/Dockerfile` でビルドする (役割が違うレイヤ)
- **将来見直す条件**: チームに Nix 未習熟メンバーが入って学習コストが上回る場合
