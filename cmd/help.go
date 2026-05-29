package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help",
	RunE: func(cmd *cobra.Command, args []string) error {
		if searchWord != "" {
			return lookupWord(searchWord)
		}
		printHelp()
		return nil
	},
}

func init() {
	rootCmd.SetHelpCommand(helpCmd)
	configCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printConfigHelp()
	})
}

func printHelp() {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	section := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	cmd := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	flag := lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(60)

	row := func(left, right string) string {
		return fmt.Sprintf("  %-28s %s", left, right)
	}

	var b strings.Builder
	b.WriteString(title.Render("vtdict") + muted.Render(" — terminal dictionary") + "\n\n")

	b.WriteString(section.Render("USAGE") + "\n")
	b.WriteString(row(cmd.Render("vtdict <word>"), desc.Render("look up a word")) + "\n")
	b.WriteString(row(cmd.Render("vtdict -s <word>"), desc.Render("force-lookup (e.g. reserved names)")) + "\n")
	b.WriteString("\n")

	b.WriteString(section.Render("COMMANDS") + "\n")
	b.WriteString(row(cmd.Render("vtdict log"), desc.Render("show lookup history")) + "\n")
	b.WriteString(row(cmd.Render("vtdict config"), desc.Render("interactive config editor")) + "\n")
	b.WriteString(row(cmd.Render("vtdict version"), desc.Render("show current version")) + "\n")
	b.WriteString(row(cmd.Render("vtdict update"), desc.Render("update to latest release")) + "\n")
	b.WriteString(row(cmd.Render("vtdict help"), desc.Render("show this help")) + "\n")
	b.WriteString("\n")

	b.WriteString(section.Render("FLAGS (any command)") + "\n")
	b.WriteString(row(flag.Render("-s, --search <word>"), desc.Render("force word lookup")) + "\n")
	b.WriteString(row(flag.Render("-l, --lang <code>"), desc.Render("translate to language (e.g. fr, ja)")) + "\n")
	b.WriteString(row(flag.Render("--no-translate"), desc.Render("skip translation")) + "\n")
	b.WriteString("\n")

	b.WriteString(section.Render("LOG FLAGS") + "\n")
	b.WriteString(row(flag.Render("--clear"), desc.Render("interactive picker to remove words")) + "\n")
	b.WriteString(row(flag.Render("--clear --all"), desc.Render("wipe all history")) + "\n")
	b.WriteString(row(flag.Render("--clear --before <dur>"), desc.Render("remove older than 7d / 2w / 1m")) + "\n")
	b.WriteString(row(flag.Render("--clear --word <word>"), desc.Render("remove one word")) + "\n")
	b.WriteString(row(flag.Render("-n <num>"), desc.Render("limit entries shown")) + "\n")
	b.WriteString("\n")

	b.WriteString(muted.Render("config: ~/.config/vtdict/config.toml"))

	fmt.Println(box.Render(b.String()))
}

func printConfigHelp() {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	section := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	cmd := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	flag := lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(60)

	row := func(left, right string) string {
		return fmt.Sprintf("  %-28s %s", left, right)
	}

	var b strings.Builder
	b.WriteString(title.Render("vtdict config") + muted.Render(" — view or edit configuration") + "\n\n")

	b.WriteString(section.Render("SUBCOMMANDS") + "\n")
	b.WriteString(row(cmd.Render("vtdict config"), desc.Render("interactive TUI editor")) + "\n")
	b.WriteString(row(cmd.Render("vtdict config set <key> <val>"), desc.Render("set a value")) + "\n")
	b.WriteString(row(cmd.Render("vtdict config get <key>"), desc.Render("read a value")) + "\n")
	b.WriteString(row(cmd.Render("vtdict config list"), desc.Render("show all + config file path")) + "\n")
	b.WriteString(row(cmd.Render("vtdict config reset"), desc.Render("restore defaults")) + "\n")
	b.WriteString("\n")

	b.WriteString(section.Render("KEYS") + "\n")
	b.WriteString(row(flag.Render("lang"), desc.Render("translation language (default: vi)")) + "\n")
	b.WriteString(row(flag.Render("translate"), desc.Render("auto-translate on lookup (default: true)")) + "\n")
	b.WriteString(row(flag.Render("phonetic"), desc.Render("show phonetic transcription (default: true)")) + "\n")
	b.WriteString(row(flag.Render("examples"), desc.Render("show usage examples (default: true)")) + "\n")
	b.WriteString(row(flag.Render("pos"), desc.Render("show part of speech (default: true)")) + "\n")
	b.WriteString("\n")

	b.WriteString(section.Render("AI FALLBACK KEYS") + "\n")
	b.WriteString(row(flag.Render("ai-fallback"), desc.Render("use AI when word not found (default: false)")) + "\n")
	b.WriteString(row(flag.Render("ai-provider"), desc.Render("openai | claude | groq | ollama")) + "\n")
	b.WriteString(row(flag.Render("ai-key"), desc.Render("API key (or set via env var)")) + "\n")
	b.WriteString(row(flag.Render("ai-model"), desc.Render("override provider default model")) + "\n")
	b.WriteString(row(flag.Render("ai-base-url"), desc.Render("custom OpenAI-compatible endpoint")) + "\n")
	b.WriteString("\n")

	b.WriteString(muted.Render("config: ~/.config/vtdict/config.toml"))

	fmt.Println(box.Render(b.String()))
}
