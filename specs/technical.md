# 技術仕様書 - Ninjo

## 技術スタック

| レイヤー | 技術 |
|---|---|
| Frontend | Next.js 15 (App Router) + TypeScript |
| UI | Chakra UI |
| サーバー状態管理 | TanStack Query |
| クライアント状態管理 | Jotai |
| API クライアント生成 | Orval（OpenAPI → hooks + 型を自動生成） |
| Backend | Go 1.22+ / Echo v4 |
| ORM | ent (by Meta)（スキーマからコード自動生成） |
| API コード生成 (Go) | oapi-codegen（OpenAPI → ServerInterface + 型 + ルーティング） |
| バリデーション | go-playground/validator |
| AI | Anthropic Go SDK (Claude API) |
| DB | PostgreSQL 16 |
| Infra | AWS (Amplify + RDS) |
| 開発環境 | Docker Compose |
| CI/CD | GitHub Actions |

### Go ライブラリ選定理由

| ライブラリ | 理由 |
|---|---|
| **Echo** | 軽量・高速。ミドルウェアが豊富。Go Web FWの定番 |
| **ent** | スキーマ定義からGoコード自動生成。型安全なクエリ。マイグレーションも自動生成。SDDの「仕様→コード生成」思想と合致 |
| **go-playground/validator** | 構造体タグでバリデーション定義。Echoとの統合が容易 |

## システムアーキテクチャ

```
┌─────────────────┐     HTTPS      ┌─────────────────┐
│   Next.js App   │  ◄──────────►  │   Go / Echo     │
│   (Frontend)    │   REST/JSON    │   (Backend)     │
│   port:3000     │                │   port:8080     │
└─────────────────┘                └───────┬──────────┘
                                           │
                                  ┌────────┴────────┐
                                  │                 │
                             ┌────▼────┐      ┌────▼─────┐
                             │  PG     │      │  Claude  │
                             │  :5432  │      │  API     │
                             └─────────┘      └──────────┘
```

## アーキテクチャ: DDD + クリーンアーキテクチャ（4層構成）

```
依存の方向:

  handler → usecase → domain ← infra
  (HTTP)    (business)  (model/IF)  (persistence/external)
```

### 各層の責務

| 層 | 責務 | 依存先 |
|---|---|---|
| **Domain** | エンティティ、値オブジェクト、リポジトリIF、サービスIF。ビジネスルールを表現 | なし（最内側） |
| **Usecase** | ビジネスロジックの実行。ドメインオブジェクトを組み合わせてフローを制御 | Domain のみ |
| **Handler** | HTTPハンドラ、ミドルウェア、ルーティング、リクエスト/レスポンスDTO | Usecase |
| **Infra** | DB永続化（ent）、外部API（Claude）の実装。DomainのリポジトリIFを実装 | Domain |

### ルール

- Domain 層は外部パッケージに依存しない（最内側）
- Infra は Domain が定義したインターフェースを実装する（依存性逆転）
- ent の生成コードは Infra 層でのみ使用。Domain に漏らさない

### DI（依存性注入）: ファサードパターン

DIライブラリは使わず、2つのファサードに依存を集約する:

```
main.go
  └── DataStore（リポジトリのファサード）  ← infra/datastore.go
        └── .User() → UserRepository
        └── .Plan() → PlanRepository
        └── ...
  └── Usecase（ビジネスロジックのファサード） ← usecase/usecase.go
        └── .Auth() → AuthUsecase
        └── .Plan() → PlanUsecase
        └── ...
  └── handler.RegisterRoutes(e, uc)
```

- 機能追加時は `datastore.go` と `usecase.go` にメソッドを足すだけ
- main.go は変更不要

## バックエンド ディレクトリ構成

