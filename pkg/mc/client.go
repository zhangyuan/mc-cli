package mc

import (
	"database/sql"

	_ "github.com/aliyun/aliyun-odps-go-sdk/sqldriver"
)

func NewDB(dsn string) (*sql.DB, error) {
	return sql.Open("odps", dsn)
}
