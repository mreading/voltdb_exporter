package main

import (
    "log"

    "github.com/tidwall/gjson"
    "github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
    up                  prometheus.Gauge
    cpu_percent_used    prometheus.Gauge
    cluster_ready       prometheus.Gauge
    cluster_txns        prometheus.Gauge
    cluster_latency99   prometheus.Gauge
}

func NewVoltDBExporter() *Exporter {
    return &Exporter{
        up:                 prometheus.NewGauge(prometheus.GaugeOpts{
                                Namespace: namespace,
                                Name:      "up",
                                Help:      "Whether the VoltDB cluster is up or not",
                            }),
        cluster_ready:      prometheus.NewGauge(prometheus.GaugeOpts{
                                Namespace: namespace,
                                Name:      "cluster_ready",
                                Help:      "Whether the VoltDB cluster is running or paused",
                            }),
        cpu_percent_used:   prometheus.NewGauge(prometheus.GaugeOpts{
                                Namespace: namespace,
                                Name:      "cpu_percent_used",
                                Help:      "The percentage of total CPU available used by the database server process",
                            }),
        cluster_txns:       prometheus.NewGauge(prometheus.GaugeOpts{
                                Namespace: namespace,
                                Name:      "cluster_txns",
                                Help:      "The number of transactions per second during the measurement interval (5000ms)",
                            }),
        cluster_latency99:  prometheus.NewGauge(prometheus.GaugeOpts{
                                Namespace: namespace,
                                Name:      "cluster_latency99",
                                Help:      "The 99th percentile latency, in microseconds",
                            }),
    }
}

func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
    ch <- e.up.Desc()
    ch <- e.cluster_ready.Desc()
    ch <- e.cpu_percent_used.Desc()
    ch <- e.cluster_txns.Desc()
    ch <- e.cluster_latency99.Desc()
}

func (e *Exporter) Collect(ch chan <- prometheus.Metric) {
    log.Print("Running scrape")

    if stats, err := getStats(); err != nil {
        log.Printf("Error while getting data from VoltDB: %s", err)

        e.up.Set(0)
        ch <- e.up
    } else {
        e.up.Set(1)
        ch <- e.up

        e.cluster_ready.Set(getClusterState(stats))
        ch <- e.cluster_ready

        e.cpu_percent_used.Set(getCPUPercentUsed(stats))
        ch <- e.cpu_percent_used

        e.cluster_txns.Set(getClusterTxns(stats))
        ch <- e.cluster_txns

        e.cluster_latency99.Set(getClusterLatency99(stats))
        ch <- e.cluster_latency99
    }

    log.Print("Scrape successful")
}

func getClusterState(stats *Stats) float64 {
    var json string
    gjson.Unmarshal(stats.state, &json)
    value := gjson.Get(json, "results.0.data.19.2")
    if value.Str == "RUNNING" {
        return float64(1)
    }
    return float64(0)
}

func getCPUPercentUsed(stats *Stats) float64 {
    var json string
    gjson.Unmarshal(stats.cpu, &json)
    value := gjson.Get(json, "results.0.data.0.3")

    return value.Num
}

func getClusterTxns(stats *Stats) float64 {
    var json string
    gjson.Unmarshal(stats.txns, &json)
    value := gjson.Get(json, "results.0.data.0.5")

    return value.Num
}

func getClusterLatency99(stats *Stats) float64 {
    var json string
    gjson.Unmarshal(stats.latency, &json)
    value := gjson.Get(json, "results.0.data.0.8")

    return value.Num
}
