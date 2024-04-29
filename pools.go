package godb

import (
	"errors"
	log "github.com/kordar/gologger"
	"sync"
)

var instanceOfDbPool *DbConnPool
var onceOfDb sync.Once

type DBItem interface {
	GetName() string
	GetInstance() interface{}
	Close() error
}

type DbConnPool struct {
	handle map[string]DBItem
}

func GetDbPool() *DbConnPool {
	onceOfDb.Do(func() {
		instanceOfDbPool = &DbConnPool{handle: make(map[string]DBItem)}
	})
	return instanceOfDbPool
}

// InitDataPool /*
// 初始化数据库连接(可在mail()适当位置调用)
func (m *DbConnPool) InitDataPool(items ...DBItem) (issucc bool) {
	for _, item := range items {
		if m.handle[item.GetName()] != nil {
			log.Errorf("[godb] 实例[%s]已存在", item.GetName())
			continue
		}
		var err error
		err = m.Add(item)
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	// 关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

// Add 添加数据库实例
func (m *DbConnPool) Add(db DBItem) error {
	if m.handle[db.GetName()] != nil {
		return errors.New("MySQL实例已存在")
	}
	m.handle[db.GetName()] = db
	return nil
}

// Remove 移除句柄
func (m *DbConnPool) Remove(name string) {
	if m.handle[name] != nil {
		defer delete(m.handle, name)
		g := m.handle[name]
		if err := g.Close(); err != nil {
			log.Errorf("[godb] 移除句柄=%v", err)
		}
	}
}

// Handle /*
// 对外获取数据库连接对象db
func (m *DbConnPool) Handle(name string) (conn interface{}) {
	return m.handle[name].GetInstance()
}
