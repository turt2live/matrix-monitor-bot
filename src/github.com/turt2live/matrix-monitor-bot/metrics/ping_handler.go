package metrics

import (
	"net/http"
	"io"
)

type PingHandler struct{}

func (PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "react-sdk-config-server is alive")
	return
}
