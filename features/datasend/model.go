package datasend

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	MenuMode Mode = iota
	CreateMode
)

type Field int

const (
	AppNameField Field = iota
	PasskeyField
	DescriptionField
)

type keyMap struct {
	Create key.Binding
	Submit key.Binding
	Back key.Binding
	Quit key.Binding
	Next key.Binding
	Prev key.Binding
}

var menuKeys = keyMap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create protocol"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "b"),
		key.WithHelp("esc/b", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("quit", "ctrl+c"),
	),
}

var createKeys = keyMap{
	Submit: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "submit"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab", "down"),
		key.WithHelp("tab/↓", "next field"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab/↑", "prev field"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	if k.Create.Enabled() {
		return []key.Binding{k.Create, k.Back, k.Quit}
	}
	return []key.Binding{k.Submit, k.Next, k.Back}
}

func (k keyMap) FullHelp() [][]key.Binding {
	if k.Create.Enabled() {
		return [][]key.Binding {
			{k.Back, k.Quit},
		}
	}
	return [][]key.Binding{
		{k.Submit, k.Next, k.Prev},
		{k.Back, k.Quit},
	}
}

type Protocol struct {
	AppName string
	Passkey string
	Description string
}

type Model struct {
	Mode Mode
	help help.Model
	inputs []textinput.Model
	focusIndex int
	protocols []Protocol
	statusMsg string
	onProtocolCreated func(Protocol)
	progress progress.Model
	keys keyMap
}

func NewModel() *Model {
	m := &Model{
		Mode: MenuMode,
		help: help.New(),
		keys: menuKeys,
		protocols: []Protocol{},
	}

	m.inputs = make([]textinput.Model, 3)

	m.inputs[AppNameField] = textinput.New()
	m.inputs[AppNameField].Placeholder = "my-app"
	m.inputs[AppNameField].Focus()
	m.inputs[AppNameField].CharLimit = 50
	m.inputs[AppNameField].Width = 40

	m.inputs[PasskeyField] = textinput.New()
	m.inputs[PasskeyField].Placeholder = "secret-key-123"
	m.inputs[PasskeyField].CharLimit = 100
	m.inputs[PasskeyField].Width = 40
	m.inputs[PasskeyField].EchoMode = textinput.EchoPassword
	m.inputs[PasskeyField].EchoCharacter = '.'

	m.inputs[DescriptionField] = textinput.New()
	m.inputs[DescriptionField].Placeholder = "API for my application"
	m.inputs[DescriptionField].CharLimit = 200
	m.inputs[DescriptionField].Width = 40

	return m
}

func (m *Model) SetProtocolCreatedCallback(fn func(Protocol)) {
	m.onProtocolCreated = fn
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Mode == MenuMode {
			return m.updateMenu(msg)
		} else {
			return m.updateCreate(msg)
		}
	}

	if m.Mode == CreateMode {
		cmd := m.updateInputs(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) updateMenu(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		m.Mode = CreateMode
		m.keys = createKeys
		m.focusIndex = 0
		m.inputs[0].Focus()
		m.statusMsg = ""
		return m, nil
	}
	return m, nil
}

func (m *Model) updateCreate(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.Mode = MenuMode
		m.keys = menuKeys
		m.resetInputs()
		return m, nil
	case "ctrl+s":
		if m.validateInputs() {
			protocol := Protocol{
				AppName:     m.inputs[AppNameField].Value(),
				Passkey:     m.inputs[PasskeyField].Value(),
				Description: m.inputs[DescriptionField].Value(),
			}
			m.protocols = append(m.protocols, protocol)

			if m.onProtocolCreated != nil {
				m.onProtocolCreated(protocol)
			}

			m.statusMsg = fmt.Sprintf("✓ Protocol '%s' created successfully!", protocol.AppName)
			m.Mode = MenuMode
			m.keys = menuKeys
			m.resetInputs()
		} else {
			m.statusMsg = "All fields are required!"
		}
		return m, nil
	case "tab", "down":
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		return m, m.updateFocus()
	case "shift+tab", "up":
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		return m, m.updateFocus()
	default:
		cmd := m.updateInputs(msg)
		return m, cmd
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *Model) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) validateInputs() bool {
	return m.inputs[AppNameField].Value() != "" &&
		m.inputs[PasskeyField].Value() != "" &&
		m.inputs[DescriptionField].Value() != ""
}

func (m *Model) resetInputs() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	m.focusIndex = 0
	m.statusMsg = ""
}

func (m Model) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

	if m.Mode == MenuMode {
		return m.viewMenu(titleStyle)
	}
	return m.viewCreate(titleStyle)
}

func (m Model) viewMenu(titleStyle lipgloss.Style) string {
	title := titleStyle.Render("Send Data")

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		Padding(1, 0)
	
		header := headerStyle.Render("Custom Data")

		info := lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Render("Create custom data transfer protocols for inter-app communication.\n")

		protocolsView := ""
		if len(m.protocols) > 0 {
			protocolsView = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Render(fmt.Sprintf("\nActive Protocols: %d\n", len(m.protocols)))
			
				for _, p := range m.protocols {
					protocolsView += lipgloss.NewStyle().
						Foreground(lipgloss.Color("green")).
						Render(fmt.Sprintf(". %s - %s\n", p.AppName, p.Description))
				}
		}

		status := ""
		if m.statusMsg != "" {
			status = "\n" + m.statusMsg + "\n"
		}

		helpView := m.help.View(m.keys)

		return lipgloss.NewStyle().
			Padding(1, 2).
			Render(title + "\n" + header + "\n" + info + protocolsView + status + "\n" + helpView)
}

func (m Model) viewCreate(titleStyle lipgloss.Style) string {
	title := titleStyle.Render("Create custom data protocol")

	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("\nDefine your custom API protocol:\n\n")

	fieldStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	form := ""
	labels := []string{"App Name:", "Passkey", "Description:"}

	for i, label := range labels {
		if i == m.focusIndex {
			form += focusedStyle.Render(label) + "\n"
		} else {
			form += fieldStyle.Render(label) + "\n"
		}
		form += m.inputs[i].View() + "\n\n"
	}

	noteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("yellow")).
		Italic(true)

	note := noteStyle.Render("Note: Your protocol will be available at:\nGET http://localhost:6767/{app_name}/init\nHeaders: X-App-Name, X-Passkey")

	status := ""
	if m.statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("red")).
			Bold(true)
		status = "\n" + statusStyle.Render(m.statusMsg) + "\n"
	}

	helpView := m.help.View(m.keys)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(title + "\n" + info + form + note + status + "\n" + helpView)
}