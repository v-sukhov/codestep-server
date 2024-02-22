package main

import (
	"net/http"
)

type IndexHtmlHandler struct{}

func (h *IndexHtmlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}
