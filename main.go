package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pcfens/fireboard-exporter/pkg/fireboard"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		addr              = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
		authenticationKey = flag.String("key", "", "The authentication key used to connect to fireboard.io")
	)

	flag.Parse()
	prometheus.MustRegister(collectors.NewBuildInfoCollector())

	key := *authenticationKey
	fc := fireboard.New(key)

	prometheus.MustRegister(fc)

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
