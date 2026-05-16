// Package infra は domain で定義したポート（Repository 等）の具象実装を提供します。
//
// この層は DB・外部 API・メッセージングなどの技術詳細を閉じ込め、
// 上位層（usecase / query / registry）からは domain のインターフェース越しにのみ利用されます。
//
// 構成:
//   - 実装ごとにサブパッケージを切る（例: mariadb）。集約ごとに *Repository を置く。
//   - sqlboiler 生成物は <実装>/models に置き、手編集しない（make generate.boilerplate）。
//     生成用スキーマは migrations/schema.sql.boiler（migrations/README.md 参照）。
//
// 例:
//   - user.Repository の MariaDB 実装: internal/infra/mariadb（NewUserRepository）。
//     domain/user.User と models.User の相互変換、Insert / Update の使い分けを担当する。
//
// ルール:
//   - domain の Repository インターフェースを実装し、戻り値は domain の型・domain/errors に揃える。
//   - HTTP ハンドラ・usecase のオーケストレーション・ビジネス判断は持ち込まない。
//   - トランザクションは pkg/sqldb の Client（CurrentTx）経由で ctx から取得する。
//   - 接続先や実装の切り替えは internal/registry で行い、infra パッケージ内に環境分岐を散らさない。
package infra
