# F8: ルーム / WiFi連携 - 機能仕様書

> **スコープ: v1.1**
> MVPではデフォルト背景のみ。本機能はネイティブアプリ化以降に本格対応。

## 概要

WiFi SSIDと場所を紐付け、アバターの背景（ルーム）が自動で切り替わる。
自分がいる場所とアバターの居場所が一致することで、デジタルの分身が「今ここにいる」感覚を作る。

## ユーザーストーリー

> ユーザーとして、自分がいる場所に応じてアバターの背景が変わり、
> デジタルの分身が自分と同じ場所にいる感覚を得たい。

## 画面

### 場所登録（設定画面）

```
場所の登録
──────────────────
┌─────────────────────────┐
│ 🏠 自宅                  │
│   WiFi: MyHome-5G        │
│   ルーム: 自分の部屋      │
└─────────────────────────┘
┌─────────────────────────┐
│ 🔬 研究室                │
│   WiFi: Univ-Lab-3F      │
│   ルーム: 研究室          │
└─────────────────────────┘
┌─────────────────────────┐
│ ☕ カフェ                 │
│   WiFi: なし（手動切替）  │
│   ルーム: カフェ          │
└─────────────────────────┘

       [ + 場所を追加 ]
```

## API エンドポイント

### POST /api/v1/locations

場所を登録する。

**Request:**
```json
{
  "name": "自宅",
  "ssid": "MyHome-5G",
  "room_theme": "bedroom"
}
```

- `ssid` は optional（手動切替のみの場所も登録可能）

**Response (201):** 作成された場所オブジェクト

### GET /api/v1/locations

登録済みの場所一覧を取得する。

### PUT /api/v1/locations/{location_id}

場所情報を更新する。

### DELETE /api/v1/locations/{location_id}

場所を削除する。

### PUT /api/v1/avatar/location

現在の場所を設定する（手動 or 自動検出）。

**Request:**
```json
{
  "location_id": "uuid"
}
```

## ルームテーマ

ドット絵のルーム背景を事前定義:

| テーマID | 名前 | 説明 |
|---|---|---|
| bedroom | 自分の部屋 | デフォルト。デスクとベッドのある部屋 |
| library | 図書館 | 本棚に囲まれた静かな空間 |
| cafe | カフェ | コーヒーと窓のある明るい空間 |
| lab | 研究室 | モニターと資料が並ぶ作業場 |
| default | デフォルト | シンプルな背景 |

## 技術的制約

- **ブラウザ**: Network Information API では WiFi SSID を取得できない
- **PWA**: 同様にSSID取得は不可
- **ネイティブアプリ**: iOS/Android では WiFi 情報取得が可能（要パーミッション）
- **MVP対応**: デフォルト背景（"bedroom"）を固定で表示。場所登録・切替UIは実装しない

## データ設計

### locations テーブル

| カラム | 型 | 説明 |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → users |
| name | string | 場所名 |
| ssid | string (nullable) | WiFi SSID |
| room_theme | string | ルームテーマID |
| created_at | timestamp | 作成日時 |

- users テーブルに `current_location_id`（nullable FK → locations）を追加

## ビジネスルール

- 場所は最大10個まで登録可能
- 同一SSIDの重複登録は不可
- 場所を削除すると current_location がデフォルトに戻る
- ルーム背景の切替はフロントエンドで即時反映（API呼び出し後）
