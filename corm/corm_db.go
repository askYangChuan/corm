package corm

import "github.com/jmoiron/sqlx"


//update insert delete in not in join

type Clause struct {

}


type DB struct {
	*sqlx.DB
	Clauses map[string]Clause
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{DB: db}
}

func CloneDB(db *DB) *DB {
	return &DB{DB: db.DB}
}

func (db *DB) getInstance() *DB {
	return db
}

//获取所有数据
func (db *DB) Get(dest interface{}, conds ...interface{}) (tx *DB) {
	tx = db.getInstance()

	return
}

func (db *DB) Where(query interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	return
}

/*
//指定获取一个数据
func (db *DB) Get(dest interface{}) {
	v := getReflectValue(dest)

	table, err := getTable(v)
}


 */

func (db *DB) Create(value interface{}) (tx *DB) {
	v := getReflectValue(dest)

	table, err := getTable(v)
}
