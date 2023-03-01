/*
 * @Author: cnzf1
 * @Date: 2021-06-15 11:58:51
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-01 17:25:15
 * @Description: 数据库连接工具类
 */
package sqlx

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLOption func(*sql.DB)

// NewMySQLInstance create a db connection
func NewMySQLInstance(dsn string, opts ...MySQLOption) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(60 * time.Second)
	for _, opt := range opts {
		opt(db)
	}

	//验证连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

// WithConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
func WithConnMaxLifetime(d time.Duration) MySQLOption {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(d)
	}
}

// WithMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func WithMaxIdleConns(n int) MySQLOption {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(n)
	}
}

// WithMaxOpenConns sets the maximum number of open connections to the database.
func WithMaxOpenConns(n int) MySQLOption {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(n)
	}
}
