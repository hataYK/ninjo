# F2: デイリータスク + F3: タスク完了・メモ - 機能仕様書

## 概要

全計画を横断して「今日やること」を1画面に表示する。
AIが計画の進捗と可処分時間から今日分のタスクを自動提案し、
ユーザーはそのまま使うか調整して学習を進める。
タスク完了時にはメモを記入でき、メモはF9（スキル抽出）の素材になる。

## ユーザーストーリー

> ユーザーとして、アプリを開いたら「今日は何をすればいいか」が
> すぐ分かって、終わったらチェックを入れてメモを書きたい。
> メモがアバターのスキルになるのが楽しみ。

## 画面フロー

```
[ホーム画面]
  ┌──────────────────────────────────┐
  │ 今日のタスク（2026-05-11）        │
  │                                  │
  │ 📖 AWS SAA対策本                  │
  │ □ p.46 〜 p.50（5ページ）         │
  │                                  │
  │ 📖 Go言語入門                     │
  │ □ p.21 〜 p.30（10ページ）        │
  │                                  │
  │ [ タスクを生成する ]  ← 未生成時   │
  └──────────────────────────────────┘
                  ↓ チェック
  ┌──────────────────────────────────┐
  │ タスク完了                        │
  │                                  │
  │ 📖 AWS SAA対策本  p.46〜p.50      │
  │                                  │
  │ 実際に読んだページ:               │
  │ [ p.52 まで ]                    │
  │                                  │
  │ 学習メモ（任意）:                 │
  │ ┌──────────────────────────────┐ │
  │ │ VPCのサブネット設計を学んだ。  │ │
  │ │ パブリック/プライベートの      │ │
  │ │ 使い分けが重要。              │ │
  │ └──────────────────────────────┘ │
  │                                  │
  │ [ 完了する ]                      │
  └──────────────────────────────────┘
```

## API エンドポイント

### GET /api/v1/daily-tasks?date={YYYY-MM-DD}

指定日のデイリータスクを全計画横断で取得する。

**Query Parameters:**
- `date` (required): 対象日（YYYY-MM-DD形式）

**Response (200):**
```json
{
  "date": "2026-05-11",
  "tasks": [
    {
      "id": "uuid",
      "plan_id": "uuid",
      "plan_title": "AWS SAA対策本",
      "date": "2026-05-11",
      "start_page": 46,
      "end_page": 50,
      "actual_end_page": null,
      "is_completed": false,
      "memo": null,
      "completed_at": null,
      "created_at": "2026-05-11T07:00:00Z"
    }
  ],
  "summary": {
    "total": 2,
    "completed": 0,
    "total_pages": 15
  }
}
```

**エラー:**
- 400: date が不正な形式
- 401: 未認証

### POST /api/v1/daily-tasks/generate

指定日のデイリータスクをAIが自動生成する。
active な全計画に対してタスクを作成する。

**Request:**
```json
{
  "date": "2026-05-11"
}
```

**Response (201):**
```json
{
  "date": "2026-05-11",
  "tasks": [
    {
      "id": "uuid",
      "plan_id": "uuid",
      "plan_title": "AWS SAA対策本",
      "date": "2026-05-11",
      "start_page": 46,
      "end_page": 50,
      "actual_end_page": null,
      "is_completed": false,
      "memo": null,
      "completed_at": null,
      "created_at": "2026-05-11T07:00:00Z"
    }
  ],
  "summary": {
    "total": 2,
    "completed": 0,
    "total_pages": 15
  }
}
```

**エラー:**
- 400: date が不正な形式、または過去の日付
- 401: 未認証
- 409: 指定日のタスクが既に生成済み

### PUT /api/v1/daily-tasks/{task_id}

タスクのページ範囲を手動で調整する（未完了時のみ）。

**Request:**
```json
{
  "start_page": 46,
  "end_page": 55
}
```

**Response (200):**
```json
{
  "id": "uuid",
  "plan_id": "uuid",
  "plan_title": "AWS SAA対策本",
  "date": "2026-05-11",
  "start_page": 46,
  "end_page": 55,
  "actual_end_page": null,
  "is_completed": false,
  "memo": null,
  "completed_at": null,
  "created_at": "2026-05-11T07:00:00Z"
}
```

**エラー:**
- 400: バリデーションエラー、または完了済みタスク
- 401: 未認証
- 404: タスクが見つからない

