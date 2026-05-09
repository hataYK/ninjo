# CLAUDE.md - Ninjo プロジェクトルール

> このファイルはClaude Codeが毎回自動で参照する。ここに書かれたルールは必ず守ること。

## 仕様書

- 全体の原則: `specs/constitution.md`
- 技術仕様: `specs/technical.md`
- プロダクト仕様: `specs/product.md`
- 機能仕様: `specs/features/` 配下

## アーキテクチャ（4層構成）

```
handler → usecase → domain ← infra
```

- `internal/domain/` — エンティティ・値オブジェクト・リポジトリIF・サービスIF（最内側、依存なし）
- `internal/usecase/` — ビジネスロジック（domainに依存）
- `internal/handler/` — HTTPハンドラ・DTO・ミドルウェア（usecaseに依存）
- `internal/infra/` — DB永続化（ent）・外部API実装（domainのIFを実装）
- ent の生成コードは infra 層でのみ使用

### DI: ファサードパターン

- `infra/datastore.go` — 全リポジトリへのアクセサ（機能追加時はここにメソッド追加）
- `usecase/usecase.go` — 全ユースケースへのアクセサ（機能追加時はここにメソッド追加）
- main.go では `DataStore → Usecase → handler` の順に組み立てるだけ

## Git ルール

- コミット: Conventional Commits（詳細は `/commit` スキル）
- ブランチ: GitHub Flow（詳細は `/create-pr` スキル）
- 新機能: SDD フロー（詳細は `/new-feature` スキル）
- main への直接 push 禁止（初期セットアップ時を除く）

## 開発の原則

- SDD: 仕様が先、コードは後
- 新機能はまず `specs/features/` に仕様書を書く
- OpenAPI スキーマが BE/FE の契約
