package cmd

import (
	"github.com/spf13/cobra"

	"github.com/fogo-sh/almanac/pkg/devserver"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Args:  cobra.NoArgs,
	Short: "Start the Almanac dev server",
	Run: func(cmd *cobra.Command, args []string) {
		server := devserver.NewServer(devserver.Config{
			Addr:             must(cmd.Flags().GetString("addr")),
			ContentDir:       must(cmd.Flags().GetString("content-dir")),
			UseBundledAssets: must(cmd.Flags().GetBool("use-bundled-assets")),
		})
		err := server.Start()
		checkError(err, "failed to start dev server")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("addr", "a", ":8080", "Address to listen on")
	serveCmd.Flags().BoolP("use-bundled-assets", "b", true, "Whether to use bundled assets embedded in the binary")
}
