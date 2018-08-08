# Prometheus Exporter for VoltDB

Exports VoltDB metrics and allows for Prometheus scraping.

## Installation

You need to have Go version go1.10.3 darwin/amd64 configured (with proper $GOPATH).

```bash
mkdir -p $GOPATH/src/github.com/user/voltdb_exporter
cd $GOPATH/src/github.com/user/voltdb_exporter
git clone https://github.com/mreading/voltdb_exporter.git
go build
go install
```

## Configuration

The exporter is configured with command line arguments. Of course, start your VoltDB server before running the exporter (it will complain).

Flag|ENV variable|Default|Meaning
---|---|---|---
-h|DB.ADDRESS(ES)|localhost:8080|Address(es) of one or more nodes of the cluster, comma seperated
-u|DB.USERNAME|(empty)|Username for database authentication (required)
-p|DB.PASSWORD|(empty)|Password for database authentication (required)
-n|NAMESPACE|voltdb|Namespace for metrics
-l|LISTENADDRESS|:9469|Address to listen on for web interface and telemetry
-m|METRICPATH|/metrics|Path under which to expose metrics

Below is an example configuration to run the exporter.

```bash
voltdb_exporter -h localhost:8080 -u matt -p matt -n voltdb -l :9469 -m /metrics
```

## Prometheus

To scrape data from the VoltDB server, download and run [Prometheus](https://prometheus.io/). The default port is :9090.

```bash
prometheus --config.file=prometheus.yml
```

## Grafana

To visualize scraped VoltDB statistics from Prometheus, download and run [Grafana](https://grafana.com/). The default port is :3000.

```bash
brew update
brew install grafana
brew services start grafana
```

Next, import the VoltDB Dashboard (voltdb-grafana-dashboard.json), et voila!