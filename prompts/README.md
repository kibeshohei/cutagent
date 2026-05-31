# prompts/

CUTAGENT の AI レコメンドで使う Gemini プロンプト集。**1 機能 1 ファイル**、Markdown で記述する。

## 現状

| ファイル | 用途 | 呼び出し元 |
|---|---|---|
| [`recommend-meal.md`](recommend-meal.md) | 献立レコメンド | `api/internal/handler/ai.go` の `RecommendMeal` |
| [`recommend-workout.md`](recommend-workout.md) | ワークアウトレコメンド | `api/internal/handler/ai.go` の `RecommendWorkout` |

## プロンプトの構造

各 Markdown は以下のセクションで構成する。

```markdown
# 用途の1行説明

## System
LLM に与える役割・トーン (例: 「あなたは管理栄養士です」)

## User (template)
{{variable}} 形式のプレースホルダで動的入力を埋め込む

## 出力形式 (Output)
JSON Schema や例で出力を縛る
```

呼び出し元 (`api/internal/handler/ai.go`) はこのファイルを読み、テンプレを埋めてから Gemini に送る。

## バージョニング

- バージョン番号は**振らない**。履歴管理は Git に任せる
- 大きな改訂は commit メッセージで意図を明示する (`prompts: 献立の出力を JSON Schema 化`)
- 改訂後は必ず Evals でリグレッションを確認する ([`evals/README.md`](../evals/README.md))

## ベストプラクティス

| 項目 | 指針 |
|---|---|
| 役割固定 | System で「あなたは〇〇の専門家」を明示 |
| 出力制約 | JSON 出力を強制 (パースエラーで stub に fallback する設計) |
| 入力境界 | `<context>...</context>` のような XML 風タグで動的入力を囲む |
| 例示 (few-shot) | Gemini 2.0 以降は zero-shot で足りるケースが多い、足す前に zero-shot で測る |
| 温度 | 0.4〜0.6 を起点に。決定的すぎると単調、高すぎると逸脱する |

## Evals との関係

`evals/<同じ名前>.jsonl` に「この入力ならこうあってほしい」ケースを蓄積する。
プロンプトを変えたら CI で LLM-as-a-Judge を回し、スコアが下がったらマージブロックする (M4 タスク)。

詳細は [`../evals/README.md`](../evals/README.md)。
