package main

import (
    "flag"
    "log"
    "net/http"
    "strings"
    "github.com/mreading/voltdb_exporter/lib"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    addresses     string
    username      string
    password      string
    namespace     string
    listenAddress string
    metricPath    string
)

// Parse CLI for flags and set variables; also acts as bootstrap --help
func init() {
    flag.StringVar(&addresses, "h", "localhost:8080", "List of cluster addresses, comma seperated")
    flag.StringVar(&username, "u", "", "Username for database authentication (required)")
    flag.StringVar(&password, "p", "", "Password for database authentication (required)")
    flag.StringVar(&namespace, "n", "voltdb", "Namespace for metrics")
    flag.StringVar(&listenAddress, "l", ":9469", "Address to listen on for web interface and telemetry.")
    flag.StringVar(&metricPath, "m", "/metrics", "Path under which to expose metrics.")
}

// Check that CLI arguments are properly set
func checkConfiguration() {
    if len(username) == 0 || len(password) == 0 {
        log.Fatal("Invalid configuration: username and password must be set. See voltdb_exporter -help for guidance")
    }
    // Put more checking here
}

func main() {
    flag.Parse()

    // split comma seperated string into list of databases
    databases := strings.Split(addresses, ",")

    checkConfiguration()

    // Configure base HTTP page with link to metrics
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`<html>
        <head><title>VoltDB exporter</title></head>
        <body>
        <h1>VoltDB exporter</h1>
        <p><a href='` + metricPath + `'>Metrics</a></p>
        </body>
        </html>
        `))
    })

    // Initialize exporter, link to Prometheus, and configure metrics HTTP page
    prometheus.MustRegister(lib.NewVoltDBExporter(namespace, username, password, databases))
    http.Handle(metricPath, promhttp.Handler())

    // Start HTTP server and prepare for scraping
    log.Printf("listening at %s", listenAddress)
    log.Fatal(http.ListenAndServe(listenAddress, nil))
}