```
backend/
├── cmd/
│   └── server/
│       └── main.go                    # エントリポイント（DI・起動）
│
├── internal/
│   ├── domain/                        # ===== Domain層 =====
│   │   ├── model/                     # エンティティ・値オブジェクト
│   │   │   ├── user.go
│   │   │   ├── plan.go
│   │   │   ├── daily_task.go
│   │   │   └── availability.go
│   │   ├── repository/                # リポジトリインターフェース
│   │   │   ├── user.go
│   │   │   ├── plan.go
│   │   │   ├── daily_task.go
│   │   │   └── availability.go
│   │   └── service/                   # サービスインターフェース（AI等）
│   │       └── ai_reviewer.go
│   │
│   ├── usecase/                       # ===== Usecase層 =====
│   │   ├── usecase.go                 # ファサード（全ユースケースへのアクセサ）
│   │   ├── auth.go                    # 認証ユースケース
│   │   ├── plan.go                    # 計画作成・レビュー・一覧
│   │   ├── daily_task.go              # タスク生成・完了
│   │   └── availability.go            # 可処分時間更新
│   │
│   ├── handler/                       # ===== Handler層 =====
│   │   ├── router.go                  # ルーティング定義
│   │   ├── validator.go               # リクエストバリデーター
│   │   ├── middleware/                # ミドルウェア
│   │   │   └── auth.go               # JWT認証
│   │   ├── auth.go
│   │   ├── plan.go
│   │   ├── daily_task.go
│   │   ├── availability.go
│   │   └── dto/                       # リクエスト/レスポンスDTO
│   │       ├── auth.go
│   │       ├── plan.go
│   │       ├── daily_task.go
│   │       └── availability.go
│   │
│   ├── infra/                         # ===== Infra層 =====
│   │   ├── datastore.go               # ファサード（全リポジトリへのアクセサ）
│   │   ├── persistence/               # DB永続化（entを使う）
│   │   │   ├── user_repository.go
│   │   │   ├── plan_repository.go
│   │   │   ├── daily_task_repository.go
│   │   │   └── availability_repository.go
│   │   └── external/                  # 外部サービス
│   │       └── claude_client.go       # Claude API クライアント
│   │
│   └── config/
│       └── config.go
│
├── ent/                               # ent スキーマ & 生成コード
│   ├── schema/                        # スキーマ定義（手動編集）
│   │   ├── user.go
│   │   ├── plan.go
│   │   ├── dailytask.go
│   │   └── availability.go
│   ├── generate.go
│   └── ...                            # 自動生成コード（手動編集禁止）
│
├── go.mod
├── go.sum
├── Dockerfile
└── Makefile
```

## フロントエンド ディレクトリ構成

```
frontend/
├── src/
│   ├── app/                       # Next.js App Router
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   ├── (auth)/login/
│   │   ├── (auth)/signup/
│   │   ├── dashboard/             # 今日のタスク一覧
│   │   ├── plans/                 # 計画一覧
│   │   ├── plans/new/             # 計画作成
│   │   ├── plans/[id]/            # 計画詳細
│   │   └── settings/              # 可処分時間設定
│   │
│   ├── components/
│   │   ├── task/
│   │   ├── plan/
│   │   └── layout/
│   │
│   ├── api/generated/             # Orval 自動生成（手動編集禁止）
│   ├── stores/
│   ├── hooks/
│   └── lib/
│       └── fetch.ts               # ネイティブfetchラッパー
│
├── orval.config.ts
└── package.json
```

## データベース設計

### users

| カラム | 型 | 制約 | 説明 |
|---|---|---|---|
| id | UUID | PK | |
| email | VARCHAR(255) | UNIQUE, NOT NULL | |
| password_hash | VARCHAR(255) | NOT NULL | |
| display_name | VARCHAR(100) | NOT NULL | |
| created_at | TIMESTAMP | NOT NULL | |
| updated_at | TIMESTAMP | NOT NULL | |

### availability（可処分時間）

