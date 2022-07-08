package statements

import (
	"fmt"
	"github.com/askYangc/corm/parse"
	"strconv"
	"strings"
)

/**
 * @Author: yc
 * @Description: 查询一组数据
 * @File: select
 * @Date: 2022/7/7 16:49
 */

func (s *Statements) hasLimit() bool {
	sql := s.Builder.String()
	if strings.Contains(sql, "limit") {
		return true
	}
	return false
}

func (s *Statements) hasForUpdate() bool {
	//todo: for update必须中间只能一个空格，考虑用正则表达式处理
	sql := s.Builder.String()
	if strings.Contains(sql, "for update") {
		return true
	}
	return false
}


func (s *Statements) Select(value interface{}, args ...interface{}) error {
	//value is a arrays. convert it to internal object
	s.Value = parse.GetSlicePtrReflectValue(value)
	s.Table = parse.GetTable(s.Value)
	s.DoAction = ACTION_SELECT

	s.Omit = make([]string, 0)
	s.FuncArgs = args

	return nil
}

func (s *Statements) GenerateSelectSql() string {
	s.Builder.WriteString("select * from ")
	s.tableNameJoin()

	if len(s.FuncArgs) > 0 {
		//has condtions
		s.Builder.WriteString(" where ")
		//use args
		argWhere, ok := s.FuncArgs[0].(string)
		if !ok {
			res := fmt.Errorf("s.FuncArgs[0] is not string, is %+V", s.FuncArgs[0])
			panic(res)
		}
		s.Builder.WriteString(argWhere)
		s.Builder.WriteByte(' ')
	}

	//set limit
	//check limit
	if !s.hasLimit() && !s.hasForUpdate() {
		if s.LimitNum == 0 {
			return s.Builder.String()
		}
		if s.LimitOffset != 0 {
			s.Builder.WriteString(" limit " + strconv.FormatUint(uint64(s.LimitOffset), 10) +
				"," + strconv.FormatUint(uint64(s.LimitNum), 10))
		}else {
			s.Builder.WriteString(" limit " + strconv.FormatUint(uint64(s.LimitNum), 10))
		}
	}

	return s.Builder.String()
}


func (s *Statements) GenerateSelectArgs() (args []interface{}){
	//from args
	if len(s.FuncArgs) > 1 {
		for _, v := range s.FuncArgs[1:] {
			args = append(args, v)
		}
	}
	return
}