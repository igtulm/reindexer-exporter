# Reindexer Exporter

Prometheus exporter for [Reindexer](https://github.com/Restream/reindexer) metrics.
Data is scraped by [Prometheus](https://prometheus.io).

Metrics are exposed through [http://localhost:9451/metrics](http://localhost:9451/metrics).


## How to start
- [Build and run the exporter](#build-and-run-the-exporter)
- [Prometheus configuration](#prometheus-configuration)
- Run Prometheus

## Build and run the exporter

```bash
dep ensure
make
cmd/reindexer_exporter
```

## Prometheus configuration

```yml
# my global config
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'
    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.
    static_configs:
    - targets: ['localhost:9090']

  # Add this job rule for Reindexer metrics scraping
  - job_name: 'reindexer'
    static_configs:
      - targets:
        - localhost:9451
```
