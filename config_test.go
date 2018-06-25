package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestEnvironmentSettingAddr(t *testing.T) {
	newValue := "localhost"
	os.Setenv("PUBLISH_ADDR", newValue)
	defer os.Unsetenv("PUBLISH_ADDR")
	initConfig()
	if config.PublishAddr != newValue {
		t.Errorf("Expected config.PUBLISH_ADDR to be modified. Found=%v, expected=%v", config.PublishAddr, newValue)
	}
}

func TestEnvironmentSettingPort(t *testing.T) {
	newValue := "1234"
	os.Setenv("PUBLISH_PORT", newValue)
	defer os.Unsetenv("PUBLISH_PORT")
	initConfig()
	if config.PublishPort != newValue {
		t.Errorf("Expected config.PUBLISH_PORT to be modified. Found=%v, expected=%v", config.PublishPort, newValue)
	}
}

func TestEnvironmentSettingFormat(t *testing.T) {
	newValue := "json"
	os.Setenv("LOG_FORMAT", newValue)
	defer os.Unsetenv("LOG_FORMAT")
	initConfig()
	if config.LogFormat != newValue {
		t.Errorf("Expected config.OUTPUT_FORMAT to be modified. Found=%v, expected=%v", config.LogFormat, newValue)
	}
}

func TestEnvironmentSettingRequestTimeout(t *testing.T) {
	newValue := "60"
	os.Setenv("REINDEXER_REQUEST_TIMEOUT", newValue)
	defer os.Unsetenv("REINDEXER_REQUEST_TIMEOUT")
	initConfig()
	val, err := strconv.Atoi(newValue)
	if err != nil {
		t.Errorf("Expected config.REINDEXER_REQUEST_TIMEOUT is a nubmer")
	}
	if config.RequestTimeout != val {
		t.Errorf("Expected config.REINDEXER_REQUEST_TIMEOUT to be modified. Found=%v, expected=%v", config.RequestTimeout, val)
	}
}

func TestEnvironmentSettingURL_HTTP(t *testing.T) {
	newValue := "http://testURL"
	os.Setenv("REINDEXER_URL", newValue)
	defer os.Unsetenv("REINDEXER_URL")
	initConfig()
	if config.ReindexerURL != newValue {
		t.Errorf("Expected config.REINDEXER_URL to be modified. Found=%v, expected=%v", config.ReindexerURL, newValue)
	}
}

func TestEnvironmentSettingUser(t *testing.T) {
	newValue := "username"
	os.Setenv("REINDEXER_USER", newValue)
	defer os.Unsetenv("REINDEXER_USER")
	initConfig()
	if config.ReindexerUsername != newValue {
		t.Errorf("Expected config.REINDEXER_USER to be modified. Found=%v, expected=%v", config.ReindexerUsername, newValue)
	}
}

func TestEnvironmentSettingPassword(t *testing.T) {
	newValue := "password"
	os.Setenv("REINDEXER_PASSWORD", newValue)
	defer os.Unsetenv("REINDEXER_PASSWORD")
	initConfig()
	if config.ReindexerPassword != newValue {
		t.Errorf("Expected config.REINDEXER_PASSWORD to be modified. Found=%v, expected=%v", config.ReindexerPassword, newValue)
	}
}

func TestEnvironmentSettingDbName(t *testing.T) {
	newValue := "password"
	os.Setenv("REINDEXER_DB", newValue)
	defer os.Unsetenv("REINDEXER_DB")
	initConfig()
	if config.ReindexerDBName != newValue {
		t.Errorf("Expected config.REINDEXER_DB to be modified. Found=%v, expected=%v", config.ReindexerDBName, newValue)
	}
}

func TestConfig_Port(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("initConfig should panic on invalid port config")
		}
	}()
	port := config.PublishPort
	os.Setenv("PUBLISH_PORT", "noNumber")
	defer os.Unsetenv("PUBLISH_PORT")
	initConfig()
	if config.PublishPort != port {
		t.Errorf("Invalid Portnumber. It should not be set. expected=%v,got=%v", port, config.PublishPort)
	}
}

func TestConfig_Addr(t *testing.T) {
	addr := config.PublishAddr
	os.Setenv("PUBLISH_ADDR", "")
	defer os.Unsetenv("PUBLISH_ADDR")
	initConfig()
	if config.PublishAddr != addr {
		t.Errorf("Invalid Addrress. It should not be set. expected=%v,got=%v", addr, config.PublishAddr)
	}
}

func TestConfig_RequestTimeout(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("initConfig should panic on invalid request timeout")
		}
	}()
	value := config.RequestTimeout
	os.Setenv("REINDEXER_REQUEST_TIMEOUT", "noNumber")
	defer os.Unsetenv("REINDEXER_REQUEST_TIMEOUT")
	initConfig()
	if config.RequestTimeout != value {
		t.Errorf("Invalid request timeout. It should not be set. expected=%v,got=%v", value, config.RequestTimeout)
	}
}

func TestConfig_Http_URL(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("initConfig should panic on invalid url config")
		}
	}()
	url := config.ReindexerURL
	os.Setenv("REINDEXER_URL", "ftp://test")
	defer os.Unsetenv("REINDEXER_URL")
	initConfig()
	if config.ReindexerURL != url {
		t.Errorf("Invalid URL. It should start with http://. expected=%v,got=%v", url, config.ReindexerURL)
	}
}

func TestConfig_EnabledExporters(t *testing.T) {
	enabledExporters := []string{"perfstats", "memstats"}
	os.Setenv("REINDEXER_EXPORTERS", "perfstats,memstats")
	defer os.Unsetenv("REINDEXER_EXPORTERS")
	initConfig()
	if diff := pretty.Compare(config.EnabledExporters, enabledExporters); diff != "" {
		t.Errorf("Invalid Exporters list. diff\n%v", diff)
	}
}
