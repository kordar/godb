package godb_test

import (
	"bytes"
	"github.com/kordar/godb"
	"log/slog"
	"testing"
)

type TestItem struct {
}

func (t TestItem) GetName() string {
	return "test"
}

func (t TestItem) GetInstance() interface{} {
	return "AAAA"
}

func (t TestItem) Close() error {
	return nil
}

func TestName(t *testing.T) {
	var buf bytes.Buffer
	godb.SetLogger(slog.New(slog.NewTextHandler(&buf, nil)))

	mysqlpool := godb.NewDbPool()
	if err := mysqlpool.Add(TestItem{}); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	conn := mysqlpool.Handle("test")
	if conn == nil {
		t.Fatalf("expected conn not nil")
	}
	if conn.(string) != "AAAA" {
		t.Fatalf("unexpected conn: %v", conn)
	}
}
