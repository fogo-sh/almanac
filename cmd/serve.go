package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"pkg.fogo.sh/almanac/pkg/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Args:  cobra.NoArgs,
	Short: "Start the Almanac server",
	Run: func(cmd *cobra.Command, args []string) {
		serverInstance := server.NewServer(server.Config{
			Addr:             must(cmd.Flags().GetString("addr")),
			ContentDir:       must(cmd.Flags().GetString("content-dir")),
			UseBundledAssets: must(cmd.Flags().GetBool("use-bundled-assets")),

			UseDiscordOAuth:     must(cmd.Flags().GetBool("use-discord-oauth")),
			DiscordClientId:     viper.GetString("discord.client_id"),
			DiscordClientSecret: viper.GetString("discord.client_secret"),
			DiscordCallbackUrl:  viper.GetString("discord.callback_url"),
			DiscordGuildId:      viper.GetString("discord.guild_id"),
			SessionSecret:       viper.GetString("discord.session_secret"),
		})
		err := serverInstance.Start()
		checkError(err, "failed to start server")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("addr", "a", ":8080", "Address to listen on")
	serveCmd.Flags().BoolP("use-bundled-assets", "b", true, "Whether to use bundled assets embedded in the binary")
	serveCmd.Flags().Bool("use-discord-oauth", false, "Whether to use Discord OAuth for authentication")
}
