package metrics

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitServer(mux *http.ServeMux) {
	initMetrics()
	mux.Handle("/_monitorbot/metrics", promhttp.Handler())
	mux.Handle("/_monitorbot/ping", PingHandler{})
}
