package ui

import (
	"fmt"

	"github.com/NBVTien/vtdict/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	checkedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	hintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type pickerModel struct {
	entries  []storage.Entry
	cursor   int
	selected map[int]bool
	done     bool
	aborted  bool
}

func newPicker(entries []storage.Entry) pickerModel {
	return pickerModel{
		entries:  entries,
		selected: make(map[int]bool),
	}
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.aborted = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "a":
			if len(m.selected) == len(m.entries) {
				m.selected = make(map[int]bool)
			} else {
				for i := range m.entries {
					m.selected[i] = true
				}
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	if m.done || m.aborted {
		return ""
	}

	s := wordStyle.Render("Select words to remove") + "\n"
	s += hintStyle.Render("↑/↓ navigate  space select  a select all  enter confirm  q cancel") + "\n\n"

	for i, e := range m.entries {
		cursor := "  "
		if m.cursor == i {
			cursor = cursorStyle.Render("▶ ")
		}

		check := "[ ]"
		name := dimStyle.Render(e.Word)
		if m.selected[i] {
			check = checkedStyle.Render("[✓]")
			name = selectedStyle.Render(e.Word)
		}

		s += fmt.Sprintf("%s%s %s\n", cursor, check, name)
	}

	count := len(m.selected)
	if count > 0 {
		s += "\n" + countStyle.Render(fmt.Sprintf("%d selected", count))
	}

	return s
}

// PickWordsToDelete runs interactive picker, returns selected words. Empty = cancelled.
func PickWordsToDelete(entries []storage.Entry) ([]string, error) {
	if len(entries) == 0 {
		return nil, nil
	}

	m, err := tea.NewProgram(newPicker(entries)).Run()
	if err != nil {
		return nil, err
	}

	final := m.(pickerModel)
	if final.aborted || len(final.selected) == 0 {
		return nil, nil
	}

	var words []string
	for i := range final.selected {
		words = append(words, entries[i].Word)
	}
	return words, nil
}
