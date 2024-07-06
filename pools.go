package godb

import (
	"errors"
	"github.com/kordar/gologger"
	"sync"
)

var instanceOfDbPool *DbConnPool
var onceOfDb sync.Once

type DbItem interface {
	GetName() string
	GetInstance() interface{}
	Close() error
}

type DbConnPool struct {
	handle map[string]DbItem
}

func GetDbPool() *DbConnPool {
	onceOfDb.Do(func() {
		instanceOfDbPool = &DbConnPool{handle: make(map[string]DbItem)}
	})
	return instanceOfDbPool
}

func NewDbPool() *DbConnPool {
	return &DbConnPool{handle: make(map[string]DbItem)}
}

// InitDataPool /*
// 初始化数据库连接(可在mail()适当位置调用)
func (m *DbConnPool) InitDataPool(items ...DbItem) (issucc bool) {
	for _, item := range items {
		if m.handle[item.GetName()] != nil {
			logger.Errorf("[godb] the db-%s already exists", item.GetName())
			continue
		}
		var err error
		err = m.Add(item)
		if err != nil {
			logger.Fatal(err)
			return false
		}
	}

	// 关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

// Add 添加数据库实例
func (m *DbConnPool) Add(db DbItem) error {
	if m.handle[db.GetName()] != nil {
		return errors.New("[godb] the db already exists")
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
			logger.Errorf("[godb] remove db err，%v", err)
		}
	}
}

// Handle /*
// 对外获取数据库连接对象db
func (m *DbConnPool) Handle(name string) (conn interface{}) {
	if m.Has(name) {
		return m.handle[name].GetInstance()
	}
	return nil
}

// Has 是否存在
func (m *DbConnPool) Has(name string) bool {
	item := m.handle[name]
	return item != nil
}
