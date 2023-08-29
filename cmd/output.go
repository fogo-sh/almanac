package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/fogo-sh/almanac/pkg/content"
)

var outputCmd = &cobra.Command{
	Use:   "output",
	Args:  cobra.NoArgs,
	Short: "Output all pages to disk",
	Run: func(cmd *cobra.Command, args []string) {
		contentDir := must(cmd.Flags().GetString("content-dir"))

		pages, err := content.DiscoverPages(contentDir)
		checkError(err, "failed to discover pages")

		outputDir := must(cmd.Flags().GetString("output-dir"))

		slog.Info(fmt.Sprintf("discovered %d pages, outputting to %s", len(pages), outputDir))

		err = content.OutputAllPagesToDisk(pages, outputDir)
		checkError(err, "failed to output pages")

		slog.Info("done!")
	},
}

func init() {
	rootCmd.AddCommand(outputCmd)

	outputCmd.Flags().StringP("output-dir", "o", "./output/", "Directory to output pages to")
}
