# GCP セットアップ手順

CUTAGENT を Cloud Run (バックエンド) + Firebase Hosting (フロント) で稼働させるための初期設定。
`.github/workflows/deploy.yml` が動くために必要な権限・シークレット・Secret Manager 登録を含む。

## 前提

- GCP アカウントと支払いプロファイル
- `gcloud` CLI が認証済み (`gcloud auth login`)
- リポジトリの GitHub Secrets を編集できる権限
- Gemini API キー (Google AI Studio で取得)

---

## 1. プロジェクト作成

```sh
export PROJECT_ID="cutagent-prod"   # 任意の ID
export REGION="asia-northeast1"
gcloud projects create "$PROJECT_ID"
gcloud config set project "$PROJECT_ID"
```

支払い口座をリンクする (コンソールから、または `gcloud billing` で)。

## 2. API 有効化

```sh
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  firestore.googleapis.com \
  iamcredentials.googleapis.com \
  sts.googleapis.com
```

## 3. Artifact Registry リポジトリ作成

```sh
gcloud artifacts repositories create cutagent \
  --repository-format=docker \
  --location="$REGION"
```

## 4. Firestore データベース作成

```sh
gcloud firestore databases create --location="$REGION"
```

> Native モードで作成される。

## 5. Workload Identity Federation 設定

GitHub Actions が**サービスアカウントキーなし**で GCP にアクセスできるようにする ([ADR 検討予定]の設計)。

```sh
# Workload Identity Pool
gcloud iam workload-identity-pools create github \
  --location=global

# Provider (GitHub OIDC)
gcloud iam workload-identity-pools providers create-oidc github-provider \
  --location=global \
  --workload-identity-pool=github \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" \
  --attribute-condition="assertion.repository=='kibeshohei/cutagent'"

# サービスアカウント (デプロイ用)
gcloud iam service-accounts create github-deployer \
  --display-name="GitHub Actions Deployer"
export SA_EMAIL="github-deployer@${PROJECT_ID}.iam.gserviceaccount.com"

# 必要なロールを付与
for role in roles/run.admin roles/artifactregistry.writer roles/iam.serviceAccountUser; do
  gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SA_EMAIL" --role="$role"
done

# GitHub からの impersonation を許可
export PROJECT_NUMBER="$(gcloud projects describe $PROJECT_ID --format='value(projectNumber)')"
gcloud iam service-accounts add-iam-policy-binding "$SA_EMAIL" \
  --role=roles/iam.workloadIdentityUser \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github/attribute.repository/kibeshohei/cutagent"
```

## 6. Gemini API キーを Secret Manager に登録

```sh
echo -n "YOUR_GEMINI_API_KEY" | gcloud secrets create gemini-api-key --data-file=-

# Cloud Run の実行時 SA から読めるようにする
# (デプロイ後に Cloud Run サービスが使う SA に変更。デフォルトは {PROJECT_NUMBER}-compute@developer.gserviceaccount.com)
export RUNTIME_SA="${PROJECT_NUMBER}-compute@developer.gserviceaccount.com"
gcloud secrets add-iam-policy-binding gemini-api-key \
  --member="serviceAccount:$RUNTIME_SA" \
  --role=roles/secretmanager.secretAccessor
```

> 本番化前に実行時用の専用 SA を作って `RUNTIME_SA` を差し替えるのが望ましい (M3 タスクで検討)。

## 7. GitHub Secrets 登録

GitHub の `Settings → Secrets and variables → Actions` で以下を登録:

| 名前 | 値 |
|---|---|
| `GCP_PROJECT_ID` | `$PROJECT_ID` |
| `GCP_REGION` | `$REGION` |
| `WIF_PROVIDER` | `projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github/providers/github-provider` |
| `WIF_SERVICE_ACCOUNT` | `${SA_EMAIL}` |

## 8. Firestore のアクセス権限

Cloud Run の実行時 SA に Firestore 読み書き権限を付与。

```sh
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${RUNTIME_SA}" \
  --role="roles/datastore.user"
```

## 9. Firebase Hosting 初期化

```sh
# Firebase CLI を入れていない場合
pnpm dlx firebase-tools login

cd web
pnpm dlx firebase-tools init hosting
# - Use existing project → $PROJECT_ID を選択
# - Public directory: dist
# - Single-page app: Yes
# - GitHub Actions: 任意 (今回は deploy.yml で扱う想定)
```

`web/firebase.json` に `/api/*` を Cloud Run へ rewrite する設定:

```json
{
  "hosting": {
    "public": "dist",
    "ignore": ["firebase.json", "**/.*", "**/node_modules/**"],
    "rewrites": [
      { "source": "/api/**", "run": { "serviceId": "cutagent-api", "region": "asia-northeast1" } },
      { "source": "**", "destination": "/index.html" }
    ]
  }
}
```

## 10. 動作確認

1. `main` への PR をマージ → GitHub Actions `deploy.yml` が走る
2. `gcloud run services describe cutagent-api --region $REGION --format='value(status.url)'` で URL 取得
3. `curl <URL>/api/health` で 200 が返れば成功
4. Firebase Hosting の URL でフロントが表示され、`/api/*` の呼び出しが Cloud Run に届くことを確認

---

## トラブルシュート

| 症状 | 原因の候補 |
|---|---|
| WIF で「failed to exchange token」 | Provider の `attribute-mapping` / `attribute-condition` のタイポ、リポジトリ名違い |
| Cloud Run デプロイで「permission denied」 | デプロイ SA に `roles/run.admin` がない、Artifact Registry の `roles/artifactregistry.writer` がない |
| Gemini API で 403 | `gemini-api-key` の Secret Manager 権限を Cloud Run 実行 SA に付与しているか |
| Firestore で「PERMISSION_DENIED」 | 実行 SA に `roles/datastore.user` がない |
| Firebase Hosting で `/api` が 404 | `firebase.json` の `rewrites` の Cloud Run サービス名・リージョンを確認 |

---

## 関連ドキュメント

- 全体アーキテクチャ: [`architecture.md`](architecture.md)
- WIF を選んだ理由: 今後の ADR で追加予定
- デプロイワークフロー: [`../.github/workflows/deploy.yml`](../.github/workflows/deploy.yml)
