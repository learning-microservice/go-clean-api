// Package registry は依存性注入（DI）とコンポーネント組み立てを担当します。
//
// この層は Registry を構築し、ドメインのポート（例: user.Repository）に対する
// infra の実装（例: mariadb）を束ね、usecase / query を配線します。
//
// 例:
//   - 認証まわり: internal/app/usecase/auth（Login / Register）に user.Repository・JWT・bcrypt 等を注入。
//   - 読み取り: internal/app/query/user にメモリ等の Read モデル実装を注入。
//
// 横断関心の適用（wire.go）:
//   - applyStandard: logging → validation（Query や Tx 不要なユースケース）
//   - applyStandardWithRequiredTx: 上記 + RequiredTx（DB 書き込みを伴うユースケース）
//   - applyStandardWithRequiredNewTx: 上記 + RequiredNewTx（常に新規トランザクション）
//   - internal/app/interceptor … logging / validation
//   - internal/registry/interceptor … DB トランザクション（sqldb）
//
// ルール:
//   - ビジネスロジックは持たない（組み立て専用）。
//   - 環境差分（dev/prod）の実装切り替えポイントは registry に集約する。
package registry
