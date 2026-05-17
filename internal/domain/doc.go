// Package domain はプロジェクトの中核となるドメインルールを定義します。
//
// 構成:
//   - 集約ごとにサブパッケージを置く（例: user）。各パッケージにエンティティ・値オブジェクト・
//     その集約用の Repository インターフェース（永続化のポート）をまとめる。
//   - ドメイン全体で共有するエラー種別は domain/errors に置く。
//
// 外部実装（DB/HTTP/framework）に依存しない純粋なモデルを維持します。
//
// ルール:
//   - infra の実装詳細や delivery の transport 型に依存しない。
//   - 永続化や外部通信は、各集約の Repository インターフェース経由で抽象化し、実装は infra に置く。
//   - ドメイン概念に関する語彙・制約・不変条件を、対応する集約パッケージに集約する。
//   - エラーは domain/errors の TypeXxx（pkg/errors.Type[FieldError](code, name)）で表現する。HTTP 相当コードはここで定義する。
package domain
