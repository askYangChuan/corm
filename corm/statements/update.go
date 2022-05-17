package statements

import "github.com/askYangc/corm/corm/parse"

/**
 * @Author: yc
 * @Description:
 * @File: update
 * @Date: 2022/5/17 15:09
 */

func (s *Statements) joinUpdatePrefix() {
	s.Builder.WriteString("update ")
	s.tableNameJoin()
	s.Builder.WriteString(" set ")
}

func (s *Statements) fillUpdateColumns() {
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


func (s *Statements) joinUpdateColumns() {
	if len(s.Columns) == 0 {
		//all
		s.fillUpdateColumns()
	}

	for i, _ := range s.Columns {
		if i == 0 {
			s.Builder.WriteString(s.Columns[i])
			s.Builder.WriteString("=?")
			continue
		}
		s.Builder.WriteByte(',')
		s.Builder.WriteString(s.Columns[i])
		s.Builder.WriteString("=?")
	}
	s.Builder.WriteString(" where ")
	s.Builder.WriteString(s.Table.PrimaryTag)
	s.Builder.WriteString("=? limit 1")
}

func (s *Statements) GenerateUpdateSql() string {
	//update %s set
	s.joinUpdatePrefix()
	//xx=?,xx=?
	s.joinUpdateColumns()
	return s.Builder.String()
}

func (s *Statements) GenerateUpdateArgs() (args []interface{}){
	for _, v := range s.Columns {
		args = append(args, s.GetColumnsArgs(v))
	}

	//primaryKey
	args = append(args, s.GetColumnsArgs(s.Table.PrimaryTag))
	return
}

func (s *Statements) Update(value interface{}) {
	s.Value = parse.GetReflectValue(value)
	s.Table = parse.GetTable(s.Value)
	s.DoAction = ACTION_UPDATE

	s.Omit = make([]string, 0)
	s.Omit = append(s.Omit, s.Table.PrimaryTag)
	s.Omit = append(s.Omit, "ctime")
	s.Omit = append(s.Omit, "mtime")
}