# F6: 可処分時間設定 - 機能仕様書

## 概要

曜日ごとに勉強に使える時間を設定する。
この情報はAIの計画レビューやデイリータスク配分の基礎データとなる。

## ユーザーストーリー

> ユーザーとして、曜日ごとの勉強可能時間を設定し、
> AIに自分の生活リズムを理解してもらいたい。

## 画面

設定画面に曜日ごとの時間入力フォームを表示。

```
可処分時間の設定
──────────────────
日曜  [ 4.0 ] 時間
月曜  [ 2.0 ] 時間
火曜  [ 2.0 ] 時間
水曜  [ 1.5 ] 時間
木曜  [ 2.0 ] 時間
金曜  [ 1.0 ] 時間
土曜  [ 5.0 ] 時間
──────────────────
週合計: 17.5 時間
──────────────────
        [ 保存する ]
```

## API エンドポイント

### GET /api/v1/availability

現在の可処分時間設定を取得する。未設定の曜日は 0h として返す（DBにレコードがなくてもデフォルト値を返す）。

**Response (200):**
```json
{
  "availability": [
    { "day_of_week": 0, "hours": 4.0 },
    { "day_of_week": 1, "hours": 2.0 },
    { "day_of_week": 2, "hours": 2.0 },
    { "day_of_week": 3, "hours": 1.5 },
    { "day_of_week": 4, "hours": 2.0 },
    { "day_of_week": 5, "hours": 1.0 },
    { "day_of_week": 6, "hours": 5.0 }
  ],
  "weekly_total": 17.5
}
```

- `day_of_week`: Go `time.Weekday` 準拠（0=日曜, 1=月曜, ..., 6=土曜）
- `label` はフロントエンド側で day_of_week から変換する（バックエンドは返さない）

**エラー:**
- 401: 未認証

### PUT /api/v1/availability

可処分時間を一括更新する（7曜日分まとめて）。
レコードが存在しない曜日は新規作成、存在する曜日は更新（upsert）。

**Request:**
```json
{
  "availability": [
    { "day_of_week": 0, "hours": 4.0 },
    { "day_of_week": 1, "hours": 2.0 },
    { "day_of_week": 2, "hours": 2.0 },
    { "day_of_week": 3, "hours": 1.5 },
    { "day_of_week": 4, "hours": 2.0 },
    { "day_of_week": 5, "hours": 1.0 },
    { "day_of_week": 6, "hours": 5.0 }
  ]
}
```

**Response (200):**
```json
{
  "availability": [
    { "day_of_week": 0, "hours": 4.0 },
    { "day_of_week": 1, "hours": 2.0 },
    { "day_of_week": 2, "hours": 2.0 },
    { "day_of_week": 3, "hours": 1.5 },
    { "day_of_week": 4, "hours": 2.0 },
    { "day_of_week": 5, "hours": 1.0 },
    { "day_of_week": 6, "hours": 5.0 }
  ],
  "weekly_total": 17.5
}
```

**バリデーション:**
- hours: 0.0〜24.0 の範囲（0.5刻み）
- 7曜日分すべて必須
- day_of_week: 0〜6（Go `time.Weekday` 準拠）
- day_of_week の重複不可

**エラー:**
- 400: バリデーションエラー
- 401: 未認証

## データ設計

- ent スキーマ: `hours` に `Max(24)` を追加
- 0.5刻みバリデーションは usecase 層で実施
- DBにレコードがない曜日はデフォルト 0h（ユーザー作成時にレコードは挿入しない）
- PUT 時に upsert（存在すれば更新、なければ作成）

## 認証

- GET/PUT ともに認証必須
- oapi-codegen の `OperationMiddlewares` で JWT 認証ミドルウェアを適用
- Cookie から access_token を取得し、ユーザー ID をコンテキストに設定

## ビジネスルール

- 0h の曜日はその日にタスクを配分しない
- 更新すると既存の計画のAIレビューには影響しない（次回レビュー時に反映）
