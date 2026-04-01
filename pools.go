package godb

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

var poolLogger = slog.Default()

type DbItem interface {
	GetName() string
	GetInstance() interface{}
	Close() error
}

type DbConnPool struct {
	handle map[string]DbItem
	locker sync.RWMutex
}

func NewDbPool() *DbConnPool {
	return &DbConnPool{
		handle: make(map[string]DbItem),
		locker: sync.RWMutex{},
	}
}

func SetLogger(logger *slog.Logger) {
	if logger != nil {
		poolLogger = logger
	}
}

func (m *DbConnPool) InitDataPool(items ...DbItem) (issucc bool) {
	if err := m.InitDataPoolE(items...); err != nil {
		poolLogger.Error("[godb] init data pool failed", slog.Any("err", err))
		return false
	}
	return true
}

func (m *DbConnPool) InitDataPoolE(items ...DbItem) error {
	for _, item := range items {
		if item == nil {
			continue
		}
		if _, ok := m.Item(item.GetName()); ok {
			poolLogger.Warn("[godb] db already exists", slog.String("name", item.GetName()))
			continue
		}
		if err := m.Add(item); err != nil {
			return err
		}
	}
	return nil
}

func (m *DbConnPool) Add(db DbItem) error {
	if db == nil {
		return errors.New("[godb] db item is nil")
	}
	name := db.GetName()
	if name == "" {
		return errors.New("[godb] db name is empty")
	}

	m.locker.Lock()
	defer m.locker.Unlock()
	if m.handle[name] != nil {
		return errors.New("[godb] the db already exists")
	}
	m.handle[name] = db
	return nil
}

func (m *DbConnPool) Remove(name string) {
	m.locker.Lock()
	item := m.handle[name]
	delete(m.handle, name)
	m.locker.Unlock()

	if item == nil {
		return
	}
	if err := item.Close(); err != nil {
		poolLogger.Error("[godb] remove db failed", slog.String("name", name), slog.Any("err", err))
	}
}

func (m *DbConnPool) RemoveE(name string) error {
	m.locker.Lock()
	item := m.handle[name]
	delete(m.handle, name)
	m.locker.Unlock()

	if item == nil {
		return nil
	}
	if err := item.Close(); err != nil {
		return fmt.Errorf("[godb] remove db %s failed: %w", name, err)
	}
	return nil
}

func (m *DbConnPool) Handle(name string) (conn interface{}) {
	item, exists := m.Item(name)
	if !exists {
		return nil
	}
	return item.GetInstance()
}

func (m *DbConnPool) Has(name string) bool {
	_, ok := m.Item(name)
	return ok
}

func (m *DbConnPool) Item(name string) (DbItem, bool) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	item, ok := m.handle[name]
	if !ok || item == nil {
		return nil, false
	}
	return item, true
}
