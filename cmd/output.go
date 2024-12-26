package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.fogo.sh/almanac/pkg/content"
	"pkg.fogo.sh/almanac/pkg/content/extensions"
)

var outputCmd = &cobra.Command{
	Use:   "output",
	Args:  cobra.NoArgs,
	Short: "Output all pages to disk",
	Run: func(cmd *cobra.Command, args []string) {
		contentDir := must(cmd.Flags().GetString("content-dir"))

		resolver, err := extensions.NewDiscordUserResolver(
			extensions.DiscordUserResolverConfig{
				DiscordToken: viper.GetString("discord.token"),
				CachePath:    viper.GetString("discord.cache_path"),
			},
		)
		if err != nil {
			slog.Warn("Failed to create Discord user resolver, Discord user mentions will not be resolved", "error", err)
		}

		parser := content.Parser{DiscordUserResolver: resolver}

		pages, err := parser.DiscoverPages(contentDir)
		checkError(err, "failed to discover pages")

		outputDir := must(cmd.Flags().GetString("output-dir"))

		slog.Info(fmt.Sprintf("discovered %d pages, outputting to %s", len(pages), outputDir))

		err = parser.OutputAllPagesToDisk(pages, outputDir)
		checkError(err, "failed to output pages")

		slog.Info("done!")
	},
}

func init() {
	rootCmd.AddCommand(outputCmd)

	outputCmd.Flags().StringP("output-dir", "o", "./output/", "Directory to output pages to")
}
