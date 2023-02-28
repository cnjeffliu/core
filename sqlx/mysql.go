/*
 * @Author: cnzf1
 * @Date: 2021-06-15 11:58:51
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-05-27 17:37:35
 * @Description: 数据库连接工具类
 */
package sqlx

import (
	"database/sql"
	"fmt"
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

func BuildDataSourceName(uname, passwd, host, port, dbname string) (dsn string) {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", uname, passwd, host, port, dbname)
}
