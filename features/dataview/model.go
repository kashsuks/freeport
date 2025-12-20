package dataview

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type keyMap struct {
	Back key.Binding
	Quit key.Binding
}

var keys = keyMap {
	Back: key.NewBinding(
		key.WithKeys("esc", "b"),
		key.WithHelp("esc/b", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

type Model struct {
	Table table.Model
	Help help.Model
	Keys keyMap
}

func NewModel() *Model {
	columns := []table.Column{
		{Title: "Key", Width: 20},
		{Title: "Value", Width: 40},
		{Title: "Type", Width: 15},
	}

	rows := []table.Row{
		{"status", "Coming Soon", "string"},
		{"data", "No data available", "string"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	h := help.New()

	return &Model {
		Table: t,
		Help: h,
		Keys: keys,
	}
}

func (m Model) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

	title := titleStyle.Render("View Data")
	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("\nDisplays system data from the API")

	helpView := m.Help.View(m.Keys)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(title + info + "\n" + baseStyle.Render(m.Table.View()) + "\n\n" + helpView)
}