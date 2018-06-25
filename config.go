package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	config        reindexerExporterConfig
	defaultConfig = reindexerExporterConfig{
		PublishAddr:       "",
		PublishPort:       "9451",
		LogFormat:         "JSON",
		RequestTimeout:    30,
		ReindexerURL:      "http://api.staging.itv.restr.im:9088",
		ReindexerUsername: "owner",
		ReindexerPassword: "123456",
		ReindexerDBName:   "itv_api_ng",
		EnabledExporters:  []string{"memstats"}, //, "perfstats", "queriesperfstats"},
	}
)

type reindexerExporterConfig struct {
	PublishAddr       string
	PublishPort       string
	LogFormat         string
	RequestTimeout    int
	ReindexerURL      string
	ReindexerUsername string
	ReindexerPassword string
	ReindexerDBName   string
	EnabledExporters  []string
}

func initConfig() {
	config = defaultConfig

	if addr := os.Getenv("PUBLISH_ADDR"); addr != "" {
		config.PublishAddr = addr
	}

	if port := os.Getenv("PUBLISH_PORT"); port != "" {
		if _, err := strconv.Atoi(port); err == nil {
			config.PublishPort = port
		} else {
			panic(fmt.Errorf("The configured port is not a valid number: %v", port))
		}

	}

	if logformat := os.Getenv("LOG_FORMAT"); logformat != "" {
		config.LogFormat = logformat
	}

	if timeout := os.Getenv("REINDEXER_REQUEST_TIMEOUT"); timeout != "" {
		t, err := strconv.Atoi(timeout)
		if err != nil {
			panic(fmt.Errorf("timeout is not a number: %v", err))
		}
		config.RequestTimeout = t
	}

	if url := os.Getenv("REINDEXER_URL"); url != "" {
		if valid, _ := regexp.MatchString("http?://[a-zA-Z.0-9]+", strings.ToLower(url)); valid {
			config.ReindexerURL = url
		} else {
			panic(fmt.Errorf("URL must start with http://"))
		}
	}

	user := os.Getenv("REINDEXER_USER")
	if user != "" {
		config.ReindexerUsername = user
	}

	pass := os.Getenv("REINDEXER_PASSWORD")
	if pass != "" {
		config.ReindexerPassword = pass
	}

	db := os.Getenv("REINDEXER_DB")
	if db != "" {
		config.ReindexerDBName = db
	}

	if enabledExporters := os.Getenv("REINDEXER_EXPORTERS"); enabledExporters != "" {
		config.EnabledExporters = strings.Split(enabledExporters, ",")
	}
}
