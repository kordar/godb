package godb_test

import (
	"github.com/kordar/godb"
	logger "github.com/kordar/gologger"
	"testing"
)

type TestItem struct {
}

func (t TestItem) GetName() string {
	//TODO implement me
	return "test"
}

func (t TestItem) GetInstance() interface{} {
	//TODO implement me
	return "AAAA"
}

func (t TestItem) Close() error {
	//TODO implement me
	return nil
}

func TestName(t *testing.T) {
	mysqlpool := godb.NewDbPool()
	mysqlpool.Add(TestItem{})
	conn := mysqlpool.Handle("test")
	logger.Infof("-------%v", conn)
}
