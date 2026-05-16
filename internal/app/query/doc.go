// Package query は読み取り系（Query）のアプリケーションロジックを提供します。
//
// この層は CQRS の Query 側を担当し、クライアント要件に最適化された
// Input/Output を返します。業務更新（Command）は扱いません。
//
// 例:
//   - user: internal/app/query/user（ユーザー検索など）。provider/memory 等で読み取り実装を差し替える。
//
// ルール:
//   - registry から app.Interactor としてラップして合成してよい（読み取りでも logging / validation を付与可能）。
//   - 命名は Input/Output を使う（Request/Response は delivery 層で使用）。
//   - query は最短経路の読み取りを担当し、複数ユースケースにまたがるオーケストレーションは app/service を検討する。
//   - transport 型（proto generated）や infra 実装詳細はこの層に持ち込まない。
//   - ドメインの語彙が必要なら domain の集約パッケージ（例: user）から参照してよい。
package query