### PATCH /api/v1/daily-tasks/{task_id}/complete

タスクを完了にする。実際に読んだページ数とメモを記録する。

**Request:**
```json
{
  "actual_end_page": 52,
  "memo": "VPCのサブネット設計を学んだ。パブリック/プライベートの使い分けが重要。"
}
```

**Response (200):**
```json
{
  "id": "uuid",
  "plan_id": "uuid",
  "plan_title": "AWS SAA対策本",
  "date": "2026-05-11",
  "start_page": 46,
  "end_page": 50,
  "actual_end_page": 52,
  "is_completed": true,
  "memo": "VPCのサブネット設計を学んだ。...",
  "completed_at": "2026-05-11T15:30:00Z",
  "created_at": "2026-05-11T07:00:00Z"
}
```

**エラー:**
- 400: バリデーションエラー、または既に完了済み
- 401: 未認証
- 404: タスクが見つからない

## タスク自動生成のロジック

```
前提:
  - ユーザーの active な計画すべてが対象
  - 対象日にその計画のタスクがまだ存在しない場合のみ生成

各計画について:
  currentPage = その計画の最新 actual_end_page（なければ 0）
  remainPages = totalPages - currentPage
  remainDays  = targetDate - date + 1（対象日を含む残り日数）

  もし remainPages <= 0:
    → タスク生成しない（計画完了済み）

  もし remainDays <= 0:
    → 期限超過。remainPages をそのまま today の分として割り当て

  todayAvailHours = 対象日の曜日に対応する可処分時間
  weekTotalHours  = 週の可処分時間合計

  もし weekTotalHours == 0:
    → 均等割り: todayPages = ceil(remainPages / remainDays)
  そうでなければ:
    → 比例配分: todayPages = ceil(remainPages * (todayAvailHours / 残り日の可処分時間合計))

  startPage = currentPage + 1
  endPage   = min(startPage + todayPages - 1, totalPages)
```

## データ設計

### daily_tasks テーブル（既存）

| カラム | 型 | 制約 | 説明 |
|---|---|---|---|
| id | UUID | PK | |
| plan_id | UUID | FK → plans, NOT NULL | |
| date | DATE | NOT NULL | 対象日 |
| start_page | INTEGER | NOT NULL, > 0 | 開始ページ |
| end_page | INTEGER | NOT NULL, > 0 | 目標ページ |
| actual_end_page | INTEGER | NULLABLE, > 0 | 実際に読み終わったページ |
| is_completed | BOOLEAN | NOT NULL, DEFAULT false | |
| memo | TEXT | NULLABLE | 学習メモ |
| completed_at | TIMESTAMP | NULLABLE | |
| created_at | TIMESTAMP | NOT NULL | |

INDEX(plan_id, date)

> ent スキーマ `backend/ent/schema/dailytask.go` で定義済み。変更不要。

## バリデーション

### タスク生成（POST /daily-tasks/generate）
- date: 必須、YYYY-MM-DD形式、今日以降の日付
- 同一日に同一計画のタスクは1つのみ（重複生成禁止）

### タスク更新（PUT /daily-tasks/{task_id}）
- start_page: 1 以上、end_page 以下
- end_page: start_page 以上、plan.total_pages 以下
- 完了済みタスクは更新不可

### タスク完了（PATCH /daily-tasks/{task_id}/complete）
- actual_end_page: 必須、1 以上、plan.total_pages 以下
- memo: 任意、最大5000文字
- 既に完了済みのタスクは再完了不可

## ビジネスルール

1. **タスクは計画単位で生成**: 1つの計画につき1日1タスク
2. **既存タスクがある日は再生成しない**: 409 Conflict を返す
3. **完了時に actual_end_page が total_pages に達した場合**: plan.status を `completed` に自動変更
4. **進捗の導出**: plan の current_page は daily_tasks の最新 actual_end_page から導出（F1仕様と同じ）
5. **メモの活用**: メモは F9（スキル抽出）の入力データ。タスク完了時にF9が有効なら自動でスキル抽出をトリガー（v1.0ではF9と統合後に実装）
6. **タスク削除**: 個別タスクの削除APIは提供しない。計画削除時にカスケード削除
7. **可処分時間0の日**: タスク生成は可能（ユーザーが明示的に生成を要求した場合）。ページ配分は均等割りにフォールバック
8. **日付の扱い**: サーバーはUTCで管理。クライアントがローカル日付を送信する
