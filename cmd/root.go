package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gorawler",
		Short: "My Tool!",
		Long:  "Gorawler is a web crawler written in Go.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Long)
			fmt.Println(cmd.UsageString())
		},
	}

	debug bool
	dbDsn string
)

func Execute() error {
	return rootCmd.Execute()
}
