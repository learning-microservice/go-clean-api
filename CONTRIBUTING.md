# Contributing Guide

このドキュメントは、機能追加・修正を行う際の手順と設計上のルールをまとめたものです。  
環境構築や `make` コマンドの一覧は [README.md](README.md) を参照してください。

## 目次

- [レイヤーと依存の向き](#レイヤーと依存の向き)
- [パッケージの切り方](#パッケージの切り方)
- [エラー（domain/errors）](#エラーdomainerrors)
- [機能追加の流れ（チェックリスト）](#機能追加の流れチェックリスト)
- [例: 読み取り API を追加する（Query）](#例-読み取り-api-を追加するquery)
- [例: 更新 API を追加する（Command）](#例-更新-api-を追加するcommand)
- [Usecase / Query の型定義（app.Interactor）](#usecase--query-の型定義appinteractor)
- [テストの方針](#テストの方針)
- [DB スキーマ・sqlboiler](#db-スキーマsqlboiler)
- [コード品質](#コード品質)
- [やらないこと](#やらないこと)

---

## レイヤーと依存の向き

```text
delivery → app (usecase / query) → domain ← infra
                ↑
            registry（配線・Interceptor 合成）
```

| レイヤ | 責務 | 詳細 |
|--------|------|------|
| **delivery** | HTTP・OpenAPI 型・ステータス変換 | [`internal/delivery/doc.go`](internal/delivery/doc.go) |
| **usecase** | 書き込み（Command） | [`internal/app/usecase/doc.go`](internal/app/usecase/doc.go) |
| **query** | 読み取り（Query） | [`internal/app/query/doc.go`](internal/app/query/doc.go) |
| **domain** | エンティティ・Repository ポート | [`internal/domain/doc.go`](internal/domain/doc.go) |
| **infra** | Repository 実装・sqlboiler models | [`internal/infra/doc.go`](internal/infra/doc.go) |
| **registry** | DI・実装の組み立て | [`internal/registry/doc.go`](internal/registry/doc.go) |

**依存の禁止例**

- `domain` → `infra` / `delivery` / `app` に依存しない
- `usecase` / `query` → Echo や sqlboiler の型を持ち込まない
- `delivery` → `domain` / `infra` を直接呼ばない（registry 経由の Interactor のみ）

---

## パッケージの切り方

| 種類 | 命名 | 例 |
|------|------|-----|
| ドメイン集約 | 単数形 | `internal/domain/user` |
| ユースケース | 機能単位（集約と 1:1 不要） | `internal/app/usecase/auth` |
| クエリ | リソース単位 | `internal/app/query/user` |
| HTTP ハンドラ | API バージョン + リソース | `internal/delivery/restapi/v1/auth` |
| DB 実装 | 技術名 | `internal/infra/mariadb` |

- パッケージ名に複数形（`users`）は使わない（Go の慣習）。
- `domain/user` と変数名 `user` が衝突する場合は import エイリアス（`userdom`）や `entity` などで回避する。
- 複数 usecase をまたぐオーケストレーションが必要になったら [`internal/app/service`](internal/app/service/doc.go) を検討する（現状は未使用）。

---

## エラー（domain/errors）

業務エラーは [`internal/domain/errors`](internal/domain/errors/errors.go) の `TypeXxx` を使います（`pkg/errors.Type[FieldError](HTTPコード, "種別名")`）。

| 層 | やること |
|----|----------|
| **usecase / infra** | `TypeXxx.New` / `Wrap` を返す。判定は `TypeXxx.Is(err)` |
| **domain/errors** | 新規種別と **HTTP 相当コード** をここで定義（正の定義元） |
| **delivery / httperror** | `TypeCode()` / `TypeName()` を JSON に載せるだけ。**ステータス用 switch は書かない** |
| **OpenAPI** | 返すステータスを `responses` に追記（契約と code を一致させる） |

`validate.Errors` や `echo.HTTPError` は httperror が 400 等に正規化します（既存経路）。  
テスト例: [`internal/delivery/restapi/httperror/encoder_test.go`](internal/delivery/restapi/httperror/encoder_test.go)。

### Wrap したときの優先順位

`TypeXxx.Wrap(cause, ...)` では **外側（Wrap した側）の Type が優先**されます。

| 観点 | 挙動 |
|------|------|
| HTTP ステータス / JSON の `type` | チェーン先頭の `*Error`（外側）の `TypeCode` / `TypeName` |
| JSON の `error` | 外側の `Message` |
| `details` | 外側 + 内側（cause）をマージ |
| `TypeXxx.Is(err)` | **外側の Type のみ**一致（内側の Type は `Is` では拾わない） |

例: `TypeInvalidCredentials.Wrap(TypeNotFound.New(...), "…")` → クライアントには **401**（内側の 404 は出ない）。ログインで「存在しないユーザ」と「パスワード不一致」を同じ応答にそろえる意図と同じです。  
内側の種別で判定したい場合は **Wrap せず** その Type を返すか、テストでは外側の Type で `Is` してください。

---

## 機能追加の流れ（チェックリスト）

新しい API を足すときは、**契約（OpenAPI）から内側へ** 進めます。

### 共通

- [ ] 1. [`api/openapi/openapi.yaml`](api/openapi/openapi.yaml) に path / schema / エラー応答を追加（新規 `TypeXxx` なら [エラー（domain/errors）](#エラーdomainerrors) も）
- [ ] 2. `make generate` で [`api/openapi/openapi.gen.go`](api/openapi/openapi.gen.go) を再生成
- [ ] 3. 必要なら `migrations/schema.sql`（と `schema.sql.boiler` のビジネスカラム）を更新 → [`migrations/README.md`](migrations/README.md) 参照
- [ ] 4. スキーマ変更時: `make generate.boilerplate`（必要なら `.reset`）
- [ ] 5. レイヤー実装（下記 Query / Command のどちらか）
- [ ] 6. [`internal/registry/registry.go`](internal/registry/registry.go) で配線
- [ ] 7. [`internal/delivery/restapi/router.go`](internal/delivery/restapi/router.go) にルート登録
- [ ] 8. ハンドラ実装（Bind / Validate / Interactor 呼び出し / エラーは `return err`、成功は `c.JSON`）
- [ ] 9. テスト追加（[テストの方針](#テストの方針)）
- [ ] 10. `make lint` / `make test` が通ることを確認

### Command（POST / PUT / DELETE など）か Query（GET）か

| 種別 | 置き場所 | registry での合成（[`wire.go`](internal/registry/wire.go)） |
|------|----------|--------------------------------------------------------|
| **検索・参照** | `internal/app/query/<リソース>/` | `applyStandard`（ログ・バリデーション） |
| **更新・登録・削除** | `internal/app/usecase/<機能>/` | 原則 `applyStandardWithRequiredTx`（+ DB トランザクション）。読み取りのみのユースケースは `applyStandard` |

---

## 例: 読み取り API を追加する（Query）

例として **`GET /v1/users/me`**（ログイン中ユーザのプロフィール取得）を想定します。

### 1. OpenAPI

`paths` に `/v1/users/me` を追加し、レスポンススキーマを `components/schemas` に定義する。

### 2. Query（app 層）

```
internal/app/query/user/
  dto.go          # SearchInput / SearchOutput など（既存を拡張 or 追加）
  service.go      # QueryService 型（app.Interactor[...]）
  get_me.go       # 実装（例: NewGetMeQuery）
  provider/
    memory/       # サンプル用（現状 registry はここを使用）
    mariadb/      # 本番相当にする場合はこちらを追加
```

- Input/Output の命名は **Input / Output**（delivery の Request/Response とは分ける）。
- 読み取り専用。状態変更は `usecase` 側で行う。

### 3. registry

[`internal/registry/registry.go`](internal/registry/registry.go) の `New` 内で配線します。

```go
userGetMe := applyStandard("user-get-me",
    mariadb.NewGetMeQuery(dbClient),
    deps,
)
registry.QuerySet.UserGetMe = userGetMe
```

> **Note:** 現状の `UserSearch` は [`provider/memory`](internal/app/query/user/provider/memory) です。MariaDB 連携に切り替える場合は `registry` の provider 差し替えのみで済むよう、インターフェースは `query/user` に集約します。

### 4. delivery

```
internal/delivery/restapi/v1/user/me.go   # ハンドラ
```

- `api/openapi` の型で Bind / Validate。
- `service.Execute(c.Request().Context(), &Input{...})` を呼ぶ。
- **成功:** `c.JSON` で OpenAPI のレスポンス型を返す。
- **失敗:** `return err` のみ（ハンドラ内で `httperror.Encode` は呼ばない）。
- **エラー JSON:** [`HTTPErrorHandler`](internal/delivery/restapi/engine.go) が [`httperror.Encode`](internal/delivery/restapi/httperror/encoder.go) を呼び、`domain/errors` の `TypeCode` / `TypeName` 等を `ErrorResponse` に載せる（[エラー（domain/errors）](#エラーdomainerrors)）。

### 5. router

```go
v1Group.GET("/users/me", user.Me(reg.QuerySet.UserGetMe))
```

---

## 例: 更新 API を追加する（Command）

例として **`PATCH /v1/users/me`**（プロフィール更新）を想定します。

### 1. domain（必要な場合）

```
internal/domain/user/
  user.go         # 振る舞い・不変条件
  repository.go   # ポート（メソッド追加）
```

- Repository は **集約ごと** に `domain/<集約>/repository.go` に置く。
- エラーは [`internal/domain/errors`](internal/domain/errors) の `TypeXxx` を使う（[エラー（domain/errors）](#エラーdomainerrors)）。

### 2. usecase

```
internal/app/usecase/user/    # 例: auth とは別パッケージでも可
  ports.go        # UpdateMeUsecase, 外部ポート（Hasher 等）
  update_me.go    # NewUpdateMeUsecase, Input/Output, Execute
```

- ビジネスロジックはここ。DB や HTTP の型は持ち込まない。
- 依存は `ports.go` の interface か `domain` の Repository。

### 3. infra

```
internal/infra/mariadb/user_repository.go   # Repository 実装の拡張
```

- sqlboiler の `models` は **手編集しない**（`make generate.boilerplate`）。
- トランザクションは `pkg/sqldb` の `CurrentTx(ctx)` を使う。

### 4. registry

```go
updateMe := applyStandardWithRequiredTx("user-update-me",
    useruc.NewUpdateMeUsecase(userRepo),
    deps,
)
registry.UsecaseSet.UserUpdateMe = updateMe
```

- DB 書き込みを伴うユースケースは `applyStandardWithRequiredTx`（実行順: logging → validation → Tx → コア）。
- 読み取りのみ（例: 現状の Login）は `applyStandard` でよい（[`registry.go`](internal/registry/registry.go) のコメント参照）。
- ネストした別トランザクションが必要なときだけ `applyStandardWithRequiredNewTx` を検討する。

### 5. delivery + router

Query の例（§4 delivery）と同様。Bind / Validate / Execute、成功は `c.JSON`、失敗は `return err`。メソッドとパスを OpenAPI と揃える。

---

## Usecase / Query の型定義（app.Interactor）

ユースケースは [`app.Interactor`](internal/app/interactor.go) を中心にします。

```go
type Interactor[I, O any] interface {
    Execute(ctx context.Context, input I) (output O, err error)
}
```

### 推奨: 型エイリアスで `Interactor` と同義にする

本リポジトリの auth では、ユースケースの公開型を **型エイリアス** で `Interactor` と結び付けています（例: [`internal/app/usecase/auth/ports.go`](internal/app/usecase/auth/ports.go)）。

```go
type LoginUsecase = app.Interactor[*LoginInput, *LoginOutput]
type RegisterUsecase = app.Interactor[*RegisterInput, *RegisterOutput]
```

- registry や delivery の型注釈で「`app.Interactor` をラップしている」と読み取りやすい。
- Query 側も [`QueryService`](internal/app/query/user/service.go) と同様にエイリアスでよい。

### mockgen とテストの扱い

mockgen の `-source` モードは **メソッドを宣言した `interface` 型** を対象にします。**型エイリアス**はインターフェース宣言ではないため、`LoginUsecase` / `RegisterUsecase` 自体の `Mock*` は生成されません。

| 対象 | `go generate`（`ports.go` など） |
|------|----------------------------------|
| Output Port（`TokenGenerator`, `PasswordHasher` 等） | 生成される |
| ユースケース型（上記エイリアスのみ） | 生成されない |
| [`domain/user.Repository`](internal/domain/user/repository.go) など | 別ファイルの directive で生成 |

**ユースケースを差し替えたいテスト**では次のいずれかを使います。

- **delivery:** 手書き fake（`Execute` のみ実装。`var _ auth.RegisterUsecase = (*fakeRegister)(nil)` でエイリアスとの整合を確認できる）
- **usecase:** ポートと `domain.Repository` を mock し、**ユースケースの本実装**で検証する（推奨）

**どうしても mockgen でユースケース型の `Mock*` が要る場合のみ**、その型だけ別名の `interface` で `app.Interactor[...]` を埋め込むファイルに切り出す、という例外パターンもあります（通常は fake で足ります）。

### ポート（Output Port）のモック

`ports.go` 内の `PasswordHasher` など通常の `interface` は、ファイル先頭の directive で生成します。

```go
//go:generate go tool mockgen -source=$GOFILE -package=$GOPACKAGE_test -destination=./mocks/$GOFILE
```

`domain` の Repository も同様です（[`internal/domain/user/repository.go`](internal/domain/user/repository.go)）。

---

## テストの方針

| 対象 | 何を mock するか | 備考 |
|------|------------------|------|
| **usecase** | `domain` の Repository、`ports` の Hasher 等 | ユースケース本体は本実装。table-driven + **testify**（`require` / `assert`）を推奨 |
| **delivery** | usecase の fake、またはポートのみ | HTTP ステータスと JSON 形状（例: [`httperror/encoder_test.go`](internal/delivery/restapi/httperror/encoder_test.go)） |
| **infra** | testcontainers または `TEST_DB_ADDRESS` | [`internal/infra/mariadb/*_test.go`](internal/infra/mariadb/main_test.go)（`-short` 時は DB セットアップをスキップ） |

**実装のたたき台:** [`internal/app/usecase/auth/register_test.go`](internal/app/usecase/auth/register_test.go)・[`login_test.go`](internal/app/usecase/auth/login_test.go)（`package auth_test` + gomock）。  
Cursor 向け詳細は [`.cursor/rules/testing.mdc`](.cursor/rules/testing.mdc)（`*_test.go` 編集時）。

モック再生成:

```sh
go generate ./internal/domain/user/...
go generate ./internal/app/usecase/auth/...
```

```sh
make test    # gotestsum（-race）。MariaDB 統合テストは Docker または TEST_DB_ADDRESS が必要
```

現状テストが少ない場合でも、**新規機能には usecase のユニットテストを 1 本** 付けることを推奨します。

---

## DB スキーマ・sqlboiler

スキーマ変更時の **ローカル作業手順（AI 補助含む）** は [`migrations/README.md`](migrations/README.md) の「sqlboiler 用 DDL の更新手順」と、[`.cursor/rules/sqlboiler-schema-workflow.mdc`](.cursor/rules/sqlboiler-schema-workflow.mdc) を参照（後者は `migrations/**` を扱うとき Cursor が読み込むルール）。

| ファイル | 用途 |
|----------|------|
| `migrations/schema.sql` | 本番・開発 DB（System Versioning あり） |
| `migrations/schema.sql.boiler` | sqlboiler 生成専用（版管理なし） |

- `schema.sql` のビジネスカラムを変えたら、`schema.sql.boiler` も揃える。
- モデル生成: `make generate.boilerplate`（詳細は [migrations/README.md](migrations/README.md)）。
- 版管理テーブルへの Upsert は避け、Insert / Update を明示する（[`internal/infra/mariadb/user_repository.go`](internal/infra/mariadb/user_repository.go) 参照）。

---

## コード品質

```sh
make format   # golangci-lint fmt
make lint     # golangci-lint run
make test
```

- import 順序は [`.golangci.yml`](.golangci.yml) の `gci` / `gofmt` に従う。
- コメントは日本語でよい。公開型・関数には godoc を 1 行以上付ける。
- 外部ライブラリの新規導入は、チーム合意のうえで行う。

---

## やらないこと

- `domain` に Echo / sqlboiler / `openapi` 型を持ち込む
- `delivery` にビジネスルールや SQL を書く
- `httperror` にドメインエラー種別ごとの HTTP ステータス switch を書く（code は `domain/errors` で定義）
- 学習用の過剰な抽象化（イベントソーシング、厳密 CQRS の読み書き DB 分離、DI フレームワークの導入など）
- `internal/infra/mariadb/models/*.go` の手編集
- **本番用の** 秘密情報のコミット（開発用プレースホルダ `your-secret-key-change-in-production` は `config` / Compose の意図的な既定値）

---

## 参考リンク

- [README.md](README.md) — クイックスタート・環境変数
- [migrations/README.md](migrations/README.md) — DDL と sqlboiler
- [api/openapi/openapi.yaml](api/openapi/openapi.yaml) — API 契約
- 各レイヤの `internal/**/doc.go` — パッケージ単位のルール
- [`.cursor/rules/go-clean-architecture.mdc`](.cursor/rules/go-clean-architecture.mdc) — Go / アーキ（`internal` / `cmd` / `pkg`）
- [`.cursor/rules/openapi-contract.mdc`](.cursor/rules/openapi-contract.mdc) — **OpenAPI・生成コード・Echo ルータの整合**（`api/openapi` / `internal/delivery/restapi`）
- [`.cursor/rules/testing.mdc`](.cursor/rules/testing.mdc) — **テストの書き方**（`*_test.go`）
- [`.cursor/rules/sqlboiler-schema-workflow.mdc`](.cursor/rules/sqlboiler-schema-workflow.mdc) — DDL / `schema.sql.boiler`（`migrations/**`）
