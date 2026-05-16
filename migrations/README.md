# migrations

チームおよび Cursor の AI エージェントが従う **sqlboiler 用スキーマ変更フローの要約**は、リポジトリの [`.cursor/rules/sqlboiler-schema-workflow.mdc`](../.cursor/rules/sqlboiler-schema-workflow.mdc) にも定義しています（`migrations/**` を扱うときに Cursor が読み込みます）。

## ファイルの役割

| ファイル | 用途 |
|----------|------|
| `schema.sql` | **本番・開発 DB**（System Versioning あり） |
| `schema.sql.boiler` | **sqlboiler 生成専用**（版管理なし・論理テーブルと同型） |

拡張子が `.sql` で終わらないため、`docker compose` の MariaDB init（`docker-entrypoint-initdb.d`）では **自動実行されません**。本番用 compose で `migrations/` を丸ごとマウントしても、init で流れるのは基本的に `schema.sql` のみです。

sqlboiler は `information_schema` からスキーマを読むため、版管理付きテーブル（`row_start` / `row_end`、複合 PK など）をそのまま introspect するとモデルが使いづらくなります。  
そのため **生成用 DB だけ** `schema.sql.boiler` を流し、Go モデルは通常テーブル前提で生成します。

## モデル生成

```sh
make generate.boilerplate
```

- `docker compose -f deployments/docker-compose.yml --profile boiler` で `mariadb-boiler` を起動し、`schema.sql.boiler` を init として適用（`make generate.boilerplate` が実行）
- `sqlboiler` で `internal/infra/mariadb/models` を生成

### スキーマを変えたあと生成がおかしいとき

MariaDB の init スクリプトは **データディレクトリが空のときだけ** 実行されます。

```sh
make generate.boilerplate.reset
```

（内部で `docker compose -f deployments/docker-compose.yml --profile boiler down -v` のあと `make generate.boilerplate` を実行します。）

## `schema.sql` を変更したとき

`schema.sql.boiler` の **ビジネスカラム**（`id`, `name`, `email`, `password_hash`, `created_at`, インデックス）を `schema.sql` と揃えてください。版管理まわりの定義だけが違います。

### sqlboiler 用 DDL の更新手順（ローカル・AI 補助）

以下は **各メンバーのローカル環境** で行う作業の例です。CI や共有サーバーでの自動実行を想定したものではありません。

1. メンバーが `schema.sql` を修正する。
2. メンバーが AI に対し、修正後の `schema.sql` を根拠に **sqlboiler 向け DDL**（`schema.sql.boiler` と同種）の作成・更新を依頼する。
3. メンバーが AI の出力（大枠は Versioning 除去・単一 PK・ユニーク等の変換）を **Diff で確認**する。
4. 修正が必要な場合はメンバーが直し、**合意した内容を `migrations/schema.sql.boiler` に反映する**（AI にファイル更新を任せる場合も、最終的な中身はこのパスに置く）。
5. ローカルで `make generate.boilerplate` を実行し、**`internal/infra/mariadb/models` の Diff** を確認して問題なければ作業を終える（必要に応じてコミット）。

**手で編集するのは常に `migrations/schema.sql.boiler`** です（生成用 DB への流し込みは Compose / Makefile が担当）。
