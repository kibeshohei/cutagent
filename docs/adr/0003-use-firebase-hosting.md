# ADR-0003: SPA 配信に Firebase Hosting を採用

- **Status**: Accepted
- **Date**: 2026-05-29

## 文脈 (Context)

React SPA を CDN 付きで配信する場所が必要。ハッカソンが Google Cloud 主催のため GCP エコシステム内で完結させたい。`/api/*` は Cloud Run 上の Go バックエンドに転送する必要がある。

## 決定 (Decision)

**Firebase Hosting** を React SPA の配信先として採用する。

## 根拠 (Rationale)

- CDN・HTTPS が標準で付く (追加設定なし)
- `firebase.json` の `rewrites` で `/api/*` を Cloud Run サービスに転送可能 → SPA と API が同一オリジンで動作 (CORS 不要)
- 無料枠が個人プロジェクト規模に十分
- GCP エコシステム内なので、ハッカソンの「Google Cloud 技術を使う」要件と一貫
- GitHub Actions の `firebase-deploy-action` で main push 時のデプロイが簡潔

## 検討した代替案 (Alternatives)

| 候補 | 採用しなかった理由 |
|---|---|
| Cloud Storage + Cloud CDN | 柔軟だが構成が複雑 (バケット / Load Balancer / SSL 証明書を別々に管理) |
| Vercel | 別ベンダー。ハッカソンが GCP 系なので統一感を欠く |
| Cloud Run で SPA も配信 | api と密結合になり、SPA だけのキャッシュ・CDN 設計ができない |
| GitHub Pages | カスタムドメインや `/api` rewrites がやりにくい |

## トレードオフ・残課題 (Consequences)

- **得たもの**: CDN / HTTPS / 同一オリジン (CORS 不要) / 安いデプロイ
- **諦めたもの**: Firebase という別概念を入れる学習コスト (1 サービスだけなので低い)
- **将来見直す条件**: マルチリージョンや細かい CDN チューニングが必要になった場合 (その時は Cloud CDN へ)
