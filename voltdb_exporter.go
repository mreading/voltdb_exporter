package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/opsgang/prometheus_voltdb_exporter/lib"

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
	useHTTP       bool
	insecureHTTPS bool
)

// Parse CLI for flags and set variables; also acts as bootstrap --help
func init() {
	flag.StringVar(&addresses, "h", "localhost:8080", "List of cluster addresses, comma seperated")
	flag.StringVar(&username, "u", "", "Username for database authentication")
	flag.StringVar(&password, "p", "", "Password for database authentication")
	flag.StringVar(&namespace, "n", "voltdb", "Namespace for metrics")
	flag.StringVar(&listenAddress, "l", ":9469", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&metricPath, "m", "/metrics", "Path under which to expose metrics.")
	flag.BoolVar(&useHTTP, "H", false, "Use http connectins")
	flag.BoolVar(&insecureHTTPS, "i", false, "Skip certificate check")
}

// Check that CLI arguments are properly set
func checkConfiguration() {
	envVarHost, isEnvVarHostSet := os.LookupEnv("VOLTDB_EXPORTER_HOST")
	envVarUser, isEnvVarUserSet := os.LookupEnv("VOLTDB_EXPORTER_USER")
	envVarPass, isEnvVarPassSet := os.LookupEnv("VOLTDB_EXPORTER_PASS")
	envVarNamespace, isEnvVarNamespaceSet := os.LookupEnv("VOLTDB_EXPORTER_NAMESPACE")
	envVarListen, isEnvVarListenSet := os.LookupEnv("VOLTDB_EXPORTER_LISTEN")
	envVarPath, isEnvVarPathSet := os.LookupEnv("VOLTDB_EXPORTER_PATH")
	envVarProto, isEnvVarProtoSet := os.LookupEnv("VOLTDB_EXPORTER_USE_HTTP")
	_, isEnvVarInsecureSet := os.LookupEnv("VOLTDB_EXPORTER_INSECURE")

	if isEnvVarHostSet && len(envVarHost) > 0 {
		addresses = envVarHost
	}
	if isEnvVarUserSet && len(envVarUser) > 0 {
		username = envVarUser
	}
	if isEnvVarPassSet && len(envVarPass) > 0 {
		password = envVarPass
	}
	if isEnvVarNamespaceSet && len(envVarNamespace) > 0 {
		namespace = envVarNamespace
	}
	if isEnvVarListenSet && len(envVarListen) > 0 {
		listenAddress = envVarListen
	}
	if isEnvVarPathSet && len(envVarPath) > 0 {
		metricPath = envVarPath
	}
	if isEnvVarProtoSet && (envVarProto == "true" || envVarProto == "1" || envVarProto == "yes") {
		useHTTP = true
	}
	if isEnvVarInsecureSet && (envVarProto == "true" || envVarProto == "1" || envVarProto == "yes") {
		insecureHTTPS = true
	}
}

func main() {
	flag.Parse()

	checkConfiguration()

	// split comma seperated string into list of databases
	databases := strings.Split(addresses, ",")

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
	prometheus.MustRegister(lib.NewVoltDBExporter(namespace, username, password, databases, useHTTP, insecureHTTPS))
	http.Handle(metricPath, promhttp.Handler())

	// Start HTTP server and prepare for scraping
	log.Printf("listening at %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
