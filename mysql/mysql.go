package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" //mysql dirver
)

var dbIns sync.Map

//Dao dao
type Dao struct {
	db *sql.DB
}

type Tx struct {
	t *sql.Tx
}

//Config config
type Config struct {
	NodeName        string `json:"nodeName"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Database        string `json:"database"`
	UserName        string `json:"userName"`
	Password        string `json:"password"`
	Charset         string `json:"charset"`
	MaxOpen         int    `json:"maxOpen"`
	MaxIdle         int    `json:"maxIdle"`
	ConnMaxLifeTime int    `json:"connMaxLifeTime"`
}

//GetDB getDB
func GetDB(config *Config) (*Dao, error) {
	dao := new(Dao)
	if db, ok := dbIns.Load(config.NodeName); ok {
		dao.db = db.(*sql.DB)
		return dao, nil
	}
	db, err := NewDB(config)
	if err != nil {
		return nil, err
	}
	dbIns.Store(config.NodeName, db)
	dao.db = db
	return dao, nil
}

//NewDB newdb
func NewDB(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", config.UserName, config.Password, config.Host, config.Port, config.Database, config.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if config.MaxOpen > 0 {
		db.SetMaxOpenConns(config.MaxOpen)
	}
	if config.MaxIdle > 0 {
		db.SetMaxIdleConns(config.MaxIdle)
	}
	if config.ConnMaxLifeTime > 0 {
		db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
	}
	return db, nil
}

//Select select
func (dao *Dao) Select(sqlStr string, args ...interface{}) ([]map[string]interface{}, error) {
	var stmt *sql.Stmt
	var err error
	var rows *sql.Rows
	stmt, err = dao.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err = stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	var result []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			continue
		}
		record := make(map[string]interface{})
		for i, val := range values {
			switch val.(type) {
			case nil:
				record[columns[i]] = nil
			case bool:
				record[columns[i]] = bool(val.(bool))
			case byte:
				record[columns[i]] = byte(val.(byte))
			case int8:
				record[columns[i]] = int8(val.(int8))
			case int16:
				record[columns[i]] = int16(val.(int16))
			case int32:
				record[columns[i]] = int32(val.(int32))
			case int:
				record[columns[i]] = int(val.(int))
			case int64:
				record[columns[i]] = int64(val.(int64))
			case float32:
				record[columns[i]] = float32(val.(float32))
			case float64:
				record[columns[i]] = float64(val.(float64))
			case []byte:
				record[columns[i]] = string(val.([]byte))
			default:
				record[columns[i]] = string(val.([]byte))
			}
		}
		result = append(result, record)
	}
	return result, nil
}

//Delete delete
func (dao *Dao) Delete(sqlStr string, args ...interface{}) (int64, error) {
	var stmt *sql.Stmt
	var err error
	var result sql.Result
	stmt, err = dao.db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err = stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

//Update update
func (dao *Dao) Update(sqlStr string, args ...interface{}) (int64, error) {
	var stmt *sql.Stmt
	var err error
	var result sql.Result
	stmt, err = dao.db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err = stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

//Insert insert
func (dao *Dao) Insert(sqlStr string, args ...interface{}) (int64, error) {
	var stmt *sql.Stmt
	var err error
	var result sql.Result
	stmt, err = dao.db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err = stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	lastInsertID, _ := result.LastInsertId()
	return lastInsertID, nil
}

//Select select
func (tx *Tx) Select(sqlStr string, args ...interface{}) ([]map[string]interface{}, error) {
	var stmt *sql.Stmt
	var err error
	var rows *sql.Rows
	stmt, err = tx.t.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err = stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	var result []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			continue
		}
		record := make(map[string]interface{})
		for i, val := range values {
			switch val.(type) {
			case nil:
				record[columns[i]] = nil
			case bool:
				record[columns[i]] = bool(val.(bool))
			case byte:
				record[columns[i]] = byte(val.(byte))
			case int8:
				record[columns[i]] = int8(val.(int8))
			case int16:
				record[columns[i]] = int16(val.(int16))
			case int32:
				record[columns[i]] = int32(val.(int32))
			case int:
				record[columns[i]] = int(val.(int))
			case int64:
				record[columns[i]] = int64(val.(int64))
			case float32:
				record[columns[i]] = float32(val.(float32))
			case float64:
				record[columns[i]] = float64(val.(float64))
			case []byte:
				record[columns[i]] = string(val.([]byte))
			default:
				record[columns[i]] = string(val.([]byte))
			}
		}
		result = append(result, record)
	}
	return result, nil
}

//Delete delete
func (tx *Tx) Delete(sqlStr string, args ...interface{}) (int64, error) {
	var stmt *sql.Stmt
	var err error
	var result sql.Result
	stmt, err = tx.t.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err = stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

//Update update
func (tx *Tx) Update(sqlStr string, args ...interface{}) (int64, error) {
	var stmt *sql.Stmt
	var err error
	var result sql.Result
	stmt, err = tx.t.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err = stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

//Insert insert
func (tx *Tx) Insert(sqlStr string, args ...interface{}) (int64, error) {
	var stmt *sql.Stmt
	var err error
	var result sql.Result
	stmt, err = tx.t.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err = stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	lastInsertID, _ := result.LastInsertId()
	return lastInsertID, nil
}

//Begin begin
func (dao *Dao) Begin() (*Tx, error) {
	t, err := dao.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{t}, nil
}

//Commit commit
func (tx *Tx) Commit() error {
	err := tx.t.Commit()
	if err != nil {
		return err
	}
	return nil
}

//Rollback rollback
func (tx *Tx) Rollback() error {
	err := tx.t.Rollback()
	if err != nil {
		return err
	}
	return nil
}

//Close close
func (dao *Dao) Close() error {
	err := dao.db.Close()
	if err != nil {
		return err
	}
	dao.db = nil
	return nil
}
