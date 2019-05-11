package config

import (
	"testing"
)

func TestConfigManage(t *testing.T) {
	c := NewConfigManage()
	logger := map[string]interface{}{}
	arraylogger := map[string]interface{}{}
	c.Register("Logger", nil, logger, func(config *Config) error {
		t.Log(logger)
		return nil
	})
	c.Register("TestArray", nil, &arraylogger, func(config *Config) error {
		t.Log(logger)
		return nil
	})
	c.Init("test.yaml")
	if logger["Level"] != "ERROR" {
		t.Fatalf("level is no error %#v", logger)
	}
	//t.Fatal(logger)
}
