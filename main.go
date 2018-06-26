package main

import (
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

const defaultLogLevel = log.ErrorLevel

func initLogger() {
	log.SetLevel(getLogLevel())

	if strings.ToUpper(config.LogFormat) == "JSON" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
}

func main() {
	initConfig()
	initLogger()
	initClient()

	exporter := newExporter()
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Reindexer Exporter</title></head>
			<body>
			<h1>Reindexer Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Starting Reindexer exporter")
	log.Fatal(http.ListenAndServe(config.PublishAddr+":"+config.PublishPort, nil))
}

func getLogLevel() log.Level {
	lvl := strings.ToLower(os.Getenv("LOG_LEVEL"))
	level, err := log.ParseLevel(lvl)
	if err != nil {
		level = defaultLogLevel
	}
	return level
}
