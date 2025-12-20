package settings

import (
	"freeport/config"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	NormalMode Mode = iota
	EditMode
)

type keyMap struct {
	Edit key.Binding
	Back key.Binding
	Quit key.Binding
	Save key.Binding
	Cancel key.Binding
}

var NormalKeys = keyMap{
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit welcome message"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "b"),
		key.WithHelp("esc/b", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var EditKeys = keyMap{
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	if k.Edit.Enabled() {
		return []key.Binding{k.Edit, k.Back, k.Quit}
	}
	return []key.Binding{k.Save, k.Cancel}
}

func (k keyMap) FullHelp() [][]key.Binding {
	if k.Edit.Enabled() {
		return [][]key.Binding{
			{k.Edit, k.Back},
			{k.Quit},
		}
	}
	return [][]key.Binding{
		{k.Save, k.Cancel, k.Quit},
	}
}

type Model struct {
	Config *config.Config
	Mode Mode
	Input textinput.Model
	Help help.Model
	Keys keyMap
	StatusMsg string
}

func NewModel(cfg *config.Config) *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter Welcome Message..."
	ti.CharLimit = 100
	ti.Width = 50

	h := help.New()

	return &Model{
		Config: cfg,
		Mode: NormalMode,
		Input: ti,
		Help: h,
		Keys: NormalKeys,
		StatusMsg: "",
	}
}

func (m Model) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

	title := titleStyle.Render("Settings")

	currentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1, 0)

	current := currentStyle.Render("Current Welcome Message:\n" +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Render("\"" + m.Config.WelcomeMessage + "\""))

	var content string
	if m.Mode == EditMode {
		editStyle := lipgloss.NewStyle().
			Padding(1, 0)
		content = editStyle.Render("Edit Welcome Message:\n" + m.Input.View())
	} else {
		content = ""
	}

	status := ""
	if m.StatusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)
		status = "\n" + statusStyle.Render(m.StatusMsg)
	}

	helpView := m.Help.View(m.Keys)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(title + "\n" + current + content + status + "\n\n" + helpView)
}