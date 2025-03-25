package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const defaultPort = 7040

func main() {
	for _, err := range startApplication() {
		slog.With("error", err).Error("execution failed")
	}
}

func startApplication() (errors []error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	slog.Info("Starting")

	port := defaultPort
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      http.HandlerFunc(handler),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
	}

	go func() {
		slog.With("port", port).Info("will try to start http server")
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("closed http server")
			} else {
				slog.With("error", err).Error("error during srv.ListenAndServe")
				os.Exit(1)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	wait := time.Second * 15
	shutdownCtx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(shutdownCtx)

	slog.Info("shutting down")

	return
}

func handler(responseWriter http.ResponseWriter, request *http.Request) {
	var numbers SlowList
	if err := json.NewDecoder(request.Body).Decode(&numbers); err != nil {
		slog.Debug("failed to read request body: " + err.Error())
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	sortMethod := numbers.QuickSort
	switch os.Getenv("SORT_TYPE") {
	case "random":
		sortMethod = numbers.RandomSort
	case "bubble":
		sortMethod = numbers.BubbleSort
	}

	if err := sortMethod(request.Context()); err != nil {
		responseWriter.Write([]byte(("failed to sort numbers: " + err.Error())))
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(responseWriter).Encode(numbers)
}
