package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/mstcl/cider/handler"
)

const httpAddr = "0.0.0.0:8080"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler.Logger = logger

	http.HandleFunc("/", handler.Index)

	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		logger.Error("http", "error", err)
		os.Exit(1)
	}
}
