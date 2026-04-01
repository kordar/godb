# godb

轻量数据库句柄池，提供线程安全的句柄注册/获取/移除能力，日志基于标准库 `log/slog`。

## 特性

- 线程安全：`RWMutex` 保护句柄表
- 标准日志：移除第三方日志依赖，使用 `slog.Logger`
- 错误友好：支持返回错误的初始化与移除接口
- 零侵入：只要求业务句柄实现统一 `DbItem` 接口

## 安装

```bash
go get github.com/kordar/godb
```

## 核心接口

```go
type DbItem interface {
	GetName() string
	GetInstance() interface{}
	Close() error
}
```

## 快速开始

```go
package main

import (
	"log/slog"
	"os"

	"github.com/kordar/godb"
)

type MysqlItem struct {
	name string
	db   any
}

func (m MysqlItem) GetName() string        { return m.name }
func (m MysqlItem) GetInstance() interface{} { return m.db }
func (m MysqlItem) Close() error           { return nil }

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	godb.SetLogger(logger)

	pool := godb.NewDbPool()
	_ = pool.InitDataPoolE(
		MysqlItem{name: "main", db: "gorm-db-instance"},
	)

	conn := pool.Handle("main")
	_ = conn

	_ = pool.RemoveE("main")
}
```

## API 说明

- `NewDbPool() *DbConnPool`：创建连接池
- `SetLogger(*slog.Logger)`：设置库内日志器
- `InitDataPool(items ...DbItem) bool`：兼容旧接口，失败返回 false
- `InitDataPoolE(items ...DbItem) error`：推荐，返回明确错误
- `Add(db DbItem) error`：添加句柄
- `Handle(name string) interface{}`：获取句柄实例，不存在返回 nil
- `Has(name string) bool`：句柄是否存在
- `Remove(name string)`：移除句柄并记录错误日志
- `RemoveE(name string) error`：移除句柄并返回错误
- `Item(name string) (DbItem, bool)`：获取原始 `DbItem`