| カラム | 型 | 制約 | 説明 |
|---|---|---|---|
| id | UUID | PK | |
| user_id | UUID | FK → users, NOT NULL | |
| day_of_week | SMALLINT | NOT NULL | 0=日〜6=土（Go time.Weekday準拠） |
| hours | DECIMAL(3,1) | NOT NULL | 勉強可能時間（h） |

UNIQUE(user_id, day_of_week)

### plans（学習計画）

| カラム | 型 | 制約 | 説明 |
|---|---|---|---|
| id | UUID | PK | |
| user_id | UUID | FK → users, NOT NULL | |
| title | VARCHAR(200) | NOT NULL | 教材名 |
| total_pages | INTEGER | NOT NULL | 総ページ数 |
| start_date | DATE | NOT NULL | 開始日 |
| target_date | DATE | NOT NULL | 目標期限 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'active' | active / completed / paused |
| ai_review | TEXT | NULLABLE | AIレビュー結果 |
| created_at | TIMESTAMP | NOT NULL | |
| updated_at | TIMESTAMP | NOT NULL | |

### daily_tasks（デイリータスク）

| カラム | 型 | 制約 | 説明 |
|---|---|---|---|
| id | UUID | PK | |
| plan_id | UUID | FK → plans, NOT NULL | |
| date | DATE | NOT NULL | 対象日 |
| start_page | INTEGER | NOT NULL | 開始ページ |
| end_page | INTEGER | NOT NULL | AIが算出した目標ページ |
| actual_end_page | INTEGER | NULLABLE | 実際に読み終わったページ（完了時に記入） |
| is_completed | BOOLEAN | NOT NULL, DEFAULT false | |
| memo | TEXT | NULLABLE | 学習メモ（完了時に記入） |
| completed_at | TIMESTAMP | NULLABLE | |
| created_at | TIMESTAMP | NOT NULL | |

INDEX(plan_id, date)

### learning_summaries（学習サマリー）※v1.1

| カラム | 型 | 制約 | 説明 |
|---|---|---|---|
| id | UUID | PK | |
| plan_id | UUID | FK → plans, UNIQUE, NOT NULL | |
| content | TEXT | NOT NULL | AI生成サマリー |
| created_at | TIMESTAMP | NOT NULL | |

## 認証方式

- JWT ベース（アクセストークン1h + リフレッシュトークン30日）
- パスワードは bcrypt でハッシュ化
- httpOnly Cookie でトークン管理
- Go: golang-jwt/jwt でトークン生成・検証

## OpenAPI コード生成 (SDD)

SDD（Specification-Driven Development）のフローに基づき、OpenAPI スキーマを起点にコードを自動生成する。

### バックエンド: oapi-codegen

`docs/openapi/openapi.yaml` を Single Source of Truth として、Go のサーバーインターフェース・型・ルーティングを自動生成する。

```bash
cd backend && oapi-codegen -config oapi-codegen.yaml ../docs/openapi/openapi.yaml
```

生成物（`internal/handler/oapi/server.gen.go`）:
- `ServerInterface`: ハンドラが実装すべきインターフェース
- リクエスト/レスポンスの型（SignupRequest, UserResponse 等）
- `RegisterHandlers()`: Echo へのルーティング自動登録
- 埋め込み OpenAPI spec

ハンドラは `Handler` 構造体で `ServerInterface` を実装し、コンパイル時に検証する:
```go
var _ oapi.ServerInterface = (*Handler)(nil)
```

### フロントエンド: Orval

同じ `docs/openapi/openapi.yaml` から TypeScript 型・API クライアントを自動生成:

```bash
cd frontend && npx orval
```

## ent によるコード生成

```bash
# スキーマからGoコードを自動生成
cd backend && go generate ./ent

# マイグレーション生成（差分検知して自動生成）
go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema
```

entのワークフロー:
1. `ent/schema/` にスキーマを定義（フィールド、エッジ、バリデーション）
2. `go generate` でCRUD・クエリビルダー・マイグレーションを自動生成
3. リポジトリ実装で生成されたクライアントを使って型安全にDB操作
