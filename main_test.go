package main

import (
	"reflect"
	"testing"
)

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

func TestLoadConfigContent_ParsingError(t *testing.T) {
	configFile := `
DUMMY
:FOO
`
	_, err := loadConfigContent([]byte(configFile))
	if err == nil {
		t.Errorf("Should have an error parsing unparseable content")
	}
}

func TestLoadConfigContent_Ok_Standard(t *testing.T) {
	configFile := `
resource_tags:
  - tag_selections:
    - tag_name: "Client"
      tag_value: "Alice"
    - tag_name: "Env"
      tag_value: "Prod"
`
	want := Config{
		[]ResourceTag{
			ResourceTag{
				TagSelections: []TagSelection{
					TagSelection{
						TagName:  "Client",
						TagValue: "Alice",
					},
					TagSelection{
						TagName:  "Env",
						TagValue: "Prod",
					},
				},
			},
		},
	}
	got, err := loadConfigContent([]byte(configFile))
	if err != nil {
		t.Errorf("Error on loading config content %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error in getting config Got:%v, Expected config:%v", got, want)
	}
}
