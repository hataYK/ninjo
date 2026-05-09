# F0. ユーザー認証

## 概要

メール + パスワードによるユーザー認証。JWT（アクセストークン + リフレッシュトークン）で認証状態を管理する。

## エンドポイント

### POST /api/v1/auth/signup

ユーザー新規登録。

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "display_name": "田中太郎"
}
```

**Response (201):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "display_name": "田中太郎"
}
```

**Cookie:**
- `access_token` (httpOnly, 1h)
- `refresh_token` (httpOnly, 30日)

**エラー:**
- 409: メールアドレスが既に登録されている
- 400: バリデーションエラー

### POST /api/v1/auth/login

ログイン。

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "display_name": "田中太郎"
}
```

**Cookie:**
- `access_token` (httpOnly, 1h)
- `refresh_token` (httpOnly, 30日)

**エラー:**
- 401: メールアドレスまたはパスワードが間違っている

### POST /api/v1/auth/logout

ログアウト。Cookieを削除する。

**Response (204):** No Content

**Cookie:**
- `access_token` を削除
- `refresh_token` を削除

### POST /api/v1/auth/refresh

アクセストークンを更新。

**Cookie (入力):**
- `refresh_token`

**Response (200):**
```json
{
  "message": "token refreshed"
}
```

**Cookie (出力):**
- `access_token` (httpOnly, 新しいトークン, 1h)

**エラー:**
- 401: リフレッシュトークンが無効または期限切れ

## バリデーション

| フィールド | ルール |
|-----------|--------|
| email | 必須、メール形式、255文字以内 |
| password | 必須、8文字以上 |
| display_name | 必須、1〜100文字 |

## セキュリティ

- パスワードは bcrypt でハッシュ化して保存
- トークンは httpOnly Cookie で管理（XSS対策）
- アクセストークン: 1時間で期限切れ
- リフレッシュトークン: 30日で期限切れ
- JWT署名キーは環境変数 `JWT_SECRET` から読み取り

## 認証ミドルウェア

認証が必要なエンドポイント（auth以外すべて）に適用。

1. Cookie から `access_token` を取得
2. JWT を検証（署名 + 有効期限）
3. ペイロードから user_id を取得し、コンテキストにセット
4. トークンがない or 無効 → 401 Unauthorized
