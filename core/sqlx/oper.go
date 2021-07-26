package sqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type dbRow map[string]interface{}

type Dblib struct {
	db *sql.DB
}

const (
	USER_NAME = "root"
	PASS_WORD = "123456"
	HOST      = "localhost"
	PORT      = "3306"
	DATABASE  = "test"
)

// 初始化链接
func NewMysqlConn() (*Dblib, error) {
	dbDSN := BuildDataSourceName(USER_NAME, PASS_WORD, HOST, PORT, DATABASE)
	db, err := NewMySQLInstance(dbDSN, WithMaxIdleConns(20),
		WithMaxOpenConns(100),
		WithConnMaxLifetime(100*time.Second))
	if nil != err {
		panic(err.Error())
		return nil, err
	}

	p := new(Dblib)
	p.db = db
	return p, nil
}

func scanRow(rows *sql.Rows) (dbRow, error) {
	columns, _ := rows.Columns()

	vals := make([]interface{}, len(columns))
	valsPtr := make([]interface{}, len(columns))

	for i := range vals {
		valsPtr[i] = &vals[i]
	}

	err := rows.Scan(valsPtr...)

	if err != nil {
		return nil, err
	}

	r := make(dbRow)

	for i, v := range columns {
		if va, ok := vals[i].([]byte); ok {
			r[v] = string(va)
		} else {
			r[v] = vals[i]
		}
	}

	return r, nil

}

// 获取一行记录
func (d *Dblib) GetOne(sql string, args ...interface{}) (dbRow, error) {
	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	rows.Next()
	result, err := scanRow(rows)
	return result, err
}

// 获取多行记录
func (d *Dblib) GetAll(sql string, args ...interface{}) ([]dbRow, error) {
	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]dbRow, 0)

	for rows.Next() {
		r, err := scanRow(rows)
		if err != nil {
			continue
		}

		result = append(result, r)
	}

	return result, nil

}

// 写入记录
func (d *Dblib) Insert(table string, data dbRow) (int64, error) {
	fields := make([]string, 0)
	vals := make([]interface{}, 0)
	placeHolder := make([]string, 0)

	for f, v := range data {
		fields = append(fields, f)
		vals = append(vals, v)
		placeHolder = append(placeHolder, "?")
	}

	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s) ", table, strings.Join(fields, ","), strings.Join(placeHolder, ","))
	result, err := d.db.Exec(sql, vals...)
	if err != nil {
		return 0, err
	}

	lID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lID, nil
}

// 更新记录
func (d *Dblib) Update(table, condition string, data dbRow, args ...interface{}) (int64, error) {
	params := make([]string, 0)
	vals := make([]interface{}, 0)

	for f, v := range data {
		params = append(params, f+"=?")
		vals = append(vals, v)
	}

	sql := "UPDATE %s SET %s"
	if condition != "" {
		sql += " WHERE %s"
		sql = fmt.Sprintf(sql, table, strings.Join(params, ","), condition)
		vals = append(vals, args...)
	} else {
		sql = fmt.Sprintf(sql, table, strings.Join(params, ","))
	}

	result, err := d.db.Exec(sql, vals...)
	if err != nil {
		return 0, err
	}

	aID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return aID, nil
}

// 删除记录
func (d *Dblib) Delete(table, condition string, args ...interface{}) (int64, error) {
	sql := "DELETE FROM %s "
	if condition != "" {
		sql += "WHERE %s"
		sql = fmt.Sprintf(sql, table, condition)
	} else {
		sql = fmt.Sprintf(sql, table)
	}

	result, err := d.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	aID, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return aID, nil

}
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(obj interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
