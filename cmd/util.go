package cmd

import (
	"log/slog"
	"os"
)

func checkError(err error, message string) {
	if err != nil {
		slog.Error(message, "error", err)
		os.Exit(1)
	}
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
