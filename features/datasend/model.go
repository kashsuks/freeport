package datasend

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Back key.Binding
	Quit key.Binding
}

var keys = keyMap{
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
	return [][]key.Binding {
		{k.Back, k.Quit},
	}
}

type Model struct {
	progress progress.Model
	help help.Model
	keys keyMap
}

func NewModel() *Model {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	h := help.New()

	return &Model{
		progress: prog,
		help: h,
		keys: keys,
	}
}

func (m Model) View(width, height int) string {
	titledStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("178")).
		Padding(1, 0)

	title := titledStyle.Render("Send Data")
	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("\nThis will send data through the API bus\n")

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1, 0)
	
	status := statusStyle.Render("Status: Ready\nProgress: " + m.progress.ViewAs(0.0))

	helpView := m.help.View(m.keys)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(title + info + "\n" + status + "\n\n[Coming Soon]\n\n" + helpView)
}