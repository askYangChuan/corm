package parse

import (
	"fmt"
	"github.com/askYangc/corm/utils"
	"log"
	"reflect"
	"sync"
)

var customTables = CustomTables{
	Tables: sync.Map{},
}


type CustomTables struct {
	Tables sync.Map		//map [structName]*SqlTable
}

func (c *CustomTables) Show() {
	fmt.Println("=========")
	c.Tables.Range(func(key, value interface{}) bool {
		v := value.(*SqlTable)
		v.Show()
		return true
	})

	fmt.Println("=========")
}

type SqlField struct {
	FiledName string		//变量名
	ColumnName string		//列名 tags
	Field reflect.StructField
}

func (field *SqlField) Show() {
	fmt.Printf("\t=====%s=====\n", field.FiledName)
	fmt.Printf("\tFiledName: %s, columnName: %s\n", field.FiledName, field.ColumnName)
}

type SqlTable struct {
	StructName string
	TableName string
	PrimaryTag string	//tags id
	Fields map[string]SqlField		//key tags
}

func (table *SqlTable) Show() {
	fmt.Printf("=====%s==%s==\n", table.StructName, table.TableName)

	for _, v := range table.Fields {
		v.Show()
	}

	fmt.Println("=====end====")
}

func (table *SqlTable) parseStructFieldWithTag(t reflect.StructField, tags string) {
	switch t.Type.Kind() {
	case reflect.Ptr:
		log.Panicf("corm not support tag with ptr, failed Filed: %s\n", t.Name)
	}

	pkey := t.Tag.Get("corm")
	if pkey == "primaryKey" {
		table.PrimaryTag = tags
	}

	table.Fields[tags] = SqlField{
		FiledName:  t.Name,
		ColumnName: tags,
		Field:      t,
	}
}

func (table *SqlTable) parseStructFieldWithNoTag(t reflect.StructField) {
	typ := t.Type

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		log.Panicf("parseStructFieldWithNoTag typ Kind not struct, is %v", typ.Kind())
	}

	table.ParseStruct(typ)
}

func (table *SqlTable) parseStructField(t reflect.StructField) {
	tags := t.Tag.Get("db")
	if tags != "" {
		//get tags
		table.parseStructFieldWithTag(t, tags)
		return
	}


	table.parseStructFieldWithNoTag(t)
}

func (table *SqlTable) ParseStruct(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		table.parseStructField(t.Field(i))
	}
}

func GetReflectValue(val interface{}) reflect.Value {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	//case reflect.Struct:
	default:
		log.Panicf("v Kind not ptr or struct, is %v", v.Kind())
	}

	return v
}

//return a new reflect.Value
func GetSlicePtrReflectValue(val interface{}) reflect.Value {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		log.Panicf("v Kind not ptr, is %v", v.Kind())
	}

	v = v.Elem()
	//now is slice
	if v.Kind() != reflect.Slice {
		log.Panicf("v Kind not slice, is %v", v.Kind())
	}

	newV := reflect.New(v.Type().Elem())
	return newV.Elem()
}

func getTableNameByVal(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Struct:
		vPtr := reflect.New(v.Type())
		vPtr.Elem().Set(v)
		v = vPtr
	case reflect.Ptr:
	default:
		log.Panicf("getTableNameByVal v Kind not ptr or struct, is %v", v.Kind())
	}

	nameFunc := v.MethodByName("TableName")
	if nameFunc.IsValid() {
		res := nameFunc.Call([]reflect.Value{})
		return res[0].Interface().(string)
	}

	return utils.SnakeString(v.Elem().Type().Name())
}

//v is struct
func ParseTable(v reflect.Value) *SqlTable{
	t := v.Type()


	sqlTableVal, ok := customTables.Tables.Load(t.Name())
	if ok {
		return sqlTableVal.(*SqlTable)
	}

	table := &SqlTable{
		StructName: t.Name(),
		TableName:  getTableNameByVal(v),
		PrimaryTag: "id",
		Fields:     make(map[string]SqlField, 0),
	}

	table.ParseStruct(t)
	customTables.Tables.Store(t.Name(), table)
	return table
}

func ShowAll() {
	customTables.Show()
}

func Show(args ...interface{}) {
	if len(args) == 0 {
		ShowAll()
		return
	}

	for i, _ := range args {
		v := GetReflectValue(args[i])
		val, ok := customTables.Tables.Load(v.Type().Name())
		if ok {
			realVal := val.(*SqlTable)
			realVal.Show()
		}
	}
}

func GetTable(v reflect.Value) *SqlTable {
	return ParseTable(v)
}


//========================do========================

func (table *SqlTable) GetSqlField(tag string) *SqlField{
	if v, ok := table.Fields[tag]; ok {
		return &v
	}

	return nil
}