/*
* 数据库连接工具类
* Jeff.Liu <zhifeng172@163.com> 2021.06.15
**/
package sqlx

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const driverName = "mysql"

type MySQLOption func(*sql.DB)

func NewMySQLInstance(dataSourceName string, opts ...MySQLOption) (db *sql.DB, err error) {
	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(60000000000)

	for _, opt := range opts {
		opt(db)
	}

	//验证连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func WithConnMaxLifetime(d time.Duration) MySQLOption {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(d)
	}
}

func WithMaxIdleConns(n int) MySQLOption {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(n)
	}
}

func WithMaxOpenConns(n int) MySQLOption {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(n)
	}
}
