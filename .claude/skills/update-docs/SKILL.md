---
name: update-docs
description: セッション中の作業内容・学びをdocs/changelog/とdocs/learnings/に記録する
allowed-tools: Read Write Edit Glob
---

## ドキュメント更新

セッション中にまとまった作業が完了したタイミングや、セッション終了時に以下を更新する。

### 1. docs/changelog（開発プロセスログ）

`docs/changelog/YYYY-MM-DD.md` に追記または新規作成する。

記録する内容（セクションごとに）:
- **AIに依頼したこと**: 調査、比較、提案、レビューなど何をさせたか
- **自分が判断したこと**: AIの提案に対してどう意思決定したか
- **学び**: その過程で得た知見

既にその日のファイルがある場合は、新しいセクションを追記する。

`docs/changelog/README.md` のファイル一覧も更新する。

### 2. docs/learnings（学びメモ）

技術的に新しく学んだことがあれば、該当するファイルに追記するか、新しいトピックのファイルを作成する。

例: `architecture.md`, `go-basics.md`, `ent-orm.md`, `auth-jwt.md` 等

`docs/learnings/README.md` のファイル一覧も更新する。

### 書き方のルール

- 事実だけでなく「なぜその判断をしたか」を残す
- AIが出した提案をそのまま採用したのか、修正したのかを明記する
- 完成後にSDD・バイブコーディングの発表資料にまとめることを意識して書く
