# ADR-0001: バックエンドに Go を採用

- **Status**: Accepted
- **Date**: 2026-05-28

## 文脈 (Context)

Google Cloud DevOps AI Agent Hackathon 2026 の出展用バックエンドを選定する。Cloud Run でホストし、Gemini API を呼び、Firestore を使う。「web 側 (React) との完全疎結合」「Google Cloud との親和性」「再現可能な開発環境」が主な要求。

## 決定 (Decision)

**Go 1.24** をバックエンド言語として採用する。

## 根拠 (Rationale)

- Cloud Run のコールドスタート時間がもっとも短いランタイムの一つ (Python の 1/3 以下、500ms 目標を満たしやすい)
- Google Cloud SDK / Firestore / Gemini いずれも公式 Go SDK あり
- 言語シンタックスが小さく、両言語ともほぼ未経験のプロジェクトで学習コストが低い
- web と別言語にすることで「型を共有しない」物理的な疎結合が成立 (OpenAPI 経由のみで結合)
- コンテナイメージが小さく (distroless で数十 MB)、デプロイが速い

## 検討した代替案 (Alternatives)

| 候補 | 採用しなかった理由 |
|---|---|
| Python (FastAPI) | AI/LLM エコシステム最厚だが、Cloud Run コールドスタートが 1〜2 秒。今回 AI は単発レコメンド 2 本のみで、Python のエコシステム優位を使い切らない |
| TypeScript (Hono) | web (React) と同言語になり「疎結合」方針と矛盾。型共有の誘惑も発生する |
| Rust (Axum) | パフォーマンスは最強だが学習コストが高すぎる |
| Kotlin (Ktor) | JVM ベースで起動時間とコンテナサイズが Go より大きい |

## トレードオフ・残課題 (Consequences)

- **得たもの**: 起動の速さ、Cloud Run / Firestore との親和性、小さいイメージ
- **諦めたもの**: Pydantic 級の宣言的バリデーション (Huma で部分的に補う / [ADR-0002](0002-use-huma-v2-31.md))、ADK / Vertex AI Eval の Python ファースト機能
- **将来見直す条件**: AI エージェント機能が大幅に高度化し、Python 側 SDK でしか提供されない機能 (例: ADK の新機能) が必須になった場合
