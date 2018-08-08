package main

import (
    "flag"
    "log"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    addr          string
    user          string
    pass          string
    namespace     string
    listenAddress string
    metricPath    string
)

func init() {
    flag.StringVar(&addr, "h", "localhost:8080", "Address of cluster")
    flag.StringVar(&user, "u", "", "Username for database authentication (required)")
    flag.StringVar(&pass, "p", "", "Password for database authentication (required)")
    flag.StringVar(&namespace, "n", "voltdb", "Namespace for metrics")
    flag.StringVar(&listenAddress, "l", ":9469", "Address to listen on for web interface and telemetry.")
    flag.StringVar(&metricPath, "m", "/metrics", "Path under which to expose metrics.")
}

func checkConfiguration() {
    if len(user) == 0 || len(pass) == 0 {
        log.Fatal("Invalid configuration: username and password must be set. See voltdb_exporter -help for guidance")
    }
}

func serveLandingPage() {
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
}

func serveMetrics() {
    prometheus.MustRegister(NewVoltDBExporter())
    http.Handle(metricPath, promhttp.Handler())
}

func startHTTPServer() {
    log.Printf("listening at %s", listenAddress)
    log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func main() {
    flag.Parse()

    checkConfiguration()
    initializeClient()

    serveLandingPage()
    serveMetrics()

    startHTTPServer()
}
