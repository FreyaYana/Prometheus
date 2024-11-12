package metrics

// Латенси:
// histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path))
// Показывает 95-й процентиль времени выполнения запросов по маршрутам.

// Количество запросов по методам и статусам:
// sum by (method, status) (rate(http_requests_total[5m]))
// Отображает распределение количества запросов по методам и статусам.

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
)

func Init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestCount)
}

func PrometheusMiddleware(handler http.Handler) http.Handler {
	instrumentedHandler := promhttp.InstrumentHandlerDuration(
		requestDuration,
		promhttp.InstrumentHandlerCounter(
			requestCount,
			handler,
		),
	)
	return instrumentedHandler
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Prometheus!"))
}

// func main() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/hello", helloHandler)
// 	mux.Handle("/metrics", promhttp.Handler())

// 	http.ListenAndServe(":8080", wrappedMux)
// }

// func main() {

// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		message := r.URL.Query().Get("message")
// 		if message == "" {
// 			message = "Hello, World!"
// 		}

// 		// Симуляция ошибки
// 		if message == "error" {
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			return
// 		}

// 		fmt.Fprintf(w, "Echo: %s", message)
// 	})

// 	http.Handle("/echo", instrumentedHandler)
// 	http.Handle("/metrics", promhttp.Handler())

// 	http.ListenAndServe(":1234", nil)
// }

// func (handler http.Handler) UseMetrcis(w *http.Request) http.HandleFunc {

// 	// Оборачиваем обработчик в middleware для измерения количества запросов и латенции
// 	instrumentedHandler := promhttp.InstrumentHandlerDuration(
// 		requestDuration,
// 		promhttp.InstrumentHandlerCounter(
// 			requestCount,
// 			handler,
// 		))

// 	return instrumentedHandler
// }
