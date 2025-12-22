package datasend

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	MenuMode Mode = iota
	CreateMode
	SuccessMode
	ManageMode
	CreateMethodMode
)

type Field int

const (
	AppNameField Field = iota
	PasskeyField
	DescriptionField
)

type MethodField int

const (
	MethodNameField MethodField = iota
	MethodDescField
)

type FocusButton int

const (
	OkButton FocusButton = iota
	BackButton
)

type keyMap struct {
	Create key.Binding
	Submit key.Binding
	Back   key.Binding
	Quit   key.Binding
	Next   key.Binding
	Prev   key.Binding
	Select key.Binding
	Left   key.Binding
	Right  key.Binding
}

var menuKeys = keyMap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create protocol"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
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

var successKeys = keyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

var manageKeys = keyMap{
	Create: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new method"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "b"),
		key.WithHelp("esc/b", "back"),
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
	if k.Left.Enabled() {
		return []key.Binding{k.Left, k.Right, k.Select}
	}
	if k.Submit.Enabled() {
		return []key.Binding{k.Submit, k.Next, k.Back}
	}
	return []key.Binding{k.Select, k.Back, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	if k.Create.Enabled() {
		return [][]key.Binding{
			{k.Create, k.Back, k.Quit},
		}
	}
	if k.Left.Enabled() {
		return [][]key.Binding{
			{k.Left, k.Right, k.Select},
		}
	}
	if k.Submit.Enabled() {
		return [][]key.Binding{
			{k.Submit, k.Next, k.Prev},
			{k.Back, k.Quit},
		}
	}
	return [][]key.Binding{
		{k.Select, k.Back, k.Quit},
	}
}

type CustomMethod struct {
	Name        string
	Description string
}

type Protocol struct {
	AppName     string
	Passkey     string
	Description string
	Methods     []CustomMethod
}

type Model struct {
	Mode                  Mode
	help                  help.Model
	inputs                []textinput.Model
	methodInputs          []textinput.Model
	focusIndex            int
	protocols             []Protocol
	currentProtocol       *Protocol
	statusMsg             string
	onProtocolCreated     func(Protocol)
	onMethodCreated       func(string, string, string)
	keys                  keyMap
	focusedButton         FocusButton
	selectedProtocolIndex int
}

func NewModel() *Model {
	m := &Model{
		Mode:      MenuMode,
		help:      help.New(),
		keys:      menuKeys,
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
	m.inputs[PasskeyField].EchoCharacter = '•'

	m.inputs[DescriptionField] = textinput.New()
	m.inputs[DescriptionField].Placeholder = "API for my application"
	m.inputs[DescriptionField].CharLimit = 200
	m.inputs[DescriptionField].Width = 40

	m.methodInputs = make([]textinput.Model, 2)

	m.methodInputs[MethodNameField] = textinput.New()
	m.methodInputs[MethodNameField].Placeholder = "get-data"
	m.methodInputs[MethodNameField].Focus()
	m.methodInputs[MethodNameField].CharLimit = 50
	m.methodInputs[MethodNameField].Width = 40

	m.methodInputs[MethodDescField] = textinput.New()
	m.methodInputs[MethodDescField].Placeholder = "Retrieves data from the API"
	m.methodInputs[MethodDescField].CharLimit = 200
	m.methodInputs[MethodDescField].Width = 40

	return m
}

func (m *Model) SetProtocolCreatedCallback(fn func(Protocol)) {
	m.onProtocolCreated = fn
}

func (m *Model) SetMethodCreatedCallback(fn func(string, string, string)) {
	m.onMethodCreated = fn
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch m.Mode {
	case MenuMode:
		return m.updateMenu(msg)
	case CreateMode:
		return m.updateCreate(msg)
	case SuccessMode:
		return m.updateSuccess(msg)
	case ManageMode:
		return m.updateManage(msg)
	case CreateMethodMode:
		return m.updateCreateMethod(msg)
	}

	return m, nil
}

func (m *Model) updateMenu(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			m.Mode = CreateMode
			m.keys = createKeys
			m.focusIndex = 0
			m.inputs[0].Focus()
			m.statusMsg = ""
			return m, nil
		case "enter":
			if len(m.protocols) > 0 && m.selectedProtocolIndex < len(m.protocols) {
				m.currentProtocol = &m.protocols[m.selectedProtocolIndex]
				m.Mode = ManageMode
				m.keys = manageKeys
				return m, nil
			}
		case "down", "j":
			if m.selectedProtocolIndex < len(m.protocols)-1 {
				m.selectedProtocolIndex++
			}
		case "up", "k":
			if m.selectedProtocolIndex > 0 {
				m.selectedProtocolIndex--
			}
		}
	}
	return m, nil
}

