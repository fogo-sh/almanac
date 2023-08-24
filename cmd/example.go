package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
)

var mdString = `# Almanac

This is some
- example
- markdown
- code
`

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Example markdown rendering command",
	RunE: func(cmd *cobra.Command, args []string) error {
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(mdString), &buf); err != nil {
			return fmt.Errorf("failed to convert markdown: %w", err)
		}

		fmt.Println(buf.String())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
