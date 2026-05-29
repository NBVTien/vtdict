package cmd

import (
	"fmt"
	"os"

	"github.com/NBVTien/vtdict/internal/config"
	"github.com/NBVTien/vtdict/internal/storage"
	"github.com/spf13/cobra"
)

var searchWord string

var rootCmd = &cobra.Command{
	Use:   "vtdict [word]",
	Short: "Terminal dictionary with Vietnamese translation",
	Args:  cobra.ArbitraryArgs,
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if searchWord != "" {
			return lookupWord(searchWord)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if searchWord != "" {
			return nil // already handled in PersistentPreRunE
		}
		if len(args) == 0 {
			return cmd.Help()
		}
		return lookupWord(args[0])
	},
}

func Execute() {
	if err := config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}
	if err := storage.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "storage init failed: %v\n", err)
		os.Exit(1)
	}
	defer storage.Close()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
