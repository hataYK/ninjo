# アーキテクチャ・設計の学び

## 2026-05-03: Go DDD + クリーンアーキテクチャのディレクトリ構成

### 3層ではなく4層が標準

最初は「UI / Usecase / Datastore」の3層で考えていたが、調査の結果 **domain を独立した層として分離する4層構成**が Go 界隈の標準だとわかった。

```
handler → usecase → domain ← infra
```

- `domain/`: エンティティ・値オブジェクト・リポジトリIF・サービスIF（最内側、依存なし）
- `usecase/`: ビジネスロジック（domainに依存）
- `handler/`: HTTPハンドラ・DTO（usecaseに依存）
- `infra/`: DB永続化・外部API実装（domainのIFを実装）

### なぜ domain を分離するのか

ドメインモデルは「ビジネスルールそのもの」であり、ユースケースに従属しない。
usecase の中に model/ を置くと、ドメインがユースケースの一部に見えてしまう。
依存関係的にも usecase が domain に依存するのが正しい方向。

### 命名の慣習（Go界隈）

| 層 | 一般的な命名 | 避けたほうがいい命名 |
|---|---|---|
| プレゼンテーション | `handler/`, `delivery/` | `ui/`（フロントと混同） |
| インフラ | `infra/`, `infrastructure/` | `datastore/`（DB以外を含みにくい） |
| ドメイン | `domain/` | `usecase/model/`（従属関係が誤解される） |

### 参考にしたリソース

- [bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch) — Go CleanArchの定番（Star 9k+）
- [Three Dots Labs - wild-workouts](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
