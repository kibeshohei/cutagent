# Architecture Decision Records (ADR)

このディレクトリは CUTAGENT の重要な技術的判断を **1 決定 1 ファイル**で記録するもの。
「なぜそう決めたか (Why)」を残すのが目的で、コードや要件定義書 (What) では伝わらない判断軸を補う。

## 一覧

| # | タイトル | ステータス | 日付 |
|---|---|---|---|
| [0001](0001-use-go-for-backend.md) | バックエンドに Go を採用 | Accepted | 2026-05-28 |
| [0002](0002-use-huma-v2-31.md) | API フレームワークに Huma v2.31 を採用 | Accepted | 2026-05-28 |
| [0003](0003-use-firebase-hosting.md) | SPA 配信に Firebase Hosting を採用 | Accepted | 2026-05-29 |
| [0004](0004-use-pnpm.md) | Node パッケージマネージャに pnpm を採用 | Accepted | 2026-05-28 |
| [0005](0005-use-nix-flake.md) | 開発環境に Nix flake を採用 | Accepted | 2026-05-28 |

## 新しい ADR を書くとき

1. [`template.md`](template.md) をコピーして連番ファイルを作る (`0006-...md`)
2. ステータスは `Proposed` で開始、合意後 `Accepted` に更新
3. 上記の一覧表にエントリを追加
4. 既存の判断を覆す場合は古い ADR を `Superseded by #NNNN` に更新

## ステータスの意味

- **Proposed**: 提案中。まだ実装に反映されていない
- **Accepted**: 採用。コードや構成に反映済み
- **Deprecated**: 古くなった。新しい判断に置き換えられた
- **Superseded by #NNNN**: 別の ADR に置き換えられた (リンクを残す)
