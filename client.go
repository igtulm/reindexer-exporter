package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
)

var client = &http.Client{Timeout: 15 * time.Second}

func initClient() {
	client = &http.Client{
		Timeout: time.Duration(config.RequestTimeout) * time.Second,
	}
}

func apiRequest(req_url string) ([]byte, error) {
	req, err := http.NewRequest("GET", req_url, nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "request": req_url}).Error("Error while requiesting")
		return nil, errors.New("Error while constructing request")
	}

	req.SetBasicAuth(config.ReindexerUsername, config.ReindexerPassword)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil || resp == nil || resp.StatusCode != 200 {
		status := 0
		if resp != nil {
			status = resp.StatusCode
		}
		log.WithFields(log.Fields{"error": err, "request": req_url, "statusCode": status}).Error("Error while retrieving data")
		return nil, errors.New("Error while retrieving data")
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"body": string(body), "request": req_url}).Debug("Metrics loaded")

	return body, err
}

func apiGetNamespacesList(config reindexerExporterConfig) ([]byte, error) {
	req_url := config.ReindexerURL + "/api/v1/db/" + config.ReindexerDBName + "/namespaces"

	return apiRequest(req_url)
}

func apiGetQuery(config reindexerExporterConfig, endpoint string) ([]byte, error) {
	req_url := config.ReindexerURL + "/api/v1/db/" + config.ReindexerDBName + "/query?q=" + url.QueryEscape("SELECT * FROM "+endpoint)

	return apiRequest(req_url)
}
