# 要件定義書 — CUTAGENT

> バージョン: 0.6.0（下書き）
> 更新日: 2026-05-29
> ステータス: レビュー中
> 主用途: Google Cloud DevOps AI Agent Hackathon 2026 出展

---

## 0. このドキュメントの位置づけ

ハッカソン出展を前提に、v0.4 で盛り込みすぎた機能を MVP まで絞り込んだ版（v0.5）に対して、**バックエンドの言語/フレームワークを Go (net/http + Huma) に変更**したもの。Cloud Run/Firestore とのインフラ親和性、起動速度、コンテナサイズ、web 側との完全疎結合を優先した。

---

## 1. プロダクト概要

### 1.1 ビジョン
「目標体重と目標日を設定するだけで、AIが食事・運動のレコメンドを通じてゴール達成を支援するダイエットアプリ」

### 1.2 開発目的
- 自分自身のダイエット管理ツールとして実際に使う
- Google Cloud DevOps AI Agent Hackathon 2026 への出展
- 個人開発の実績として残す

### 1.3 ターゲットユーザー
- 期限付きの目標体重がある人（イベント・大会・結婚式など）
- 食事管理を続けたいが手間をかけたくない人

---

## 2. ハッカソン要件マッピング

### 2.1 必須要件の対応
| 必須項目 | 本プロダクトでの選定 |
|---|---|
| Google Cloud アプリ実行プロダクト | **Cloud Run**（バックエンド） |
| Google Cloud AI 技術 | **Gemini API**（献立・ワークアウトレコメンド） |

### 2.2 任意で利用する技術
- **Firestore**: DB
- **Secret Manager**: API キー管理
- **GitHub Actions**: CI/CD + Evals ゲート
- **Firebase Hosting**: React SPA の静的ファイル配信（CDN 付き）

### 2.3 審査基準との対応
| 審査基準 | 本プロダクトでどう満たすか |
|---|---|
| AIエージェントが価値の中心になっているか | Gemini を中核に、目標と日次収支を踏まえた**目的合致レコメンド**を提供（自律的振る舞いは将来拡張） |
| 課題へのアプローチ力 | 「期限付き目標 × 手間ゼロ」を軸にターゲットと提供価値を一貫させる |
| ユーザビリティ | モバイルファースト・PWA、プリセット選択中心の入力で負荷を最小化 |
| 実用性・体験価値 | 開発者自身が日次で利用 |
| 実装力 | Go + Cloud Run + Gemini + Evals/CI で**継続改善できる本番品質**を提示 |

---

## 3. AI レコメンド機能

### 3.1 機能の位置づけ
シンプルな単発レコメンドとして実装する。自律ループ・通知・自動リプランは MVP のスコープ外（§13 を参照）。

### 3.2 レコメンドの種類
| 名称 | 入力コンテキスト | 出力 |
|---|---|---|
| 献立レコメンド | 目標、現在体重、当日の残カロリー、簡易な好み | 1〜複数の食事候補（メニュー名 + 概算カロリー） |
| ワークアウトレコメンド | 目標、当日の収支、簡易な好み | 1〜複数の運動候補（種目 + 時間 + 概算消費kcal） |

### 3.3 セキュリティ
- Gemini API キーはバックエンドのみで保持。Cloud Run の Secret Manager 統合で注入

---

## 4. DevOps ループ「まわす」

### 4.1 プロンプトのバージョン管理
- プロンプトを `prompts/` 配下に Markdown で配置し、Git で履歴管理

### 4.2 Evals データセット
- `evals/` に評価ケースを JSON で蓄積（入力 + 期待挙動 + 採点ルーブリック）
- カテゴリ: 献立レコメンド品質 / ワークアウトレコメンド品質

### 4.3 PR 時の自動評価（GitHub Actions）
- PR 作成時に Evals を実行し、LLM-as-a-Judge で各ケースをスコアリング
- スコア低下時は CI を fail とし、マージをブロック
- スコアレポートを PR コメントに自動投稿

### 4.4 自動デプロイ
- `main` への push で Cloud Run に自動デプロイ

---

## 5. 機能要件

### 5.1 ダッシュボード
- 当日の摂取・消費・収支をカード表示
- 体重推移グラフ（目標ラインと重ねて表示）
- 目標日まで残り日数・残り減量 kg
- ペース警告（このままでは間に合わない場合にアラート表示）

### 5.2 体重ログ
- 毎日の体重入力（小数点1桁）
- 体重推移グラフ

### 5.3 食事記録
- プリセットリストから選択して記録（朝/昼/夜/間食の区分）
- 1日の摂取カロリー合計を自動集計

### 5.4 運動記録
- プリセットリストから選択して記録
- 消費カロリーは固定値（体重連動の自動計算は将来対応)
- 1日の消費カロリー合計を自動集計

### 5.5 AI レコメンド
- 献立レコメンド（ボタン1回で生成）
- ワークアウトレコメンド（ボタン1回で生成）

### 5.6 目標設定
- 現在体重・目標体重・目標日の設定
- 必要な1日あたりの赤字カロリーを自動計算

---

## 6. 非機能要件

| 項目 | 内容 |
|------|------|
| パフォーマンス | 初回ロード3秒以内、AI レスポンス10秒以内、API のコールドスタート 500ms 以内 |
| デバイス | モバイルファースト・PWA 対応 |
| 将来対応 | React Native でネイティブアプリ化できる設計 |
| セキュリティ | Gemini API キーはバックエンドで管理（Secret Manager） |
| 疎結合 | web と api は OpenAPI 経由のみで結合。TS の型は OpenAPI からの生成物を使う |

