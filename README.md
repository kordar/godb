# godb

定义抽象接口，提供`pool`对象，开放`add`、`remove`、`get`、`exists`操作。

## `DbItem`接口

相关实现：
- [gorm-mysql](https://github.com/kordar/goframework-gorm-mysql)
- [gorm-sqlite](https://github.com/kordar/goframework-gorm-sqlite)
- [redis](https://github.com/kordar/goframework-redis)
- [leveldb](https://github.com/kordar/goframework-leveldb)

```go
type DbItem interface {
	GetName() string
	GetInstance() interface{}
	Close() error
}
```

## 使用

- 初始化

```go
mysqlpool  = godb.NewDbPool()
```

- 功能实现

```go
// 添加mysql句柄
func AddMysqlInstance(db string, cfg map[string]string) error {
    ins := NewGormConnIns(db, cfg, gormConfig())
    return mysqlpool.Add(ins)
}

func AddMysqlInstanceWithDsn(db string, dsn string) error {
    ins := NewGormConnInsWithDsn(db, dsn, gormConfig())
    return mysqlpool.Add(ins)
}

// 移除mysql句柄
func RemoveMysqlInstance(db string) {
    mysqlpool.Remove(db)
}

// mysql句柄是否存在
func HasMysqlInstance(db string) bool {
    return mysqlpool != nil && mysqlpool.Has(db)
}

// 获取mysql实例
func GetMysqlDB(db string) *gorm.DB {
    return mysqlpool.Handle(db).(*gorm.DB)
}
```
