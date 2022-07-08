package corm

import (
	"errors"
	"fmt"
	"github.com/askYangc/corm/logging"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"reflect"
	"time"
)

func (db *DB) updateInsertValue() {
	tx := db.getInstance()
	lastId, _ := tx.Result.LastInsertId()
	val := tx.Statement.Value

	//id
	field := tx.Statement.Table.GetSqlField(tx.Statement.Table.PrimaryTag)
	if field != nil {
		value := val.FieldByName(field.FiledName)
		value.Set(reflect.ValueOf(uint64(lastId)))
	}

	t := time.Now()

	//ctime
	field = tx.Statement.Table.GetSqlField("ctime")
	if field != nil {
		value := val.FieldByName(field.FiledName)
		value.Set(reflect.ValueOf(t))
	}

	//mtime
	field = tx.Statement.Table.GetSqlField("mtime")
	if field != nil {
		value := val.FieldByName(field.FiledName)
		value.Set(reflect.ValueOf(t))
	}
}

func (db *DB) updateInsertOrUpdateValue() {
	tx := db.getInstance()
	lastId, err := tx.Result.LastInsertId()
	if err != nil {
		return
	}
	affected, err := tx.Result.RowsAffected()
	if err != nil {
		logging.ZError("corm updateInsertOrUpdateValue failed", zap.Error(err))
		return
	}

	if affected == 0 {
		//nothing to update, but not get id
		logging.ZDebug("corm updateInsertOrUpdateValue get affected 0")
		return
	}
	if affected == 2 {
		//do update
		logging.ZDebug("corm updateInsertOrUpdateValue get affected 2")
		tx.updateUpdateValue()
		return
	}

	//insert
	val := tx.Statement.Value
	if lastId != 0 {
		//id
		field := tx.Statement.Table.GetSqlField(tx.Statement.Table.PrimaryTag)
		if field != nil {
			value := val.FieldByName(field.FiledName)
			value.Set(reflect.ValueOf(uint64(lastId)))
		}
	}

	t := time.Now()
	//ctime
	field := tx.Statement.Table.GetSqlField("ctime")
	if field != nil {
		value := val.FieldByName(field.FiledName)
		value.Set(reflect.ValueOf(t))
	}

	//mtime
	field = tx.Statement.Table.GetSqlField("mtime")
	if field != nil {
		value := val.FieldByName(field.FiledName)
		value.Set(reflect.ValueOf(t))
	}
}


func (db *DB) updateUpdateValue() {
	tx := db.getInstance()
	val := tx.Statement.Value

	affected, err := tx.Result.RowsAffected()
	if err != nil || affected == 0 {
		logging.ZWarn("updateUpdateValue get err", zap.Int64("affected", affected), zap.Error(err))
		return
	}

	//mtime
	field := tx.Statement.Table.GetSqlField("mtime")
	if field != nil {
		value := val.FieldByName(field.FiledName)
		value.Set(reflect.ValueOf(time.Now()))
	}
}


//insert into xxx
func (db *DB) Insert(value interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Insert(value)

	var useTx zap.Field
	//insert
	sqlStr, args := tx.Join()
	if tx.Tx != nil {
		useTx = zap.Bool("useTx", true)
		tx.Result, tx.Error = tx.Tx.Exec(sqlStr, args...)
	}else {
		useTx = zap.Bool("useTx", false)
		tx.Result, tx.Error = tx.DB.Exec(sqlStr, args...)
	}

	if tx.Error != nil {
		logging.ZError("corm Insert failed",
			zap.String("sql", sqlStr), zap.Any("args", args), useTx, zap.Error(tx.Error))
		tx.Statement.Reset()
		return tx
	}
	logging.ZDebug("Insert sql", zap.String("sql", sqlStr), zap.Any("args", args), useTx)

	tx.updateInsertValue()
	tx.Statement.Reset()
	return tx
}

//if db.ID == 0, insert
func (db *DB) Update(value interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Update(value)

	var useTx zap.Field
	if tx.Statement.PrimaryKeyIsZero() {
		//is zero
		logging.ZDebug("corm Update PrimaryKey is zero. do Insert")
		tx.Insert(value)
		return tx
	}

	//update
	sqlStr, args := tx.Join()
	if tx.Tx != nil {
		useTx = zap.Bool("useTx", true)
		tx.Result, tx.Error = tx.Tx.Exec(sqlStr, args...)
	}else {
		useTx = zap.Bool("useTx", false)
		tx.Result, tx.Error = tx.DB.Exec(sqlStr, args...)
	}

	if tx.Error != nil {
		logging.ZError("corm Update failed",
			zap.String("sql", sqlStr), zap.Any("args", args), useTx, zap.Error(tx.Error))
		tx.Statement.Reset()
		return tx
	}

	logging.ZDebug("Update sql", zap.String("sql", sqlStr), zap.Any("args", args), useTx)
	tx.updateUpdateValue()
	tx.Statement.Reset()
	return tx
}

