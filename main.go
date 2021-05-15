package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
)

var (
	v2rayAddr       string
	listenAddr      string
	metricsEndpoint string
	debugMode       bool
)

func main() {
	pflag.StringVar(&v2rayAddr, "target", "127.0.0.1:10150", "v2ray grpc api endpoint")
	pflag.StringVar(&listenAddr, "listen", "127.0.0.1:9100", "address exporter to listen")
	pflag.StringVar(&metricsEndpoint, "endpoint", "/metrics", "enpoint for metrics")
	pflag.BoolVar(&debugMode, "debug", false, "print debug log")
	pflag.Parse()
	initLogger()

	signals := make(chan os.Signal, 1)
	v2c, err := NewClient(v2rayAddr)
	if err != nil {
		logger.Fatalf("dial V2Ray gRPC server: %v", err)
	}
	prometheus.MustRegister(NewExporter(v2c))
	defer v2c.Close()
	go serveHTTP(listenAddr, metricsEndpoint)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	logger.Warn("v2ray_exporter exit")
}

func serveHTTP(listenAddress, metricsEndpoint string) {
	http.Handle(metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>v2ray exporter</title></head>
			<body>
			<h1>v2ray exporter</h1>
			<p><a href="` + metricsEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
	})
	logger.Infoln("Starting HTTP server on ", listenAddress)
	logger.Fatal(http.ListenAndServe(listenAddress, nil))
}
