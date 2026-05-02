---
name: commit
description: Conventional Commitsに従ったコミットを作成する
allowed-tools: Bash(git *) Read Glob Grep
---

## コミット作成手順

### 1. 変更内容を確認

!`git status`
!`git diff --stat`

### 2. コミットルール

Conventional Commits フォーマットに従うこと。

```
<type>(<scope>): <英語で簡潔な説明>

<補足（日本語OK、必要な場合のみ）>
```

#### type 一覧

| type | 用途 |
|------|------|
| feat | 新機能 |
| fix | バグ修正 |
| docs | ドキュメント変更 |
| refactor | 機能変更なしのコード改善 |
| test | テスト追加・修正 |
| chore | ビルド・依存関係・雑務 |
| ci | CI/CD設定変更 |
| perf | パフォーマンス改善 |
| style | フォーマット変更（動作に影響なし） |

### 3. ルール

- 1コミット = 1つの論理的変更
- リファクタリングと機能追加は別コミット
- テストと実装は同じコミットでOK
- .env やシークレットを含むファイルは絶対にコミットしない
- 変更内容を分析して適切な type と scope を選ぶ
- サブジェクト行は50文字以内
- 引数 $ARGUMENTS が指定されている場合はそれをコミットメッセージのヒントにする
