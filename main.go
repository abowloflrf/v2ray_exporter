package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	v2c             *Client
	v2rayAddr       string
	listenAddr      string
	metricsEndpoint string
	debugMode       bool
)

var cmd = &cobra.Command{
	Use:   "v2ray-exporter",
	Short: "v2ray-exporter is a exporter to collect traffic usage by each vmess user which can be collected by prometheus",
	PreRun: func(cmd *cobra.Command, args []string) {
		// init zap logger and exporter
		initLogger()
		prometheus.MustRegister(NewExporter())
	},
	Run: func(cmd *cobra.Command, args []string) {
		signals := make(chan os.Signal, 1)
		v2c = NewClient(v2rayAddr)
		go serveHTTP(listenAddr, metricsEndpoint)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		v2c.Close()
	},
}

func init() {
	cmd.PersistentFlags().StringVar(&v2rayAddr, "target", "127.0.0.1:10150", "v2ray grpc api endpoint")
	cmd.PersistentFlags().StringVar(&listenAddr, "listen", "127.0.0.1:9100", "address exporter to listen")
	cmd.PersistentFlags().StringVar(&metricsEndpoint, "endpoint", "/metrics", "enpoint for metrics")
	cmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "open debug mode")
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
	sugar.Info("Starting HTTP server on ", listenAddress)
	sugar.Fatal(http.ListenAndServe(listenAddress, nil))
}

func main() {
	if err := cmd.Execute(); err != nil {
		sugar.Error("Failed to start server", err)
		os.Exit(1)
	}
}
