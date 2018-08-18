package sqler

import (
	"strings"
	"errors"
)

// @Date：   2018/8/18 0018 18:35
// @Author:  Joshua Conero
// @Name:    sql 生成器

type Creator interface {
	Field(fields ...string) Creator
	// table, alias
	From(tables ...string) Creator

	// table, alias, cond
	Join(tables ...string) Creator

	LeftJoin(tables ...string) Creator
	RightJoin(tables ...string) Creator
	Where(cond string) Creator
	OrWhere(cond string) Creator

	Select() string
	Count(fields ...string) string
}

func NewCreator(driver string) (error, Creator)  {
	var creator Creator
	var err error
	switch strings.ToLower(driver) {
	case "mysql":
		creator = new(MysqlCreator)
	default:
		err = errors.New("driver not find")
	}
	return err, creator
}