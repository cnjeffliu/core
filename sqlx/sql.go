package sqlx

import "fmt"

// BuildDataSourceName return dsn to username:password@tcp(ip:port)/database？charset=utf8mb4
func BuildDataSourceName(uname, passwd, ip string, port int, dbname string) (dsn string) {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", uname, passwd, ip, port, dbname)
}
