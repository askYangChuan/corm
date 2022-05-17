package statements

import "github.com/askYangc/corm/corm/parse"

/**
 * @Author: yc
 * @Description:
 * @File: delete
 * @Date: 2022/5/17 15:48
 */


func (s *Statements) joinDeletePrefix() {
	s.Builder.WriteString("delete from ")
	s.tableNameJoin()
	s.Builder.WriteString(" where ")
}


func (s *Statements) GenerateDeleteSql() string {
	s.Builder.WriteString("delete from ")
	s.tableNameJoin()
	s.Builder.WriteString(" where ")
	s.Builder.WriteString(s.Table.PrimaryTag)
	s.Builder.WriteString("=? limit 1")
	return s.Builder.String()
}

func (s *Statements) GenerateDeleteArgs() (args []interface{}){
	//primaryKey
	args = append(args, s.GetColumnsArgs(s.Table.PrimaryTag))
	return
}

func (s *Statements) Delete(value interface{}) {
	s.Value = parse.GetReflectValue(value)
	s.Table = parse.GetTable(s.Value)
	s.DoAction = ACTION_DELETE

	s.Omit = make([]string, 0)
}