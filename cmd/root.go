package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var logLevelString string
var enableRichLogs bool
var disableRichLogs bool

var rootCmd = &cobra.Command{
	Use: "almanac",
	// TODO: Add description
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevelString, "log-level", "info", "Log level (debug, info, warn, error)")
	_ = rootCmd.RegisterFlagCompletionFunc(
		"log-level",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var matching []string

			for _, level := range []string{"debug", "info", "warn", "error"} {
				if strings.HasPrefix(level, toComplete) {
					matching = append(matching, level)
				}
			}

			return matching, cobra.ShellCompDirectiveNoFileComp
		},
	)

	rootCmd.PersistentFlags().BoolVar(
		&enableRichLogs,
		"enable-rich-logs",
		false,
		"Enable rich logs, even in non-terminal environments",
	)
	rootCmd.PersistentFlags().BoolVar(
		&disableRichLogs,
		"disable-rich-logs",
		false,
		"Disable rich logs, even in terminal environments",
	)
	rootCmd.PersistentFlags().String("content-dir", "content", "Directory containing content files")

	cobra.OnInitialize(setupLogging)
}

func setupLogging() {
	level, validLevel := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}[strings.ToLower(logLevelString)]

	if !validLevel {
		slog.Error(
			"Invalid log level",
			"log-level", logLevelString,
		)
		os.Exit(1)
	}

	if (term.IsTerminal(int(os.Stderr.Fd())) && !disableRichLogs) || enableRichLogs {
		slog.SetDefault(slog.New(
			tint.NewHandler(
				os.Stderr,
				&tint.Options{
					TimeFormat: time.Kitchen,
					Level:      level,
				},
			),
		))
	} else {
		slog.SetDefault(slog.New(
			slog.NewJSONHandler(
				os.Stderr,
				&slog.HandlerOptions{
					Level: level,
				},
			),
		))
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
