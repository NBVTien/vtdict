package ui

import (
	"fmt"
	"strings"

	"github.com/NBVTien/vtdict/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type field struct {
	label string
	key   string
	kind  string // "bool" | "string"
}

var fields = []field{
	{"Translation language", "lang", "string"},
	{"Auto-translate", "translate", "bool"},
	{"Show phonetic", "phonetic", "bool"},
	{"Show examples", "examples", "bool"},
	{"Show part of speech", "pos", "bool"},
}

type configModel struct {
	cfg      config.Config
	cursor   int
	editing  bool   // editing a string field
	input    string // current input buffer
	saved    bool
	aborted  bool
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	activeRow    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	inactiveRow  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	valueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	mutedValue   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Underline(true)
	savedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	configBox    = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			Width(48)
)

func newConfigEditor(cfg config.Config) configModel {
	return configModel{cfg: cfg}
}

func (m configModel) Init() tea.Cmd { return nil }

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				m.cfg.Lang = m.input
				m.editing = false
			case "esc":
				m.editing = false
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.input += msg.String()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.aborted = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(fields)-1 {
				m.cursor++
			}
		case " ", "enter":
			f := fields[m.cursor]
			if f.kind == "bool" {
				m.cfg = toggleBool(m.cfg, f.key)
			} else {
				m.editing = true
				m.input = fieldStringValue(m.cfg, f.key)
			}
		case "s", "ctrl+s":
			m.saved = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m configModel) View() string {
	if m.saved || m.aborted {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("vtdict config") + "\n\n")

	for i, f := range fields {
		cursor := "  "
		if m.cursor == i {
			cursor = cursorStyle.Render("▶ ")
		}

		label := inactiveRow.Render(fmt.Sprintf("%-22s", f.label))
		if m.cursor == i {
			label = activeRow.Render(fmt.Sprintf("%-22s", f.label))
		}

		var val string
		if f.kind == "bool" {
			v := fieldBoolValue(m.cfg, f.key)
			if v {
				val = valueStyle.Render("on")
			} else {
				val = mutedValue.Render("off")
			}
		} else {
			sv := fieldStringValue(m.cfg, f.key)
			if m.cursor == i && m.editing {
				val = inputStyle.Render(m.input + "▌")
			} else {
				val = valueStyle.Render(sv)
			}
		}

		b.WriteString(fmt.Sprintf("%s%s  %s\n", cursor, label, val))
	}

	b.WriteString("\n" + hintStyle.Render("↑/↓ navigate  space/enter toggle or edit  s save  q cancel"))

	return configBox.Render(b.String())
}

func toggleBool(cfg config.Config, key string) config.Config {
	switch key {
	case "translate":
		cfg.Translate = !cfg.Translate
	case "phonetic":
		cfg.Phonetic = !cfg.Phonetic
	case "examples":
		cfg.Examples = !cfg.Examples
	case "pos":
		cfg.POS = !cfg.POS
	}
	return cfg
}

func fieldBoolValue(cfg config.Config, key string) bool {
	switch key {
	case "translate":
		return cfg.Translate
	case "phonetic":
		return cfg.Phonetic
	case "examples":
		return cfg.Examples
	case "pos":
		return cfg.POS
	}
	return false
}

func fieldStringValue(cfg config.Config, key string) string {
	switch key {
	case "lang":
		return cfg.Lang
	}
	return ""
}

// RunConfigEditor opens interactive config TUI, returns updated config + whether saved
func RunConfigEditor(cfg config.Config) (config.Config, bool, error) {
	m, err := tea.NewProgram(newConfigEditor(cfg)).Run()
	if err != nil {
		return cfg, false, err
	}
	final := m.(configModel)
	return final.cfg, final.saved, nil
}
