package utils

import (
	"io"
	"log/slog"
)

func DeferredClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		slog.Error("Failed to close file", "error", err)
	}
}
