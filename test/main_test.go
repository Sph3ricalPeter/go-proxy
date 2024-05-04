package main

import (
	"reflect"
	"testing"
)

func TestMain(t *testing.T) {
	configPath := "test_config.json"

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := Config{
		Name:    "go-proxy",
		Host:    "localhost",
		Port:    8080,
		Targets: map[string]string{"/somepath": "someurl"},
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("unexpected config: got %+v, want %+v", config, expected)
	}
}
