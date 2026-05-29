# CUTAGENT Web

React + Vite + Tailwind CSS 製フロントエンド。モバイルファースト PWA を目標に構築。

## 起動

```bash
# 依存インストール（初回のみ）
pnpm install

# 開発サーバー起動（バックエンド api/ も別途起動しておくこと）
pnpm dev
```

`http://localhost:5173` で確認できる。`/api/*` へのリクエストは開発時のみ `http://localhost:8080` にプロキシされる（`vite.config.ts`）。

## コマンド一覧

| コマンド | 内容 |
|----------|------|
| `pnpm dev` | 開発サーバー起動 |
| `pnpm build` | 型チェック + 本番ビルド（`dist/` に出力） |
| `pnpm preview` | `dist/` をローカルでプレビュー |
| `pnpm lint` | Biome で lint チェック |
| `pnpm format` | Biome で自動フォーマット |

## ディレクトリ構成

```
web/
├── src/
│   ├── main.tsx            # エントリポイント
│   ├── App.tsx             # ルーティング（ページ切り替え）+ ボトムナビ
│   ├── api.ts              # バックエンド API クライアント（全エンドポイントをまとめて管理）
│   ├── index.css           # グローバルスタイル（Tailwind インポート）
│   ├── vite-env.d.ts       # Vite の型定義
│   └── pages/
│       ├── Dashboard.tsx   # ダッシュボード（摂取・消費・収支カード、残日数）
│       ├── Meals.tsx       # 食事記録（プリセット選択 + 当日一覧）
│       ├── Workouts.tsx    # 運動記録（プリセット選択 + 当日一覧）
│       ├── Weight.tsx      # 体重記録（入力 + 推移一覧）
│       ├── AIRecommend.tsx # AI レコメンド（献立 / ワークアウト提案）
│       └── Settings.tsx    # 設定（目標体重・目標日）
├── index.html
├── vite.config.ts
├── tsconfig.json
├── biome.json
└── package.json
```

## API クライアント（`src/api.ts`）

バックエンドとの通信はすべて `src/api.ts` の `api` オブジェクト経由。型定義もここに集約している。

```ts
// 使用例
import api from "./api";

const goal = await api.getGoal();
const meals = await api.listMeals("2026-05-29");
await api.addMeal({ date: "2026-05-29", mealType: "lunch", name: "サラダ", calories: 350 });
```

**現状の制約:** 型は手書き（`api/internal/schema/schema.go` の構造に準拠）。  
**TODO:** バックエンドの `/openapi.json` から `openapi-typescript` で自動生成に切り替える。

## 技術スタック

| ライブラリ | バージョン | 用途 |
|------------|-----------|------|
| React | 19 | UI |
| Vite | 8 | バンドラ・開発サーバー |
| Tailwind CSS | v4 | スタイリング |
| TypeScript | 6 | 型チェック |
| Biome | latest | lint / format |
| pnpm | 10 | パッケージマネージャ |
