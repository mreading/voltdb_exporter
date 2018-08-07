package main

import (
    "log"

    "github.com/tidwall/gjson"
    "github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
    up                  prometheus.Gauge
    cpu_percent_used    prometheus.Gauge
}

func NewVoltDBExporter() *Exporter {
    return &Exporter{
        up:                 prometheus.NewGauge(prometheus.GaugeOpts{
                                Namespace: namespace,
                                Subsystem: "",
                                Name:      "up",
                                Help:      "Whether the VoltDB scrape was successful",
                            }),
        cpu_percent_used:   prometheus.NewGauge(prometheus.GaugeOpts{
                                Namespace: namespace,
                                Subsystem: "",
                                Name:      "cpu_percent_used",
                                Help:      "The percentage of total CPU available used by the database server process",
                            }),
    }
}

func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
    ch <- e.up.Desc()
    ch <- e.cpu_percent_used.Desc()
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

        e.cpu_percent_used.Set(getCPUPercentUsed(stats))
        ch <- e.cpu_percent_used
    }

    log.Print("Scrape successful")
}

func getCPUPercentUsed(stats *Stats) float64 {
    var json string
    gjson.Unmarshal(stats.cpu, &json)
    value := gjson.Get(json, "results.0.data.0.3")
    
    return value.Num
}
