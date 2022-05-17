package statements

import "github.com/askYangc/corm/parse"

/**
 * @Author: yc
 * @Description:
 * @File: insert_or_update
 * @Date: 2022/5/17 16:00
 */


func (s *Statements) joinDuplicate() {
	s.Builder.WriteString(" ON DUPLICATE KEY UPDATE ")
	for k, v := range s.Columns {
		if k == 0 {
			s.Builder.WriteString(v)
			s.Builder.WriteString("=values(")
			s.Builder.WriteString(v)
			s.Builder.WriteString(")")
			continue
		}
		s.Builder.WriteByte(',')
		s.Builder.WriteString(v)
		s.Builder.WriteString("=values(")
		s.Builder.WriteString(v)
		s.Builder.WriteString(")")
	}
}

func (s *Statements) GenerateInsertorUpdateSql() string {
	//insert into xx
	s.joinInsertPrefix()
	//(id,xx,xx)
	s.joinInsertColumns()
	//ON DUPLICATE KEY UPDATE
	s.joinDuplicate()
	return s.Builder.String()
}

func (s *Statements) GenerateInsertorUpdateArgs() (args []interface{}){
	for _, v := range s.Columns {
		args = append(args, s.GetColumnsArgs(v))
	}
	return
}

func (s *Statements) InsertOrUpdate(value interface{}) {
	s.Value = parse.GetReflectValue(value)
	s.Table = parse.GetTable(s.Value)
	s.DoAction = ACTION_INSERTORUPDATE

	s.Omit = make([]string, 0)
	s.Omit = append(s.Omit, s.Table.PrimaryTag)
	s.Omit = append(s.Omit, "ctime")
	s.Omit = append(s.Omit, "mtime")
}