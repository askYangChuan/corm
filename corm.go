package corm

import (
	"database/sql"
	"github.com/askYangc/corm/statements"
	"github.com/jmoiron/sqlx"
)

//update insert delete in not in join
type DB struct {
	*sqlx.DB
	clone int
	Statement    statements.Statements

	//result
	Result sql.Result
	Error error
}

var (
	db *DB
)

func NewDB(db *sqlx.DB) *DB {
	return &DB{DB: db, clone : 1}
}

func CloneDB(db *DB) *DB {
	return &DB{DB: db.DB}
}

func (db *DB) getInstance() *DB {
	if db.clone > 0 {
		return CloneDB(db)
	}
	return db
}

func Insert(value interface{}) error {
	tx := db.Insert(value)
	return tx.Error
}

func Update(value interface{}) error {
	tx := db.Update(value)
	return tx.Error
}

func InsertOrUpdate(value interface{}) error {
	tx := db.InsertOrUpdate(value)
	return tx.Error
}

func Delete(value interface{}) error {
	tx := db.Delete(value)
	return tx.Error
}

//add get, condtions, args
func Get(value interface{}, args ...interface{}) error {
	tx := db.Get(value, args...)
	return tx.Error
}

//add select, condtions, args, just support all cloumns
func Select(value interface{}, args ...interface{}) error {
	tx := db.Select(value, args...)
	return tx.Error
}

func Limit(num uint32, args ...uint32) (tx *DB) {
	tx = db.Limit(num, args...)
	return tx
}



func CormInit(sqlDb *sqlx.DB) {
	db = NewDB(sqlDb)
}


