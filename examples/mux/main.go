package main

import (
	"fmt"
	"log"
	"net/http"

	mobiledetect "github.com/housemecn/go-mobile-detect"
)

// Handler .
type Handler struct{}

// ServeHTTP .
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%#v", mobiledetect.Device(r))
}

func main() {
	log.Println("Starting local server http://localhost:10001/check (cmd+click to open from terminal)")
	mux := http.NewServeMux()
	h := &Handler{}
	mux.Handle("/check", h)
	http.ListenAndServe(":10001", mobiledetect.HandlerMux(mux, nil))
}
