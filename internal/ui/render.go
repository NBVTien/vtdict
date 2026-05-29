package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/NBVTien/vtdict/internal/ai"
	"github.com/NBVTien/vtdict/internal/dictionary"
	"github.com/NBVTien/vtdict/internal/storage"
	"github.com/charmbracelet/lipgloss"
)

var (
	wordStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	phoneStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	posStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	defStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	exStyle    = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("244"))
	transStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	aiStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Italic(true)

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(60)

	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	countStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
)

type RenderOpts struct {
	Phonetic bool
	Examples bool
	POS      bool
}

func RenderLookup(results []dictionary.Result, translation string, opts RenderOpts) {
	var b strings.Builder

	for i, r := range results {
		if i > 0 {
			b.WriteString("\n")
		}

		header := wordStyle.Render(r.Word)
		if opts.Phonetic {
			for _, p := range r.Phonetics {
				if p.Text != "" {
					header += "  " + phoneStyle.Render(p.Text)
					break
				}
			}
		}
		b.WriteString(header + "\n")

		for j, m := range r.Meanings {
			if j >= 3 {
				break
			}
			if opts.POS {
				b.WriteString("\n" + posStyle.Render(m.PartOfSpeech) + "\n")
			} else if j == 0 {
				b.WriteString("\n")
			}
			for k, d := range m.Definitions {
				if k >= 2 {
					break
				}
				b.WriteString(defStyle.Render("  • "+d.Definition) + "\n")
				if opts.Examples && d.Example != "" {
					b.WriteString(exStyle.Render("    \""+d.Example+"\"") + "\n")
				}
			}
		}

		if translation != "" {
			b.WriteString("\n" + labelStyle.Render("translation  ") + transStyle.Render(translation) + "\n")
		}
	}

	fmt.Println(boxStyle.Render(b.String()))
}

func RenderAILookup(def *ai.Definition, translation string, opts RenderOpts) {
	var b strings.Builder

	b.WriteString(wordStyle.Render(def.Word) + "\n")

	if opts.POS && def.PartOfSpeech != "" {
		b.WriteString("\n" + posStyle.Render(def.PartOfSpeech) + "\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(defStyle.Render("  • "+def.Definition) + "\n")

	if opts.Examples && def.Example != "" {
		b.WriteString(exStyle.Render("    \""+def.Example+"\"") + "\n")
	}

	if len(def.Synonyms) > 0 {
		b.WriteString("\n" + labelStyle.Render("synonyms  ") + dimStyle.Render(strings.Join(def.Synonyms, ", ")) + "\n")
	}

	if translation != "" {
		b.WriteString("\n" + labelStyle.Render("translation  ") + transStyle.Render(translation) + "\n")
	}

	b.WriteString("\n" + aiStyle.Render("⚡ via AI"))

	fmt.Println(boxStyle.Render(b.String()))
}

func RenderCached(definition, translation string, opts RenderOpts) {
	var b strings.Builder

	b.WriteString(defStyle.Render("  • "+definition) + "\n")

	if translation != "" {
		b.WriteString("\n" + labelStyle.Render("translation  ") + transStyle.Render(translation) + "\n")
	}

	fmt.Println(boxStyle.Render(b.String()))
}

func RenderHistory(entries []storage.Entry) {
	if len(entries) == 0 {
		fmt.Println(dimStyle.Render("No history yet. Look up a word first."))
		return
	}

	var rows []string
	rows = append(rows, lipgloss.NewStyle().Bold(true).Render(
		fmt.Sprintf("%-18s  %4s  %-12s  %s", "WORD", "CNT", "LAST SEEN", "DEFINITION"),
	))
	rows = append(rows, strings.Repeat("─", 68))

	for _, e := range entries {
		word := wordStyle.Render(fmt.Sprintf("%-18s", truncate(e.Word, 18)))
		count := countStyle.Render(fmt.Sprintf("%4d", e.LookupCount))
		ts := dimStyle.Render(fmt.Sprintf("%-12s", humanTime(e.LastLookedUp)))
		def := defStyle.Render(truncate(e.Definition, 28))
		rows = append(rows, fmt.Sprintf("%s  %s  %s  %s", word, count, ts, def))
	}

	fmt.Println(boxStyle.Width(72).Render(strings.Join(rows, "\n")))
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n-1]) + "…"
}

func humanTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return t.Format("Jan 02 2006")
	}
}
