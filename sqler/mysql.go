package sqler

import (
	"regexp"
	"strings"
	"strconv"
)

// @Date：   2018/8/18 0018 18:41
// @Author:  Joshua Conero
// @Name:    名称描述

type MysqlCreator struct {
	_table  string
	_alias  string
	_fields []string
	_join   []map[string]string
	_cond   []map[string]string
	_group []string
	_order []string
	_page int
	_size int
}

// 链式方法
func (c *MysqlCreator) From(tables ...string) Creator {
	c._table = tables[0]
	if len(tables) > 1 {
		c._alias = tables[1]
	}
	return c
}

func (c *MysqlCreator) Field(fields ...string) Creator {
	if c._fields == nil {
		c._fields = []string{}
	}
	c._fields = append(c._fields, fields...)
	return c
}

func (c *MysqlCreator) toJoin(vtype string, tables ...string) {
	if c._join == nil {
		c._join = []map[string]string{}
	}
	joinInfo := map[string]string{
		"table": tables[0],
		"alias": tables[1],
		"cond":  tables[2],
		"type":  vtype,
	}
	c._join = append(c._join, joinInfo)
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
	if c._cond == nil {
		c._cond = []map[string]string{}
	}
	c._cond = append(c._cond, map[string]string{"type": vtype, "cond": cond})
}

func (c *MysqlCreator) Where(cond string) Creator {
	c.toCond("and", cond)
	return c
}

func (c *MysqlCreator) OrWhere(cond string) Creator {
	c.toCond("or", cond)
	return c
}

func (c *MysqlCreator) Page(page, size int) Creator {
	c._page = page
	c._size = size
	return c
}

func (c *MysqlCreator) Group(fields ...string) Creator {
	if c._group == nil{
		c._group = []string{}
	}
	c._group = append(c._group, fields...)
	return c
}

func (c *MysqlCreator) Order(fields ...string) Creator {
	if c._order == nil{
		c._order = []string{}
	}
	c._order = append(c._order, fields...)
	return c
}

// 数据处理--------->

// 字段解析
func (c *MysqlCreator) parseName(name string) string {
	if strings.Index(name, "`") == -1 {
		ptCk, ptRpl := "[\\.=]", "[\\.=]"
		if matched, _ := regexp.MatchString(ptCk, name); matched {
			ptRplIns := regexp.MustCompile(ptRpl)
			regIns := regexp.MustCompile("\\.[^\\.\\s=]+[=\\s]{0,}")
			for _, str := range regIns.FindAllString(name, -1) {
				repl := ptRplIns.ReplaceAllString(str, "")
				repl = ".`" + strings.TrimSpace(repl) + "`"
				// 结尾为 =
				strLast := str[len(str)-1:]
				lastIsSpace := false
				if strLast == "=" {
					repl += "="
				}else if strLast == " "{
					lastIsSpace = true
				}
				repl = strings.TrimSpace(repl)
				if lastIsSpace{
					repl += " "
				}
				name = strings.Replace(name, str, repl, -1)
			}
		} else {
			name = "`" + strings.TrimSpace(name) + "`"
		}
	}

	return name
}

// 数据表解析
func (c *MysqlCreator) parseTable() string {
	table := c.parseName(c._table)
	if c._alias != "" {
		table += " " + c._alias
	}
	// _join 联表
	joinQueue := []string{}
	for _, joinMap := range c._join {
		joinStr := strings.ToUpper(joinMap["type"]) + " JOIN " +
			c.parseName(joinMap["table"]) + " " + joinMap["alias"] + " ON " +
			c.parseName(joinMap["cond"])
		joinQueue = append(joinQueue, joinStr)
	}
	if len(joinQueue) > 0 {
		table += " " + strings.Join(joinQueue, " ")
	}
	return table
}

// where 条件解析
func (c *MysqlCreator) parseWhere() string {
	where := ""
	whereQuque := []string{}
	for _, whMap := range c._cond {
		if len(whereQuque) == 0 {
			if whMap["vtype"] == "and" {
				whereQuque = append(whereQuque, whMap["where"])
			}
		} else {
			whereQuque = append(whereQuque, strings.ToUpper(whMap["type"])+"("+whMap["where"]+")")
		}
	}
	if len(whereQuque) > 0 {
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
	for _, fld := range c._fields {
		if fld != "*" {
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

func (c *MysqlCreator) parseGroup() string{
	if c._group == nil{
		c._group = []string{}
	}
	groupStr := ""
	groupQuque := []string{}
	for _, grp := range c._group{
		groupQuque = append(groupQuque, c.parseName(grp))
	}
	if len(groupQuque) > 0{
		groupStr = " GROUP BY " + strings.Join(groupQuque, ",")
	}
	return groupStr
}

func (c *MysqlCreator) parseOrder() string{
	if c._order == nil{
		c._order = []string{}
	}
	str := ""
	quque := []string{}
	for _, v := range c._order{
		v = c.parseName(v)
		vT := strings.ToLower(v)
		//println(vT, v, strings.Index(vT, " desc ") == -1 && strings.Index(vT, " asc ") == -1)
		// @Warning 此处可增强
		if strings.Index(vT, " desc") == -1 && strings.Index(vT, " asc") == -1{
			v = strings.TrimSpace(v) + " ASC "
		}
		quque = append(quque, v)
	}
	if len(quque) > 0{
		str = " ORDER BY  " + strings.Join(quque, ",")
	}
	return str
}

func (c *MysqlCreator) parseLimit() string{
	str := ""
	if c._size > 0{
		if c._page > 0{
			str = " LIMIT "+ strconv.Itoa((c._page - 1) * c._size) + "," + strconv.Itoa(c._size)
		}else {
			str = " LIMIT "+ strconv.Itoa(c._size)
		}
	}
	return str
}

func (c *MysqlCreator) Select() string {
	sqlStr := "SELECT " + c.parseField() + " FROM " +
		c.parseTable() + "" +
		c.parseWhere() +
		c.parseGroup() +
		c.parseOrder() +
		c.parseLimit()
	return sqlStr
}

func (c *MysqlCreator) Count(fields ...string) string {
	field := "*"
	sqlStr := "SELECT COUNT(" + field + ") FROM " + c.parseTable() + ""
	return sqlStr
}
