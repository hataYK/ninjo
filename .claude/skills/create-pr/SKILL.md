---
name: create-pr
description: GitHub Flowに従ったPRを作成する
allowed-tools: Bash(git *) Bash(gh *) Read Glob Grep
---

## PR作成手順

### 1. 現在の状態を確認

!`git branch --show-current`
!`git log main..HEAD --oneline`

### 2. PR作成ルール

- PRタイトルは70文字以内、英語
- mainブランチに対してPRを作成
- Squash merge を前提としたPR

### 3. PRテンプレート

```markdown
## Summary
<変更内容を1-3個の箇条書きで>

## Test plan
<テスト方法をチェックリストで>
```

### 4. 手順

1. 未コミットの変更がないか確認
2. リモートにブランチをpush（未pushの場合）
3. `gh pr create` でPR作成
4. PR URLを返す

引数 $ARGUMENTS が指定されている場合はPRタイトルのヒントにする。
