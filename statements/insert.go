package statements

import (
	"github.com/askYangc/corm/parse"
)

/**
 * @Author: yc
 * @Description:
 * @File: insert
 * @Date: 2022/5/17 10:56
 */


func (s *Statements) joinInsertPrefix() {
	s.Builder.WriteString("insert into ")
	s.tableNameJoin()
}



func (s *Statements) joinInsertColumns() {
	s.Builder.WriteByte('(')
	if len(s.Columns) == 0 {
		//all
		s.fillInsertColumns()
	}

	for i, _ := range s.Columns {
		if i == 0 {
			s.Builder.WriteString(s.Columns[i])
			continue
		}
		s.Builder.WriteByte(',')
		s.Builder.WriteString(s.Columns[i])
	}
	s.Builder.WriteString(") values(")

	for i, _ := range s.Columns {
		if i == 0 {
			s.Builder.WriteString("?")
			continue
		}
		s.Builder.WriteString(",?")
	}
	s.Builder.WriteByte(')')
}

func (s *Statements) fillInsertColumns() {
	for k, _ := range s.Table.Fields {
		if InSlice(k, s.Omit) {
			continue
		}
/*
		val := s.Value.FieldByName(v.FiledName)
		if s.isZero(val) {
			continue
		}
 */
		s.Columns = append(s.Columns, k)
	}
}

func (s *Statements) GenerateInsertSql() string {
	//insert into xx
	s.joinInsertPrefix()
	//(id,xx,xx)
	s.joinInsertColumns()
	return s.Builder.String()
}

func (s *Statements) GenerateInsertArgs() (args []interface{}){
	for _, v := range s.Columns {
		args = append(args, s.GetColumnsArgs(v))
	}
	return
}

func (s *Statements) Insert(value interface{}) {
	s.Value = parse.GetReflectValue(value)
	s.Table = parse.GetTable(s.Value)
	s.DoAction = ACTION_INSERT

	s.Omit = make([]string, 0)
	s.Omit = append(s.Omit, s.Table.PrimaryTag)
	s.Omit = append(s.Omit, "ctime")
	s.Omit = append(s.Omit, "mtime")
}