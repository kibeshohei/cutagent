# ADR-0002: API フレームワークに Huma v2.31 を採用

- **Status**: Accepted
- **Date**: 2026-05-28

## 文脈 (Context)

[ADR-0001](0001-use-go-for-backend.md) で Go を採用したが、Go 標準の `net/http` だけだと OpenAPI 仕様書を自動生成できず、web 側 (TypeScript) に型を渡せない。「型の二重管理を避ける」ことが疎結合の鍵なので、OpenAPI を自動生成できるフレームワークが要る。

## 決定 (Decision)

**Huma v2.31** を `net/http` の上で使う。

## 根拠 (Rationale)

- Go の構造体タグから **OpenAPI 3.1 仕様書を自動生成** できる
- Pydantic に近い宣言的バリデーション (`required`, `minimum`, `format` などをタグで指定)
- `net/http` 互換のミドルウェアがそのまま使える (Go 1.22+ の `ServeMux` 強化機能と相性◎)
- 軽量で依存が少ない
- Gin/Echo のような独自のルーティング DSL に縛られない

## 検討した代替案 (Alternatives)

| 候補 | 採用しなかった理由 |
|---|---|
| 素の `net/http` のみ | OpenAPI を手書きする羽目になり、web 側との型同期が破綻する |
| Gin / Echo / Fiber | 独自のミドルウェア体系を学ぶ必要があり、OpenAPI は別途プラグインが必要 |
| go-swagger | 強力だがコードジェネレーションが煩雑 |
| Chi + 手動 OpenAPI | 軽量だが OpenAPI と Go コードが二重管理になる |

## トレードオフ・残課題 (Consequences)

- **得たもの**: OpenAPI 自動生成、型安全な DTO、`net/http` ミドルウェアの再利用
- **諦めたもの**: Echo/Gin に比べエコシステムが小さい (ドキュメントは公式が充実しているので大きな痛みなし)
- **バージョン固定の理由**: Huma v2.38 は Go 1.25 以上を要求する。本プロジェクトは Go 1.24 ([ADR-0001](0001-use-go-for-backend.md)) なので v2.31 を使用
- **将来見直す条件**: Go を 1.25 以上に上げる時に Huma も最新へ
