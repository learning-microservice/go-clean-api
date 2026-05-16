package sqldb

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// MariaDB / MySQL の主なエラーコード定義
/*
1062: 重複エントリー（PRIMARY KEY や UNIQUE 制約に違反）
1213: デッドロック（トランザクションが競合。アプリ側でリトライが必要なケース）
1451/1452: 外部キー制約違反（子データがある親を消そうとした、または存在しない親IDを指定した）
*/
const (
	ErrCodeDuplicateKey = 1062 // 重複エラー（ユニーク制約違反）
	ErrCodeForeignKey   = 1452 // 外部キー制約違反（親レコードが存在しないなど）
)

func errorCode(err error) uint16 {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number
	}
	return 0
}

func isDuplicateKeyError(err error) bool {
	return errorCode(err) == ErrCodeDuplicateKey
}

func isForeignKeyError(err error) bool {
	return errorCode(err) == ErrCodeForeignKey
}
