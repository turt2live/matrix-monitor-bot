package webserver

import (
	"net/http"
	"io"
)

func InitServer(mux *http.ServeMux) {
	mux.Handle("/", &tempHandler{})
}

type tempHandler struct{}

func (tempHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}
