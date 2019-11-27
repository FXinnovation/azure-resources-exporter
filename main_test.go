package main

import "testing"

func TestLoadConfig_No_Config(t *testing.T) {
	_, err := loadConfig("no_config")
	if err != nil {
		t.Errorf("Error on loading config %v", err)
	}
}

func TestLoadConfig_Example_Config(t *testing.T) {
	_, err := loadConfig("config/config_example.yml")
	if err != nil {
		t.Errorf("Error on loading config %v", err)
	}
}
