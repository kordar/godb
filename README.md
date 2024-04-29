# gorm集成

## 实现gorm实例对象

```go
type GormConnIns struct {
	name        string
	ins         *gorm.DB
	mysqlConfig *gorm.Config
}

func NewGormConnIns(name string, config *gorm.Config) *GormConnIns {
	return &GormConnIns{name: name, mysqlConfig: config}
}

func (c GormConnIns) GetName() string {
	return c.name
}

func (c GormConnIns) GetInstance() interface{} {
	user := gocfg.GetSectionValue(c.name, "user")
	password := gocfg.GetSectionValue(c.name, "password")
	host := gocfg.GetSectionValue(c.name, "host")
	port := gocfg.GetSectionValue(c.name, "port")
	database := gocfg.GetSectionValue(c.name, "db")
	charset := gocfg.GetSectionValue(c.name, "charset")
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, database, "charset="+charset+"&parseTime=true")
	if obj, err := gorm.Open(mysql.Open(source), c.mysqlConfig); err == nil {
		c.ins = obj
		return obj
	} else {
		return err
	}
}

func (c GormConnIns) Close() error {
	if s, err := c.ins.DB(); err == nil {
		return s.Close()
	} else {
		return err
	}
}
```

## 初始化gorm驱动，例如mysql

```go
var mysqlpool *godb.DbConnPool

// ---------------- mysql -----------------------------

func gormConfig() *gorm.Config {
dbLogLevel := gocfg.GetSystemValue("gorm_log_level")
mysqlConfig := gorm.Config{}
if dbLogLevel == "error" {
mysqlConfig.Logger = logger.Default.LogMode(logger.Error)
}
if dbLogLevel == "warn" {
mysqlConfig.Logger = logger.Default.LogMode(logger.Warn)
}
if dbLogLevel == "info" {
mysqlConfig.Logger = logger.Default.LogMode(logger.Info)
}
return &mysqlConfig
}

// InitMysqlHandle 初始化mysql句柄
func InitMysqlHandle(dbs ...string) {
mysqlpool = godb.GetDbPool()
for _, s := range dbs {
ins := NewGormConnIns(s, gormConfig())
err := mysqlpool.Add(ins)
if err != nil {
log.Warnf("初始化异常，err=%v", err)
}
}

}

// AddMysqlInstance 添加mysql句柄
func AddMysqlInstance(db string) error {
mysqlpool = godb.GetDbPool()
ins := NewGormConnIns(db, gormConfig())
return mysqlpool.Add(ins)
}

// RemoveMysqlInstance 移除mysql句柄
func RemoveMysqlInstance(db string) {
mysqlpool.Remove(db)
}
```


## 初始化

```go
common.InitMysqlHandle("database")
```

