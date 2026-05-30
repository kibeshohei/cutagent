# ADR-0004: Node パッケージマネージャに pnpm を採用

- **Status**: Accepted
- **Date**: 2026-05-28

## 文脈 (Context)

`web/` で React + Vite + Tailwind を扱うため Node の依存を管理する必要がある。npm / yarn / pnpm / bun から選ぶ。

## 決定 (Decision)

**pnpm** を採用する。CLI は [ADR-0005](0005-use-nix-flake.md) の Nix flake で固定。

## 根拠 (Rationale)

- インストールが速い (hardlink ベース、CI も含めて短縮効果あり)
- ディスク効率がよい (グローバルストアを hardlink で参照)
- workspaces 機能がモノレポ拡張時にも有効
- `package.json` の `packageManager` フィールドで Node 24 以降のネイティブサポート対象

## 検討した代替案 (Alternatives)

| 候補 | 採用しなかった理由 |
|---|---|
| npm | 同梱で無設定だが、インストールが遅くディスク非効率 |
| yarn (Classic / Berry) | バージョン分岐があり選定がややこしい。Berry の Plug'n'Play は Vite 等と相性問題が出る |
| bun | 高速だが Node API 互換性に未解決のケースあり。ハッカソン用途で攻めすぎ |

## トレードオフ・残課題 (Consequences)

- **得たもの**: 高速インストール、ディスク効率、ロックファイル (`pnpm-lock.yaml`) の安定性
- **諦めたもの**: npm 標準でないため一部チュートリアル/README が npm 想定で書かれていることがある (大体は `pnpm` に読み替えれば動く)
- **将来見直す条件**: モノレポ機能が pnpm で不足するか、他ツールチェーンが pnpm 非対応の場合 (現状は問題なし)
