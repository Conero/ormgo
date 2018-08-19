package sqler

import (
	"testing"
	"fmt"
)

// @Date：   2018/8/18 0018 18:51
// @Author:  Joshua Conero
// @Name:    名称描述


func TestMysqlCreator_Where(t *testing.T) {
	sqler := new (MysqlCreator)

	sql := sqler.From("user").
		Where("name like ?").
		Select()
	fmt.Println(sql)

	_, crt := NewCreator("mysql")
	sql = crt.From("user", "a").Count()
	fmt.Println(sql)

	// join
	sql = new(MysqlCreator).From("user", "us").
		Field("us.name", "us.gender", "rs.text", "cp.name").
		Join("resume", "rs", "us.uid=rs.uid").
		LeftJoin("company", "cp", "rs.resume_id=cp.resume_id").
		Order("us.name", "us.gender desc", "cp.name desc").
		Group("us.name", "us.gender", "rs.text", "cp.name").
		Page(1, 50).
		Select()
	fmt.Println(sql)
}

func TestMysqlCreator_Count(t *testing.T) {

}