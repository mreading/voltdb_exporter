package lib

import (
    "github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
    client              *VoltDBClient
    databases           []string

    up                  prometheus.Gauge

    cpu_percent_used    *prometheus.GaugeVec
    cluster_ready       *prometheus.GaugeVec
    cluster_txns        *prometheus.GaugeVec
    cluster_latency99   *prometheus.GaugeVec
    server_ram_used     *prometheus.GaugeVec
    dr_role             *prometheus.GaugeVec
    dr_state            *prometheus.GaugeVec
}

func NewVoltDBExporter(namespace string, user string, pass string, dbs []string) *Exporter {
    return &Exporter{
        client:    NewVoltDBClient(user, pass, dbs),
        databases: dbs,

        up: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "up",
                Help:      "Whether the VoltDB cluster is up or not",
            }),

        cluster_ready: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "cluster_ready",
                Help:      "Whether the VoltDB cluster is running or paused",
            }, 
            []string{"database"}),

        cpu_percent_used: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "cpu_percent_used",
                Help:      "The percentage of total CPU available used by the database server process",
            }, 
            []string{"database"}),

        cluster_txns: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "cluster_txns",
                Help:      "The number of transactions per second during the measurement interval (5000ms)",
            }, 
            []string{"database"}),

        cluster_latency99: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "cluster_latency99",
                Help:      "The 99th percentile latency, in microseconds",
            }, 
            []string{"database"}),

        server_ram_used: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "server_ram_used",
                Help:      "The amount of memory (in kilobytes) allocated by Java and current in use by VoltDB",
            }, 
            []string{"database"}),

        dr_role: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "dr_role",
                Help:      "Database replication role",
            }, 
            []string{"database"}),

        dr_state: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "dr_state",
                Help:      "Database replication state",
            }, 
            []string{"database"}),
    }
}
