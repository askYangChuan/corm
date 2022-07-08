package statements

import (
	"errors"
	"fmt"
	"github.com/askYangc/corm/parse"
)

/**
 * @Author: yc
 * @Description: get one
 * @File: get
 * @Date: 2022/7/7 15:44
 */



func (s *Statements) Get(value interface{}, args ...interface{}) error {
	s.Value = parse.GetReflectValue(value)
	s.Table = parse.GetTable(s.Value)
	s.DoAction = ACTION_GET

	s.Omit = make([]string, 0)
	s.FuncArgs = args

	if len(args) == 0 {
		return errors.New("corm.Get args is empty")
	}
	return nil
}

func (s *Statements) GenerateGetSql() string {
	s.Builder.WriteString("select * from ")
	s.tableNameJoin()



	s.Builder.WriteString(" where ")
	//use args
	argWhere, ok := s.FuncArgs[0].(string)
	if !ok {
		res := fmt.Errorf("s.FuncArgs[0] is not string, is %+V", s.FuncArgs[0])
		panic(res)
	}

	s.Builder.WriteString(argWhere)

	if !s.hasLimit() && !s.hasForUpdate() {
		s.Builder.WriteString(" limit 1")
	}

	return s.Builder.String()
}


func (s *Statements) GenerateGetArgs() (args []interface{}){
	//from args
	if len(s.FuncArgs) > 1 {
		for _, v := range s.FuncArgs[1:] {
			args = append(args, v)
		}
	}
	return
}