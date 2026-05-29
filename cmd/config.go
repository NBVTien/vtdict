package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NBVTien/vtdict/internal/config"
	"github.com/NBVTien/vtdict/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit configuration",
	Long:  "Run without subcommand to open interactive editor.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, saved, err := ui.RunConfigEditor(config.Get())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if saved {
			config.Set(cfg)
			if err := config.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Config saved.")
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Example: `  vtdict config set lang fr
  vtdict config set translate false
  vtdict config set phonetic true`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, val := strings.ToLower(args[0]), args[1]
		cfg := config.Get()
		switch key {
		case "lang":
			cfg.Lang = val
		case "translate":
			cfg.Translate = parseBool(val)
		case "phonetic":
			cfg.Phonetic = parseBool(val)
		case "examples":
			cfg.Examples = parseBool(val)
		case "pos":
			cfg.POS = parseBool(val)
		default:
			fmt.Fprintf(os.Stderr, "unknown key %q\nvalid keys: lang, translate, phonetic, examples, pos\n", key)
			os.Exit(1)
		}
		config.Set(cfg)
		if err := config.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Set %s = %s\n", key, val)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		key := strings.ToLower(args[0])
		switch key {
		case "lang":
			fmt.Println(cfg.Lang)
		case "translate":
			fmt.Println(cfg.Translate)
		case "phonetic":
			fmt.Println(cfg.Phonetic)
		case "examples":
			fmt.Println(cfg.Examples)
		case "pos":
			fmt.Println(cfg.POS)
		default:
			fmt.Fprintf(os.Stderr, "unknown key %q\n", key)
			os.Exit(1)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all config values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Width(20)
		valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
		pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

		rows := [][]string{
			{"lang", cfg.Lang},
			{"translate", fmt.Sprintf("%v", cfg.Translate)},
			{"phonetic", fmt.Sprintf("%v", cfg.Phonetic)},
			{"examples", fmt.Sprintf("%v", cfg.Examples)},
			{"pos", fmt.Sprintf("%v", cfg.POS)},
		}
		for _, r := range rows {
			fmt.Printf("%s%s\n", keyStyle.Render(r[0]), valStyle.Render(r[1]))
		}
		fmt.Printf("\n%s\n", pathStyle.Render("config: "+config.Path()))
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all config to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Reset(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Config reset to defaults.")
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd, configGetCmd, configListCmd, configResetCmd)
	rootCmd.AddCommand(configCmd)
}

func parseBool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return strings.ToLower(s) == "on" || s == "1" || strings.ToLower(s) == "yes"
	}
	return v
}
