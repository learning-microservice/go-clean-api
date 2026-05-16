# go-clean-api

クリーンアーキテクチャをベースにした Go の REST API サンプルです。認証（ログイン・ユーザ登録）を中心に、レイヤー分離・CQRS 風の構成・OpenAPI 契約を学習・検証するためのリポジトリです。

## 前提条件

- Go **1.26.3**（[`.mise.toml`](.mise.toml) 参照。未導入の場合は [mise](https://mise.jdx.dev/) または同等のツールでバージョンを合わせる）
- [Docker](https://www.docker.com/) / Docker Compose（MariaDB・コード生成用 DB）
- [Make](https://www.gnu.org/software/make/)

## クイックスタート

### 1. 依存関係の取得

```sh
make deps
```

### 2. MariaDB の起動とスキーマ適用

```sh
make db-up
```

初回起動時は [`deployments/docker-compose.yml`](deployments/docker-compose.yml) が `migrations/` を `docker-entrypoint-initdb.d` にマウントするため、**`schema.sql` が自動適用**されます（`schema.sql.boiler` は拡張子のため init 対象外）。

既存ボリュームがある場合や手動で流す場合:

```sh
docker exec -i mariadb mariadb -uroot -ppassword localdb < migrations/schema.sql
```

詳細は [`migrations/README.md`](migrations/README.md)。

### 3. API サーバーの起動

```sh
# ホットリロード（開発向け）
make run

# または直接起動
go run ./cmd/app server
```

デフォルトで `http://localhost:8080` で待ち受けます。

### 4. 動作確認（例）

```sh
# ユーザ登録
curl -s -X POST http://localhost:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{
    "name":"yuki kawamura",
    "email":"yuki.kawamura@example.com",
    "password":"password1234"
  }'

# ログイン
curl -s -X POST http://localhost:8080/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
    "email":"yuki.kawamura@example.com",
    "password":"password1234"
  }'
```

## アーキテクチャ

依存の向きは外側から内側へ。**ドメインはフレームワーク・DB に依存しません**。

```text
delivery (HTTP / Echo)
    → app (usecase / query, interceptor)
        → domain (集約・Repository ポート・エラー)
            ← infra (MariaDB / sqlboiler モデル等)
```

| レイヤ | 役割 | 主なパス |
|--------|------|----------|
| **delivery** | HTTP ルーティング・リクエスト/レスポンス | `internal/delivery/restapi` |
| **app** | ユースケース（Command）・クエリ（Query） | `internal/app/usecase`, `internal/app/query` |
| **domain** | エンティティ・値オブジェクト・Repository インターフェース | `internal/domain/user` など |
| **infra** | DB 実装・生成モデル | `internal/infra/mariadb` |
| **registry** | DI・Interceptor 合成・実装の組み立て | `internal/registry` |
| **pkg** | 横断ユーティリティ（ログ・JWT・DB クライアント等） | `pkg/` |

- **Command / Query:** 更新系は `usecase`（例: `usecase/auth`）、読み取り系は `query`（例: `query/user`）。
- **集約:** ドメインは集約ごとのサブパッケージ（例: `user` に `User` / `ID` / `Repository`）。各パッケージの `doc.go` に方針を記載。
- **横断関心:** バリデーション・ログは `internal/app/interceptor`、DB トランザクション等は `internal/registry/interceptor`。

## ディレクトリ構成（抜粋）

```text
cmd/app/                 CLI エントリ（server / workflow）
config/                  設定（フラグ・環境変数）
deployments/             ローカル用 Docker Compose 等
api/openapi/             OpenAPI 定義・oapi-codegen 生成物
internal/
  delivery/restapi/      HTTP ハンドラ・ルータ
  app/usecase/           書き込み系ユースケース
  app/query/             読み取り系クエリ
  app/interceptor/       ロギング・バリデーション等
  domain/                ドメイン（集約・errors）
  infra/mariadb/         Repository 実装・sqlboiler models
  registry/              依存性注入
migrations/              DDL（本番用 schema.sql / 生成用 schema.sql.boiler）
pkg/                     共有ライブラリ
```

## API 仕様

- 契約定義: [`api/openapi/openapi.yaml`](api/openapi/openapi.yaml)
- 生成型: [`api/openapi/openapi.gen.go`](api/openapi/openapi.gen.go)（`go generate ./api/openapi/...`）

| メソッド | パス | 概要 |
|----------|------|------|
| POST | `/v1/auth/register` | ユーザ登録 |
| POST | `/v1/auth/login` | ログイン（JWT トークン返却） |

## 設定（環境変数）

CLI フラグと環境変数の両方で指定できます（未指定時はデフォルト値）。

> **本番向け:** `JWT_SECRET` は **必ず環境変数やシークレット管理で上書き**してください。リポジトリに含まれるデフォルト値（`your-secret-key-change-in-production`）は **ローカル開発専用** です。

| 変数 | 説明 | デフォルト例 |
|------|------|----------------|
| `DB_ADDRESS` | MySQL DSN | `root:password@tcp(localhost:3306)/localdb?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4` |
| `HTTP_PORT` | 待ち受けポート | `8080` |
| `JWT_SECRET` | JWT 署名鍵 | `your-secret-key-change-in-production`（開発用。本番では変更必須） |
| `JWT_TOKEN_EXPIRY` | トークンの有効期間（time.Duration 文字列。絶対時刻ではない） | `24h` |
| `LOG_LEVEL` | ログレベル | `info` |

その他は `go run ./cmd/app server --help` を参照。

## 開発コマンド

```sh
make help          # 一覧
make lint          # golangci-lint
make format        # フォーマット
make test          # ユニットテスト（gotestsum）
make generate      # go generate（OpenAPI 型など）
make generate.boilerplate        # sqlboiler モデル生成
make generate.boilerplate.reset  # 生成用 DB を作り直してから生成
```

## データベース・コード生成

- **本番・開発 DDL:** [`migrations/schema.sql`](migrations/schema.sql)（System Versioning あり）
- **sqlboiler 用 DDL:** [`migrations/schema.sql.boiler`](migrations/schema.sql.boiler)（版管理なし・論理テーブルと同型）

sqlboiler は版管理付きテーブルをそのまま introspect しないため、**生成専用スキーマ**を別管理しています。手順は [`migrations/README.md`](migrations/README.md) を参照。

## 主な利用ライブラリ

| ライブラリ | 用途 |
|------------|------|
| [Echo v5](https://github.com/labstack/echo) | HTTP サーバ |
| [sqlboiler](https://github.com/aarondl/sqlboiler) | DB モデル生成・クエリ |
| [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) | OpenAPI → Go 型 |
| [go-playground/validator](https://github.com/go-playground/validator) | リクエスト検証 |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT |
| [golangci-lint](https://golangci-lint.run/) | 静的解析・フォーマット |
| [Air](https://github.com/air-verse/air) | 開発時ホットリロード |

実行時依存の一覧は [`go.mod`](go.mod) を参照。

## Docker Compose

定義は [`deployments/docker-compose.yml`](deployments/docker-compose.yml)。日常操作は `make` 経由を推奨。

```sh
# MariaDB のみ
make db-up

# API + MariaDB（ビルドして起動）
make compose-up
```

## 関連ドキュメント

- [`docs/register-flow-5min.md`](docs/register-flow-5min.md) — **5 分で追う Register の道順**（MVC 経験者向け・初見用）
- [`CONTRIBUTING.md`](CONTRIBUTING.md) — 機能追加の手順・レイヤールール・テスト方針
- [`.cursor/rules/go-clean-architecture.mdc`](.cursor/rules/go-clean-architecture.mdc) — Go / アーキ（`internal` / `cmd` / `pkg`）
- [`.cursor/rules/openapi-contract.mdc`](.cursor/rules/openapi-contract.mdc) — OpenAPI とルータ・生成コード
- [`.cursor/rules/testing.mdc`](.cursor/rules/testing.mdc) — テストの書き方
- [`.cursor/rules/sqlboiler-schema-workflow.mdc`](.cursor/rules/sqlboiler-schema-workflow.mdc) — DDL / sqlboiler
- [`migrations/README.md`](migrations/README.md) — DDL と sqlboiler 生成の運用
- 各パッケージの `doc.go` — レイヤー・集約ごとのルール（例: `internal/domain/doc.go`）

## ライセンス

（未設定の場合はプロジェクト方針に従って追記してください。）