func (db *DB) InsertOrUpdate(value interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.InsertOrUpdate(value)
	var useTx zap.Field

	//insert or update
	sqlStr, args := tx.Join()
	if tx.Tx != nil {
		useTx = zap.Bool("useTx", true)
		tx.Result, tx.Error = tx.Tx.Exec(sqlStr, args...)
	}else {
		useTx = zap.Bool("useTx", false)
		tx.Result, tx.Error = tx.DB.Exec(sqlStr, args...)
	}

	if tx.Error != nil {
		logging.ZError("corm InsertOrUpdate failed",
			zap.String("sql", sqlStr), zap.Any("args", args), useTx, zap.Error(tx.Error))
		tx.Statement.Reset()
		return tx
	}

	logging.ZDebug("InsertOrUpdate sql", zap.String("sql", sqlStr), zap.Any("args", args), useTx)
	tx.updateInsertOrUpdateValue()
	tx.Statement.Reset()
	return tx
}

func (db *DB) Delete(value interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Statement.Delete(value)
	var useTx zap.Field

	if tx.Statement.PrimaryKeyIsZero() {
		//is zero
		tx.Error = fmt.Errorf("value  primaryKey %s is Zero", tx.Statement.Table.PrimaryTag)
		logging.ZError("corm Delete failed", zap.Error(tx.Error))
		return tx
	}

	//delete
	sqlStr, args := tx.Join()
	if tx.Tx != nil {
		useTx = zap.Bool("useTx", true)
		tx.Result, tx.Error = tx.Tx.Exec(sqlStr, args...)
	}else {
		useTx = zap.Bool("useTx", false)
		tx.Result, tx.Error = tx.DB.Exec(sqlStr, args...)
	}

	if tx.Error != nil {
		logging.ZError("corm Delete failed",
			zap.String("sql", sqlStr), zap.Any("args", args), useTx, zap.Error(tx.Error))
		tx.Statement.Reset()
		return tx
	}

	logging.ZDebug("Delete sql", zap.String("sql", sqlStr), zap.Any("args", args), useTx)
	tx.Statement.Reset()
	return tx
}

func (db *DB) Get(value interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Error = tx.Statement.Get(value, args...)
	if tx.Error != nil {
		logging.ZError("corm Get failed", zap.Error(tx.Error))
		return tx
	}
	var useTx zap.Field

	//get
	sqlStr, args := tx.Join()
	if tx.Tx != nil {
		useTx = zap.Bool("useTx", true)
		tx.Error = tx.Tx.Get(value, sqlStr, args...)
	}else {
		useTx = zap.Bool("useTx", false)
		tx.Error = tx.DB.Get(value, sqlStr, args...)
	}

	if tx.Error != nil {
		logging.ZError("corm Get failed",
			zap.String("sql", sqlStr), zap.Any("args", args), useTx, zap.Error(tx.Error))
		tx.Statement.Reset()
		return tx
	}

	logging.ZDebug("Get sql", zap.String("sql", sqlStr), zap.Any("args", args), useTx)
	tx.Statement.Reset()
	return tx
}

func (db *DB) Select(value interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	tx.Error = tx.Statement.Select(value, args...)
	if tx.Error != nil {
		logging.ZError("corm Select failed", zap.Error(tx.Error))
		return tx
	}
	var useTx zap.Field

	//select
	sqlStr, args := tx.Join()
	if tx.Tx != nil {
		useTx = zap.Bool("useTx", true)
		tx.Error = tx.Tx.Select(value, sqlStr, args...)
	}else {
		useTx = zap.Bool("useTx", false)
		tx.Error = tx.DB.Select(value, sqlStr, args...)
	}

	if tx.Error != nil {
		logging.ZError("corm Select failed",
			zap.String("sql", sqlStr), zap.Any("args", args), useTx, zap.Error(tx.Error))
		tx.Statement.Reset()
		return tx
	}

	logging.ZDebug("Select sql", zap.String("sql", sqlStr), zap.Any("args", args), useTx)
	tx.Statement.Reset()
	return tx
}

func (db *DB) Limit(num uint32, args ...uint32) (tx *DB) {
	tx = db.getInstance()
	if len(args) == 0 {
		if num == 0 {
			tx.Error = errors.New("corm.Limit is 0")
			return tx
		}
		tx.Statement.SetLimit(0, num)
	}else {
		if len(args) != 1 {
			tx.Error = errors.New("corm.Limit just accept 2 args")
			return tx
		}
		tx.Statement.SetLimit(num, args[0])
	}

	logging.ZDebug("Set Limit", zap.Uint32s("limit", []uint32{tx.Statement.LimitOffset, tx.Statement.LimitNum}))
	return tx
}


func (db *DB) Join() (string, []interface{}) {
	return db.Statement.Join()
}


//support trans
func (db *DB) Beginx() (tx *DB, err error) {
	tx = db.getInstance()

	tx.Tx, err = tx.DB.Beginx()
	if err != nil {
		logging.ZError("corm Beginx failed", zap.Error(err))
	}
	return
}

func (db *DB) Commit() (tx *DB, err error) {
	tx = db.getInstance()

	err = tx.Tx.Commit()
	if err != nil {
		logging.ZError("corm Commit failed", zap.Error(err))
	}
	return
}

func (db *DB) Rollback() (tx *DB, err error) {
	tx = db.getInstance()

	err = tx.Tx.Rollback()
	if err != nil {
		logging.ZError("corm Rollback failed", zap.Error(err))
	}
	return
}

func (db *DB) SetTx(extraTx *sqlx.Tx) (tx *DB) {
	tx = db.getInstance()

	tx.Tx = extraTx
	return
}