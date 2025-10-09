package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("component", "api")
	slog.SetDefault(logger)

	slog.With("hello", "world")

	slog.Error("thisis")
}
