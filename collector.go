package main

import (
    "strings"
    "github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
    addrs        []string
    user         string
    pass         string
    namespace    string
}

func NewVoltDBExporter(addr string, user string, pass string, namespace string) *Exporter {
    return &Exporter{
        addrs:       strings.Split(addr, ","),
        user:        user,
        pass:        pass,
        namespace:   namespace,
    }
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	// TODO
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// TODO
}