func (m *Model) updateCreate(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
					Methods: []CustomMethod{
						{Name: "init", Description: "Initialize connection"},
					},
				}
				m.protocols = append(m.protocols, protocol)
				m.currentProtocol = &m.protocols[len(m.protocols)-1]

				if m.onProtocolCreated != nil {
					m.onProtocolCreated(protocol)
				}

				m.Mode = SuccessMode
				m.keys = successKeys
				m.focusedButton = OkButton
				m.resetInputs()
				return m, nil
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

	return m, nil
}

func (m *Model) updateSuccess(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left", "h":
			m.focusedButton = OkButton
		case "right", "l":
			m.focusedButton = BackButton
		case "enter":
			if m.focusedButton == OkButton {
				m.Mode = ManageMode
				m.keys = manageKeys
			} else {
				m.Mode = MenuMode
				m.keys = menuKeys
				m.currentProtocol = nil
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) updateManage(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "b":
			m.Mode = MenuMode
			m.keys = menuKeys
			m.currentProtocol = nil
			return m, nil
		case "n":
			m.Mode = CreateMethodMode
			m.keys = createKeys
			m.focusIndex = 0
			m.methodInputs[0].Focus()
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) updateCreateMethod(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.Mode = ManageMode
			m.keys = manageKeys
			m.resetMethodInputs()
			return m, nil
		case "ctrl+s":
			if m.validateMethodInputs() {
				method := CustomMethod{
					Name:        m.methodInputs[MethodNameField].Value(),
					Description: m.methodInputs[MethodDescField].Value(),
				}

				if m.currentProtocol != nil {
					m.currentProtocol.Methods = append(m.currentProtocol.Methods, method)

					if m.onMethodCreated != nil {
						m.onMethodCreated(m.currentProtocol.AppName, method.Name, method.Description)
					}

					m.statusMsg = fmt.Sprintf("✓ Method '%s' created!", method.Name)
				}

				m.Mode = ManageMode
				m.keys = manageKeys
				m.resetMethodInputs()
			} else {
				m.statusMsg = "All fields are required!"
			}
			return m, nil
		case "tab", "down":
			m.focusIndex = (m.focusIndex + 1) % len(m.methodInputs)
			return m, m.updateMethodFocus()
		case "shift+tab", "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = len(m.methodInputs) - 1
			}
			return m, m.updateMethodFocus()
		default:
			cmd := m.updateMethodInputs(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *Model) updateMethodInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.methodInputs))
	for i := range m.methodInputs {
		m.methodInputs[i], cmds[i] = m.methodInputs[i].Update(msg)
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

func (m *Model) updateMethodFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.methodInputs))
	for i := 0; i < len(m.methodInputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.methodInputs[i].Focus()
		} else {
			m.methodInputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *Model) validateInputs() bool {
	return m.inputs[AppNameField].Value() != "" &&
		m.inputs[PasskeyField].Value() != "" &&
		m.inputs[DescriptionField].Value() != ""
}

func (m *Model) validateMethodInputs() bool {
	return m.methodInputs[MethodNameField].Value() != "" &&
		m.methodInputs[MethodDescField].Value() != ""
}

func (m *Model) resetInputs() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	m.focusIndex = 0
	m.statusMsg = ""
}

func (m *Model) resetMethodInputs() {
	for i := range m.methodInputs {
		m.methodInputs[i].SetValue("")
	}
	m.focusIndex = 0
	m.statusMsg = ""
}

func (m Model) View(width, height int) string {
	switch m.Mode {
	case MenuMode:
		return m.viewMenu()
	case CreateMode:
		return m.viewCreate()
	case SuccessMode:
		return m.viewSuccess()
	case ManageMode:
		return m.viewManage()
	case CreateMethodMode:
		return m.viewCreateMethod()
	}
	return ""
}

