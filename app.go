package main

import (
	"fmt"
	"net/http"
	"strconv"

	"log"
	"time"

	monitoring "github.com/FreyaYana/Prometheus/monitoring"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PORT is the TCP port number the server will listen to
var PORT = ":1234"

// Depending on what kind of information you want to collect and expose,
// you will have to use a different metric type.
var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "custom",
			Name:      "my_counter",
			Help:      "This is my counter",
		})
)

func main() {

	monitoring.Init()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		delayStr := r.URL.Query().Get("delay")
		if delayStr != "" {
			if delay, err := strconv.Atoi(delayStr); err == nil {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			} else {
				fmt.Fprintf(w, "Invalid delay value. Using no delay.\n")
			}
		}

		fmt.Fprintf(w, "Echo: %s", delayStr)
	})
	mux.HandleFunc("/echo2", func(w http.ResponseWriter, r *http.Request) {
		delayStr := r.URL.Query().Get("delay")
		if delayStr != "" {
			if delay, err := strconv.Atoi(delayStr); err == nil {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			} else {
				fmt.Fprintf(w, "Invalid delay value. Using no delay.\n")
			}
		}

		counter.Inc()
		fmt.Fprintf(w, "Echo: %s", delayStr)
	})

	wrappedMux := monitoring.PrometheusMiddleware(mux)

	//prometheus.MustRegister(counter)

	// Запуск HTTP сервера
	log.Println("Starting server on", PORT)
	log.Fatal(http.ListenAndServe(PORT, wrappedMux))
}
