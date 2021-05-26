# Prometheus Exporter for VoltDB

Exports VoltDB metrics and allows for Prometheus scraping.

## Installation

You need to have Go version go1.10.3 darwin/amd64 configured (with proper $GOPATH).

```bash
mkdir -p $GOPATH/src/github.com/mreading/
cd $GOPATH/src/github.com/mreading/
git clone https://github.com/mreading/voltdb_exporter.git
cd voltdb_exporter
go build
go install
```

## Dependencies

You need to install the following Go packages using ```go get```.

1. github.com/prometheus/client_golang/prometheus
2. github.com/tidwall/gjson

## Configuration

The exporter is configured with CLI arguments. Of course, start your VoltDB server before running the exporter (it will complain).

Flag|ENV variable|Default|Meaning
---|---|---|---
-h|VOLTDB_EXPORTER_HOST|localhost:8080|Address(es) of one or more nodes of the cluster, comma seperated
-u|VOLTDB_EXPORTER_USER|(empty)|Username for database authentication (required)
-p|VOLTDB_EXPORTER_PASS|(empty)|Password for database authentication (required)
-n|VOLTDB_EXPORTER_NAMESPACE|voltdb|Namespace for metrics
-l|VOLTDB_EXPORTER_LISTEN|:9469|Address to listen on for web interface and telemetry
-m|VOLTDB_EXPORTER_PATH|/metrics|Path under which to expose metrics
-H|VOLTDB_EXPORTER_USE_HTTP|false|Use HTTP connections instead HTTPS
-i|VOLTDB_EXPORTER_INSECURE|false|Skip certificate check

Below is an example configuration to run the exporter.

```bash
voltdb_exporter -h localhost:8080,localhost:8081 -u matt -p secret -n voltdb -l :9469 -m /metrics
```

However, you don't necessarily have to use args but you can override all those values with following environemt variables.
* VOLTDB_EXPORTER_HOST
* VOLTDB_EXPORTER_USER
* VOLTDB_EXPORTER_PASS
* VOLTDB_EXPORTER_NAMESPACE
* VOLTDB_EXPORTER_LISTEN
* VOLTDB_EXPORTER_PATH

## Prometheus

To scrape data from the VoltDB server, download and run [Prometheus](https://prometheus.io/). The default port is :9090.

```bash
prometheus --config.file=config/prometheus.yml
```

## Grafana

To visualize scraped VoltDB statistics from Prometheus, download and run [Grafana](https://grafana.com/). The default port is :3000.

```bash
brew update
brew install grafana
brew services start grafana
```

Next, import the VoltDB Dashboard (config/voltdb-grafana-dashboard.json), et voila! Enjoy your metrics.

## Notes

Ideas and code heavily inspired by other database exporters found under the [Exporter and Integrations](https://prometheus.io/docs/instrumenting/exporters/) page on the Prometheus website.