func (m Model) viewMenu() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

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

		for i, p := range m.protocols {
			prefix := "  "
			if i == m.selectedProtocolIndex {
				prefix = "> "
			}
			protocolsView += lipgloss.NewStyle().
				Foreground(lipgloss.Color("green")).
				Render(fmt.Sprintf("%s%s - %s (%d methods)\n", prefix, p.AppName, p.Description, len(p.Methods)))
		}
		protocolsView += "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true).
			Render("Press Enter to manage selected protocol\n")
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

func (m Model) viewCreate() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

	title := titleStyle.Render("Create Custom Data Protocol")

	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("\nDefine your custom API protocol:\n\n")

	fieldStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	form := ""
	labels := []string{"App Name:", "Passkey:", "Description:"}

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

func (m Model) viewSuccess() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("green")).
		Padding(1, 0)

	title := titleStyle.Render(fmt.Sprintf("✓ Protocol '%s' Created Successfully!", m.currentProtocol.AppName))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1, 0)

	info := infoStyle.Render("\nAvailable Methods:\n")

	methodsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(60)

	methodsList := ""
	for _, method := range m.currentProtocol.Methods {
		methodsList += lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			Render(fmt.Sprintf("• %s\n", method.Name))
		methodsList += lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Render(fmt.Sprintf("  %s\n\n", method.Description))
	}

	methods := methodsStyle.Render(methodsList)

	usageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Padding(1, 0)

	usage := usageStyle.Render("Usage:\n") +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("yellow")).
			Render(fmt.Sprintf("curl -H \"X-App-Name: %s\" -H \"X-Passkey: [your-passkey]\" \\\n  http://localhost:6767/%s/init\n\n",
				m.currentProtocol.AppName, m.currentProtocol.AppName))

	okStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("170")).
		Padding(0, 3).
		Bold(true)

	okDimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 3)

	backStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 3)

	backFocusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("243")).
		Padding(0, 3).
		Bold(true)

	var okButton, backButtonView string
	if m.focusedButton == OkButton {
		okButton = okStyle.Render("  OK  ")
		backButtonView = backStyle.Render(" Back ")
	} else {
		okButton = okDimStyle.Render("  OK  ")
		backButtonView = backFocusedStyle.Render(" Back ")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, okButton, "  ", backButtonView)

	helpView := m.help.View(m.keys)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(title + "\n" + info + methods + "\n" + usage + buttons + "\n\n" + helpView)
}

func (m Model) viewManage() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

	if m.currentProtocol == nil {
		return "No protocol selected"
	}

	title := titleStyle.Render(fmt.Sprintf("%s", m.currentProtocol.AppName))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true)

	desc := descStyle.Render(m.currentProtocol.Description + "\n")

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		Padding(1, 0)

	header := headerStyle.Render("API Methods")

	methodsView := ""
	for _, method := range m.currentProtocol.Methods {
		methodsView += lipgloss.NewStyle().
			Foreground(lipgloss.Color("green")).
			Bold(true).
			Render(fmt.Sprintf("\n• %s\n", method.Name))
		methodsView += lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Render(fmt.Sprintf("  %s\n", method.Description))
		methodsView += lipgloss.NewStyle().
			Foreground(lipgloss.Color("yellow")).
			Render(fmt.Sprintf("  GET http://localhost:6767/%s/%s\n", m.currentProtocol.AppName, method.Name))
	}

	status := ""
	if m.statusMsg != "" {
		status = "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("green")).
			Render(m.statusMsg) + "\n"
	}

	helpView := m.help.View(m.keys)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(title + "\n" + desc + "\n" + header + methodsView + status + "\n\n" + helpView)
}

func (m Model) viewCreateMethod() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

	title := titleStyle.Render("Create New API Method")

	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render(fmt.Sprintf("\nCreating method for: %s\n\n", m.currentProtocol.AppName))

	fieldStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	form := ""
	labels := []string{"Method Name:", "Description:"}

	for i, label := range labels {
		if i == m.focusIndex {
			form += focusedStyle.Render(label) + "\n"
		} else {
			form += fieldStyle.Render(label) + "\n"
		}
		form += m.methodInputs[i].View() + "\n\n"
	}

	noteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("yellow")).
		Italic(true)

	note := noteStyle.Render(fmt.Sprintf("Your method will be available at:\nGET http://localhost:6767/%s/{method_name}", m.currentProtocol.AppName))

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