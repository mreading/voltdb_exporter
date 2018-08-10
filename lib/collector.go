package lib

import (
    "log"
    "github.com/tidwall/gjson"
    "github.com/prometheus/client_golang/prometheus"
)

// Writes all descriptors to the Prometheus desc channel
func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
    ch <- e.up.Desc()
    e.cluster_ready.Describe(ch)
    e.cpu_percent_used.Describe(ch)
    e.cluster_txns.Describe(ch)
    e.cluster_latency99.Describe(ch)
    e.server_ram_used.Describe(ch)
    e.dr_role.Describe(ch)
    e.dr_state.Describe(ch)
}

// Writes latest value to Prometheus metric channel with either
// ch <- e.<metric> or <metric>.Collect(ch) 
func (e *Exporter) Collect(ch chan <- prometheus.Metric) {
    log.Print("Running scrape")

    if stats, err := getStats(e.databases, e.client); err != nil {
        log.Printf("Error while getting data from VoltDB: %s", err)

        e.up.Set(0)
        ch <- e.up
    } else {
        e.up.Set(1)
        ch <- e.up

        collectPerDatabaseGauge(stats, e.cluster_ready, getClusterState, ch)
        collectPerDatabaseGauge(stats, e.cpu_percent_used, getCPUPercentUsed, ch)
        collectPerDatabaseGauge(stats, e.cluster_txns, getClusterTxns, ch)
        collectPerDatabaseGauge(stats, e.cluster_latency99, getClusterLatency99, ch)
        collectPerDatabaseGauge(stats, e.server_ram_used, getServerRAMUsed, ch)
        collectPerDatabaseGauge(stats, e.dr_role, getDrRole, ch)
        collectPerDatabaseGauge(stats, e.dr_state, getDrState, ch)
    }

    log.Print("Scrape complete")
}

func collectPerDatabaseGauge(s *[]Stats, vec *prometheus.GaugeVec, collectFunc func(Stats) float64, ch chan<- prometheus.Metric) {
    for _, st := range *s {
        vec.WithLabelValues(st.database).Set(collectFunc(st))
    }
    vec.Collect(ch)
}

func getClusterState(stats Stats) float64 {
    var json string
    gjson.Unmarshal(stats.state, &json)
    value := gjson.Get(json, "results.0.data.19.2")
    if value.Str == "RUNNING" {
        return float64(1)
    }
    return float64(0) // PAUSED
}

func getCPUPercentUsed(stats Stats) float64 {
    var json string
    gjson.Unmarshal(stats.cpu, &json)
    value := gjson.Get(json, "results.0.data.0.3")

    return value.Num
}

func getClusterTxns(stats Stats) float64 {
    var json string
    gjson.Unmarshal(stats.txns, &json)
    value := gjson.Get(json, "results.0.data.0.5")

    return value.Num
}

func getClusterLatency99(stats Stats) float64 {
    var json string
    gjson.Unmarshal(stats.latency, &json)
    value := gjson.Get(json, "results.0.data.0.8")

    return value.Num/1000 // Converted to ms
}

func getServerRAMUsed(stats Stats) float64 {
    var json string
    gjson.Unmarshal(stats.ram, &json)
    value := gjson.Get(json, "results.0.data.0.3")

    return value.Num/1000000 // Converted to GB
}

func getDrRole(stats Stats) float64 {
    var json string
    gjson.Unmarshal(stats.dr_role, &json)
    value := gjson.Get(json, "results.0.data.0.0")
    if value.Str == "MASTER" {
        return float64(1)
    } else if value.Str == "REPLICA" {
        return float64(2)
    } else if value.Str == "XDCR" {
        return float64(3)
    }
    return float64(0) // NONE
}

func getDrState(stats Stats) float64 {
    var json string
    gjson.Unmarshal(stats.dr_state, &json)
    value := gjson.Get(json, "results.0.data.0.1")
    if value.Str == "ACTIVE" {
        return float64(1)
    } else if value.Str == "PENDING" {
        return float64(2)
    } else if value.Str == "STOPPED" {
        return float64(3)
    }
    return float64(0) // DISABLED
}
