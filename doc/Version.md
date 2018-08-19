# Version

## 0.0.x

> 实现 *sql* 字符串的生成



### 0.02/20180819

> sqler

- (优化) *Creator* 添加方法 *Group*/*Order* 接口
- (优化) *MysqlCreator*
  - 内部变量更改为 *_name* 开头 
  - (实现) 实现 *Group*/*Order*/*Page* 解析等方法

### 0.0.1/20180818

> sqler

- 创建包 *sqler* 
- 实现 *Creator* 接口，并提供根据*driver*不同的 数据库适配器
- 实现 *MySQL* 数据库生成器
  - 添加*mysql_test* 测试用例