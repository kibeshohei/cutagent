# evals/

CUTAGENT の AI レコメンド品質を継続評価するためのデータセット。
**1 ユースケース 1 JSONL ファイル**、1 行 = 1 評価ケース。

## ファイル

| ファイル | 評価対象プロンプト |
|---|---|
| [`recommend-meal.jsonl`](recommend-meal.jsonl) | [`prompts/recommend-meal.md`](../prompts/recommend-meal.md) |
| [`recommend-workout.jsonl`](recommend-workout.jsonl) | [`prompts/recommend-workout.md`](../prompts/recommend-workout.md) |

## ケースのスキーマ (JSONL 1 行)

```json
{
  "id": "meal-001-low-cal-dinner",
  "description": "夜・残カロリー300 で重い食事をすすめないか",
  "input": {
    "goal": { "targetWeight": 65, "targetDate": "2026-08-01", "currentWeight": 70 },
    "remainingKcal": 300,
    "mealType": "dinner",
    "preferences": ["高タンパク", "低脂質"]
  },
  "expected": {
    "shouldContainAny": ["鶏むね", "卵白", "白身魚"],
    "shouldNotExceedKcal": 350
  },
  "rubric": [
    "残カロリー(300kcal)を超えていないか",
    "タンパク質源が含まれているか",
    "好みの低脂質に合うか"
  ]
}
```

| フィールド | 必須 | 意味 |
|---|---|---|
| `id` | ✅ | ユニーク。`<topic>-<連番>-<簡潔な要約>` |
| `description` | ✅ | このケースで検証したいことを 1 文で |
| `input` | ✅ | API に渡す入力 (型はプロンプトに対応) |
| `expected` | ✅ | **機械的に判定できる**期待値 |
| `rubric` | ✅ | LLM-as-a-Judge が評価する観点 (3〜5 個) |

## 評価方法 (2 段階)

### 1. 機械チェック
`expected.shouldContainAny` / `shouldNotExceedKcal` などはコードで判定。Pass / Fail が明確。

### 2. LLM-as-a-Judge
Gemini の出力 + `rubric` を別の Gemini 呼び出しに渡し、各観点に 0〜5 でスコアリング。
平均スコアが閾値を下回ると CI fail。

> M4 タスクで CI に組み込む。現状は JSONL の蓄積フェーズ。

## 新しい評価ケースを追加するタイミング

**1 回の本番事故 → 1 ケース追加**を徹底する。

- 開発中に「この出力おかしい」と感じた瞬間
- 本番ログから失敗例を発見したとき
- プロンプトを改善するときは「直したい現象」のケースを先に書く

その入力と「本来こうあってほしい」を JSONL に追記する。プロンプト改善 → そのケースで合格するか確認 → コミット。

## 注意

- **個人情報を含めない** (公開リポジトリ)
- **入力は再現性のあるリテラル**で書く (時刻・乱数を含めない)
- **JSONL なのでコメント不可** — `description` に必ず意図を書く