---

## 7. 画面構成

```
├── ダッシュボード（ホーム）
├── 食事記録
│   └── 今日の記録一覧 + プリセットから追加
├── 運動記録
│   └── 今日の記録一覧 + プリセットから追加
├── 体重ログ
│   └── グラフ + 記録一覧
├── AI レコメンド
│   ├── 献立提案
│   └── ワークアウト提案
└── 設定
    └── 目標設定（体重・日付）
```

---

## 8. API インターフェース設計（案）

API は Huma により OpenAPI 3.1 仕様書を自動生成する。エンドポイント設計の案は以下。

```
# 体重
GET    /api/weight                体重記録一覧
POST   /api/weight                体重記録追加

# 食事
GET    /api/meals?date=           指定日の食事記録
POST   /api/meals                 食事記録追加
DELETE /api/meals/:id             食事記録削除

# 食べるものリスト（プリセット）
GET    /api/meal-master           プリセット食品一覧

# 運動
GET    /api/workouts?date=        指定日の運動記録
POST   /api/workouts              運動記録追加
DELETE /api/workouts/:id          運動記録削除

# ワークアウトリスト（プリセット）
GET    /api/workout-master        プリセットワークアウト一覧

# 日次サマリー
GET    /api/summary?date=         指定日の intake / burned / 収支 / 残kcal

# 目標
GET    /api/goal                  目標取得
PUT    /api/goal                  目標更新

# AI レコメンド
POST   /api/ai/recommend-meal     献立レコメンド
POST   /api/ai/recommend-workout  ワークアウトレコメンド
```

---

## 9. データモデル（案）

実装は Go の構造体で定義し、Huma が OpenAPI スキーマを自動生成する。`web/` 側は生成された OpenAPI から TS 型を生成する（型の二重管理を避ける）。下記は概念表記。

```typescript
type Goal = {
  targetWeight: number;     // 目標体重 (kg)
  targetDate: string;       // 目標日 (ISO 8601)
  currentWeight: number;    // 現在体重 (kg)
};

type WeightLog = {
  id: string;
  date: string;
  weight: number;
};

type MealLog = {
  id: string;
  date: string;
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  name: string;
  calories: number;
};

type WorkoutLog = {
  id: string;
  date: string;
  name: string;
  durationMinutes: number;
  caloriesBurned: number;
};

type MealMaster = {
  id: string;
  name: string;
  calories: number;
};

type WorkoutMaster = {
  id: string;
  name: string;
  caloriesPerHour: number;
};
```

---

## 10. 技術スタック（確定）

| レイヤー | 技術 | 備考 |
|---------|------|------|
| フロントエンド | React + Vite + Tailwind CSS | PWA 対応 |
| バックエンド | **Go 1.24 + net/http + Huma** | OpenAPI 3.1 自動生成、型安全な DTO |
| AI | **Gemini API** (`google.golang.org/genai`) | バックエンドから呼び出し |
| DB | **Firestore** (Go SDK) | サーバーレス、Cloud Run と統合容易 |
| 認証 | なし（MVP） | ローカル仮ユーザーで進める。将来 Firebase Auth |
| シークレット | **Secret Manager** | Gemini API キー等 |
| ホスティング | **Cloud Run** | 起動時間 500ms 以内を目安 |
| CI/CD | **GitHub Actions** | 自動デプロイ + Evals ゲート |
| 開発環境 | **Nix (flake) + direnv** | Go、Node、各種 CLI のバージョンを固定 |

---

## 11. マイルストーン

| ID | 目標 | 完了条件 |
|----|------|----------|
| M0 | 要件定義 v0.6 確定 | 本ドキュメント承認 |
| M1 | 最小疎通 | Cloud Run に Hello API（Go）がデプロイされ、GitHub Actions で自動デプロイされる |
| M2 | コア CRUD | 食事・運動・体重・目標の記録と日次収支の集計が動く。OpenAPI が生成され、`web/` 側で型が利用可能 |
| M3 | AI レコメンド | 献立とワークアウトのレコメンドが Gemini 経由で返る |
| M4 | Evals パイプライン | PR 時に Evals が走り、スコア低下でマージブロックされる |
| M5 | 仕上げ | UI 磨き込み・デモ動画・README 完了 |

---

## 12. 未決事項

- [ ] プリセットの食べるものリスト・ワークアウトリストの初期データ内容
- [ ] 目標カロリー計算式（基礎代謝×活動係数 or シンプルな赤字目標）
- [ ] Evals の採点ルーブリックの初版定義
- [ ] Huma のミドルウェア構成（認証、ロギング、CORS 等）

---

## 13. 将来対応（スコープ外）

エージェントである必然性を強めるため、基本機能の確立後に以下を順次追加する。

- 自律ループ（Cloud Scheduler + Pub/Sub による日次ペース判定）
- ペース逸脱時の自動リプラン（献立修正案／運動追加案／期日延伸案の比較）
- 通知（PWA Push などによる能動的通知）
- 写真→食事推定（Gemini multimodal）
- 自然文入力（記録・目標設定の NLP パース）
- カスタム食品・ワークアウトの管理
- Firebase Auth でのマルチユーザー対応
- ワークアウト消費カロリーの体重連動自動計算（体重 × METs）
- React Native によるネイティブアプリ化

---

*このドキュメントは随時更新します。*
