---
name: new-feature
description: SDD（仕様駆動開発）に従って新機能の開発を開始する
allowed-tools: Bash(git *) Read Write Edit Glob Grep
---

## SDD 新機能開発フロー

引数: $ARGUMENTS（機能名 or 機能の説明）

### 1. ブランチ作成

mainから新しいfeatureブランチを切る。

```
feature/<機能名（ハイフン区切り）>
```

### 2. 仕様書作成（コードより先）

`specs/features/` 配下に機能仕様書を作成し、最初にコミットする。

```
docs(<scope>): add feature specification
```

仕様書に含めるべき内容:
- 機能概要
- APIエンドポイント定義（メソッド、パス、リクエスト/レスポンス）
- バリデーションルール
- ビジネスルール
- 画面仕様（該当する場合）

### 3. 実装（仕様に従って）

アーキテクチャ（4層構成）に従い、以下の順で実装:

1. `internal/domain/` — モデル・リポジトリIF・サービスIF
2. `internal/usecase/` — ビジネスロジック
3. `internal/infra/` — DB永続化・外部API実装
4. `internal/handler/` — HTTPハンドラ・DTO
5. フロントエンド（該当する場合）

各レイヤーごとに意味のあるコミットを作る。

### 4. 参照すべき仕様書

- 全体原則: `specs/constitution.md`
- 技術仕様: `specs/technical.md`
- プロダクト仕様: `specs/product.md`
- 既存の機能仕様: `specs/features/`
