package corm

import (
	"corms/utils"
	"fmt"
	"log"
	"reflect"
)

var customTables = CustomTables{
	Tables: make(map[string]SqlTable, 0),
}


type CustomTables struct {
	Tables map[string]SqlTable
}

func (c *CustomTables) Show() {
	fmt.Println("=========")
	for _, v := range c.Tables {
		v.Show()
	}
	fmt.Println("=========")
}

type SqlField struct {
	FiledName string
	ColumnName string
	Field reflect.StructField
}

func (field *SqlField) Show() {
	fmt.Printf("\t=====%s=====\n", field.FiledName)
	fmt.Printf("\tFiledName: %s, columnName: %s\n", field.FiledName, field.ColumnName)
}

type SqlTable struct {
	StructName string
	TableName string
	Fields map[string]SqlField
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

func getReflectValue(val interface{}) reflect.Value {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Struct:
	default:
		log.Panicf("v Kind not ptr or struct, is %v", v.Kind())
	}

	return v
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

func ParseTable(val interface{}) {
	v := getReflectValue(val)
	t := v.Type()

	_, ok := customTables.Tables[t.Name()]
	if ok {
		return
	}

	table := SqlTable{
		StructName: t.Name(),
		TableName:  getTableNameByVal(v),
		Fields:     make(map[string]SqlField, 0),
	}

	table.ParseStruct(t)
	customTables.Tables[t.Name()] = table
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
		v := getReflectValue(args[i])
		val, ok := customTables.Tables[v.Type().Name()]
		if ok {
			val.Show()
		}
	}
}

func getTable(v reflect.Value) *SqlTable{

}
