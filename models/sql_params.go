package models

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)


type SqlParam struct {
	SqlTParamType uint8			//类型	0, 表示是条件;1,表示后缀, 2 表示limit 3 表示in
	Conditions string		//具体条件
	Value interface{}
}

type SqlIn struct {
	Condition string
	Num int		//参数个数
	Value []interface{}
}

func NewSqlIn(condition string) *SqlIn {
	return &SqlIn{
		Condition: condition,
		Num:   0,
		Value: make([]interface{}, 0),
	}
}

func (s *SqlIn) GetInSql() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s in (?", s.Condition))

	for i := 1; i < s.Num; i++ {
		b.WriteString(",?")
	}
	b.WriteByte(')')
	return b.String()
}

func (s *SqlIn) Append(v interface{}) {
	s.Value = append(s.Value, v)
	s.Num++
}



type SqlParams struct {
	TableName []string		//表名
	Columns []string		//字段
	Params []SqlParam		//查询条件
}

func NewSqlParams(tableName ...string) *SqlParams {
	return &SqlParams{
		TableName: tableName,
		Columns: make([]string, 0),
		Params: make([]SqlParam, 0),
	}
}

func CopyNewSqlParams(sqlParams *SqlParams) *SqlParams {
	//有可实现深度拷贝的通用函数，后续再加入
	p := NewSqlParams(sqlParams.TableName...)
	copy(p.Columns, sqlParams.Columns)
	copy(p.Params, sqlParams.Params)
	return p
}

func (s *SqlParams) Clear() {
	s.TableName = make([]string, 0)
	s.Columns = make([]string, 0)
	s.Params = make([]SqlParam, 0)
}

func (s *SqlParams) SetTableName(tableName ...string) {
	s.TableName = tableName
}

func (s *SqlParams) AddColumns(columns ...string) {
	for _, column := range columns {
		s.Columns = append(s.Columns, column)
	}
}

func (s *SqlParams) Append(param SqlParam) {
	s.Params = append(s.Params, param)
}

func (s *SqlParams) AppendCondition(condition string, v interface{}) {
	s.Append(SqlParam{
		SqlTParamType: 0,
		Conditions:    condition,
		Value:         v,
	})
}

func (s *SqlParams) AppendIn(in *SqlIn) {
	s.Append(SqlParam{
		SqlTParamType: 3,
		Conditions:    "",
		Value:         in,
	})
}

func (s *SqlParams) AppendGroupBy(condition string, v interface{}) {
	s.Append(SqlParam{
		SqlTParamType: 1,
		Conditions:    condition,
		Value:         v,
	})
}


func (s *SqlParams) AppendOrderBy(condition string, v interface{}) {
	s.Append(SqlParam{
		SqlTParamType: 1,
		Conditions:    condition,
		Value:         v,
	})
}


//常用设置page的函数
func (s *SqlParams) SetPageParams(page, pageSize uint32) {
	if page != 0 {
		s.Append(SqlParam{
			SqlTParamType: 2,
			Conditions:    fmt.Sprintf("limit %d,%d", (page-1)*pageSize, pageSize),
			Value:         nil,
		})
	}
}

//常用设置page的函数
func (s *SqlParams) SetLimit(n uint32) {
	s.Append(SqlParam{
		SqlTParamType: 2,
		Conditions:    fmt.Sprintf("limit %d", n),
		Value:         nil,
	})
}

func (s *SqlParams) tableNameJoin() string {
	var buffer bytes.Buffer
	for i, name := range s.TableName {
		if i == 0 {
			buffer.WriteString(name)
			continue
		}
		buffer.WriteString(",")
		buffer.WriteString(name)
	}
	return buffer.String()
}

func (s *SqlParams) columnsJoin() string {
	var columns bytes.Buffer
	for i, column := range s.Columns {
		if i == 0 {
			columns.WriteString(column)
			continue
		}
		columns.WriteString(",")
		columns.WriteString(column)
	}
	return columns.String()
}

//构建前缀
func (s *SqlParams) joinPrefix() string {
	if len(s.Columns) == 0 {
		return fmt.Sprintf("select * from %s", s.tableNameJoin())
	}
	return fmt.Sprintf("select %s from %s", s.columnsJoin(), s.tableNameJoin())
}

//构建获取记录的前缀
func (s *SqlParams) joinCntPrefix() string {
	return fmt.Sprintf("select count(1) from %s", s.tableNameJoin())
}

func (s *SqlParams) GetSql() string {
	var b strings.Builder
	b.WriteString(s.joinPrefix())

	for i, param := range s.Params {
		switch param.SqlTParamType {
		case 0:
			// where
			if i != 0 {
				b.WriteString("and")
			}else {
				b.WriteString(" where")
			}
			b.WriteByte(' ')
			b.WriteString(param.Conditions)
			b.WriteByte(' ')
		case 1:
			//suffix
			fallthrough
		case 2:
			//limit
			b.WriteByte(' ')
			b.WriteString(param.Conditions)
		case 3:
			//in
			if i != 0 {
				b.WriteString("and")
			}else {
				b.WriteString(" where")
			}
			in := param.Value.(*SqlIn)
			b.WriteByte(' ')
			b.WriteString(in.GetInSql())
		}
	}
	return b.String()
}

func (s *SqlParams) Join() (string, []interface{}, error) {
	log.Println(s.GetSql())
	log.Println(s.GetArgs())
	return s.GetSql(), s.GetArgs(), nil
}

func (s *SqlParams) GetArgs() (args []interface{}){
	for _, param := range s.Params {
		if param.Value != nil {
			if param.SqlTParamType == 3 {
				//in
				in := param.Value.(*SqlIn)
				for i := 0; i < in.Num; i++ {
					args = append(args, in.Value[i])
				}
				continue
			}

			args = append(args, param.Value)
		}


	}
	return args
}

func (s *SqlParams) Get(v interface{}) error {
	sqlStr, args, err := s.Join()
	err = db.Get(v, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlParams) Select(v interface{}) error {
	sqlStr, args, err := s.Join()
	err = db.Select(v, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}


//获取满足所有条件的记录个数
func (s *SqlParams) GetTotalCount() (totalCount uint32, err error){
	var b strings.Builder
	b.WriteString(s.joinCntPrefix())
	for i, param := range s.Params {
		switch param.SqlTParamType {
		case 0:
			// where
			if i != 0 {
				b.WriteString("and")
			}else {
				b.WriteString(" where")
			}
			b.WriteByte(' ')
			b.WriteString(param.Conditions)
			b.WriteByte(' ')
		case 1:
			//suffix
		case 2:
			//limit
		case 3:
			//in
			if i != 0 {
				b.WriteString("and")
			}else {
				b.WriteString(" where")
			}
			in := param.Value.(*SqlIn)
			b.WriteByte(' ')
			b.WriteString(in.GetInSql())
		}
	}

	err = db.QueryRow(b.String(), s.GetArgs()...).Scan(&totalCount)
	return totalCount, err
}