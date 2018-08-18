package sqler

import (
	"regexp"
	"strings"
)

// @Date：   2018/8/18 0018 18:41
// @Author:  Joshua Conero
// @Name:    名称描述

type MysqlCreator struct {
	table  string
	alias  string
	fields []string
	join   []map[string]string
	cond   []map[string]string
}

// 链式方法
func (c *MysqlCreator) From(tables ...string) Creator {
	c.table = tables[0]
	if len(tables) > 1 {
		c.alias = tables[1]
	}
	return c
}

func (c *MysqlCreator) Field(fields ...string) Creator {
	if c.fields == nil {
		c.fields = []string{}
	}
	c.fields = append(c.fields, fields...)
	return c
}

func (c *MysqlCreator) toJoin(vtype string, tables ...string) {
	if c.join == nil {
		c.join = []map[string]string{}
	}
	joinInfo := map[string]string{
		"table": tables[0],
		"alias": tables[1],
		"cond":  tables[2],
		"type":  vtype,
	}
	c.join = append(c.join, joinInfo)
}

func (c *MysqlCreator) Join(tables ...string) Creator {
	c.toJoin("inner", tables...)
	return c
}

func (c *MysqlCreator) LeftJoin(tables ...string) Creator {
	c.toJoin("left", tables...)
	return c
}

func (c *MysqlCreator) RightJoin(tables ...string) Creator {
	c.toJoin("right", tables...)
	return c
}

func (c *MysqlCreator) toCond(vtype string, cond string) {
	if c.cond == nil {
		c.cond = []map[string]string{}
	}
	c.cond = append(c.cond, map[string]string{"type": vtype, "cond": cond})
}

func (c *MysqlCreator) Where(cond string) Creator {
	c.toCond("and", cond)
	return c
}

func (c *MysqlCreator) OrWhere(cond string) Creator {
	c.toCond("or", cond)
	return c
}

func (c *MysqlCreator) Page(page, size int) Creator{
	// @TODO 分页实现
	return c
}

// 数据处理--------->

// 字段解析
func (c *MysqlCreator) parseName(name string) string {
	if strings.Index(name, "`") == -1 {
		ptCk, ptRpl := "[\\.=]", "[\\.=]"
		if matched, _ := regexp.MatchString(ptCk, name); matched{
			ptRplIns := regexp.MustCompile(ptRpl)
			regIns := regexp.MustCompile("\\.[^\\.\\s=]+[=\\s]{0,}")
			for _, str := range regIns.FindAllString(name, -1){
				repl := ptRplIns.ReplaceAllString(str, "")
				repl = ".`"+repl+"`"
				// 结尾为 =
				if str[len(str)-1:] == "="{
					repl += "="
				}
				name = strings.Replace(name, str, strings.TrimSpace(repl), -1)
			}
		}else{
			name = "`" + strings.TrimSpace(name) + "`"
		}
	}

	return name
}
// 数据表解析
func (c *MysqlCreator) parseTable() string {
	table := c.parseName(c.table)
	if c.alias != "" {
		table += " " + c.alias
	}
	// join 联表
	joinQueue := []string{}
	for _, joinMap := range c.join{
		joinStr := strings.ToUpper(joinMap["type"]) + " JOIN " +
			c.parseName(joinMap["table"]) + " " + joinMap["alias"] + " ON " +
			c.parseName(joinMap["cond"])
		joinQueue = append(joinQueue, joinStr)
	}
	if len(joinQueue) > 0{
		table += " " + strings.Join(joinQueue, " ")
	}
	return table
}
// where 条件解析
func (c *MysqlCreator) parseWhere() string {
	where := ""
	whereQuque := []string{}
	for _, whMap := range c.cond{
		if len(whereQuque) == 0{
			if whMap["vtype"] == "and"{
				whereQuque = append(whereQuque, whMap["where"])
			}
		}else {
			whereQuque = append(whereQuque, strings.ToUpper(whMap["type"]) + "("+whMap["where"]+")")
		}
	}
	if len(whereQuque) > 0{
		where = strings.Join(whereQuque, " ")
	}
	if where != "" {
		where += " WHERE " + where
	}
	return where
}
func (c *MysqlCreator) parseField() string {
	field := ""
	fieldQueue := []string{}
	pattern := "/[\\(\\)`]/"
	for _, fld := range c.fields {
		if fld != "*"{
			if matched, err := regexp.MatchString(pattern, fld); !matched && err == nil {
				idx := strings.Index(fld, ".")
				if idx > -1 {
					fld = fld[0:idx] + ".`" + fld[idx+1:] + "`"
				} else {
					fld = "`" + fld + "`"
				}
			}
		}
		fieldQueue = append(fieldQueue, fld)
	}
	if len(fieldQueue) > 0 {
		field = strings.Join(fieldQueue, ",")
	}
	if field == "" {
		field = "*"
	}
	return field
}

func (c *MysqlCreator) Select() string {
	sqlStr := "SELECT " + c.parseField() + " FROM " +
		c.parseTable() + "" +
		c.parseWhere()
	return sqlStr
}

func (c *MysqlCreator) Count(fields ...string) string {
	field := "*"
	sqlStr := "SELECT COUNT(" + field + ") FROM " + c.parseTable() + ""
	return sqlStr
}
