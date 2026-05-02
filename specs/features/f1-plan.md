# F1: 計画作成 + AIレビュー - 機能仕様書

## 概要

教材名・総ページ数・目標期限を入力して学習計画を作成する。
AIが可処分時間と過去実績から「1日あたり何ページ必要か」を算出し、
妥当性をレビューしてフィードバックする。

## ユーザーストーリー

> ユーザーとして、勉強したい本と期限を登録したら、
> AIが「そのペースで現実的か」を教えてくれて、安心して始めたい。

## 画面フロー

```
[計画作成画面]
  教材名:     [ AWS SAA対策本          ]
  総ページ数:  [ 300                   ]
  目標期限:    [ 2026-07-15            ]

              [ AIにレビューしてもらう ]
                      ↓
  ┌─────────────────────────────────┐
  │ AIレビュー結果                    │
  │                                 │
  │ 残り74日 × あなたの可処分時間から  │
  │ 1日あたり約5ページ必要です。      │
  │                                 │
  │ 平日2hなら十分いけるペース。      │
  │ 無理なく進められそうですね！       │
  └─────────────────────────────────┘

       [ 計画を作成する ]  [ 修正する ]
```

## API エンドポイント

### POST /api/v1/plans/review

計画の入力内容をAIにレビューしてもらう（まだ保存はしない）。

**Request:**
```json
{
  "title": "AWS SAA対策本",
  "total_pages": 300,
  "target_date": "2026-07-15"
}
```

**Response (200):**
```json
{
  "daily_pages": 4.8,
  "total_days": 74,
  "available_days": 62,
  "review_message": "残り74日、勉強できる日が62日あるので、1日あたり約5ページ必要です。平日2hなら十分いけるペースですね！"
}
```

- `daily_pages`: 勉強可能日で割った1日あたりの必要ページ数
- `total_days`: 今日から目標期限までの日数
- `available_days`: そのうち可処分時間 > 0 の日数
- `review_message`: AIが生成したレビューコメント

**エラー:**
- 400: バリデーションエラー
- 400: 目標期限が今日以前

### POST /api/v1/plans

レビュー確認後、計画を確定して保存する。

**Request:**
```json
{
  "title": "AWS SAA対策本",
  "total_pages": 300,
  "target_date": "2026-07-15"
}
```

**Response (201):**
```json
{
  "id": "uuid",
  "title": "AWS SAA対策本",
  "total_pages": 300,
  "current_page": 0,
  "start_date": "2026-05-02",
  "target_date": "2026-07-15",
  "status": "active",
  "ai_review": "残り74日、1日あたり約5ページ...",
  "progress_rate": 0.0,
  "created_at": "2026-05-02T10:00:00Z"
}
```

### GET /api/v1/plans

計画一覧を取得する。

**Response (200):**
```json
{
  "plans": [
    {
      "id": "uuid",
      "title": "AWS SAA対策本",
      "total_pages": 300,
      "current_page": 45,
      "start_date": "2026-05-02",
      "target_date": "2026-07-15",
      "status": "active",
      "progress_rate": 0.15,
      "created_at": "2026-05-02T10:00:00Z"
    }
  ]
}
```

### GET /api/v1/plans/{plan_id}

計画の詳細を取得する。

**Response (200):**
```json
{
  "id": "uuid",
  "title": "AWS SAA対策本",
  "total_pages": 300,
  "current_page": 45,
  "start_date": "2026-05-02",
  "target_date": "2026-07-15",
  "status": "active",
  "ai_review": "...",
  "progress_rate": 0.15,
  "daily_pages_needed": 4.8,
  "days_remaining": 74,
  "created_at": "2026-05-02T10:00:00Z"
}
```

### DELETE /api/v1/plans/{plan_id}

計画を削除する。紐づくデイリータスク・メモも削除。

**Response (204):** No Content

## AIレビューのロジック

```
基本計算:
  totalDays  = targetDate - today（日数）
  availDays  = そのうち可処分時間 > 0 の日数（曜日設定から算出）
  dailyPages = totalPages / availDays

AIへのプロンプト:
  ユーザーの学習計画をレビューしてください。

  教材: {title}
  総ページ数: {totalPages}
  期限まで: {totalDays}日（うち勉強可能日: {availDays}日）
  1日あたり必要ページ数: {dailyPages}
  ユーザーの可処分時間: {availability}

  {pastPerformance があれば}
  過去の実績: 1時間あたり平均{pagesPerHour}ページ

  短く（3-4文）でレビューしてください。
  現実的ならポジティブに、厳しければ正直に伝えてください。
```

## バリデーション

- title: 1〜200文字
- total_pages: 1〜10000
- target_date: 今日より後の日付
- 同時にactive状態の計画は最大10個

## ビジネスルール

- start_date は計画作成日（当日）を自動設定
- current_page は作成時 0、デイリータスク完了で加算
- progress_rate = current_page / total_pages
- current_page が total_pages に達したら status を completed に自動変更
- 可処分時間が未設定（全曜日0h）でもレビュー可能。その場合はtotal_daysで割る
