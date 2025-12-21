package ui

import (
	"freeport/api"
	"freeport/config"
	"freeport/features/dataview"
	"freeport/features/datasend"
	"freeport/features/settings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View int

const (
	MenuView View = iota
	DataViewView
	DataSendView
	SettingsView
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
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
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Back, k.Quit},
	}
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	list     list.Model
	help     help.Model
	keys     keyMap
	view     View
	config   *config.Config
	width    int
	height   int

	dataViewModel *dataview.Model
	dataSendModel *datasend.Model
	settingsModel *settings.Model
}

func NewModel() Model {
	cfg := config.Load()

	items := []list.Item{
		item{title: "View Data", desc: "View system data and API information"},
		item{title: "Send Data", desc: "Send data through the API bus"},
		item{title: "Settings", desc: "Configure application settings"},
		item{title: "Exit", desc: "Exit the application"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("170")).
		BorderForeground(lipgloss.Color("170"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("243")).
		BorderForeground(lipgloss.Color("170"))

	l := list.New(items, delegate, 0, 0)
	l.Title = cfg.WelcomeMessage
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true).
		Padding(0, 0, 1, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	h := help.New()

	dataSendModel := datasend.NewModel()
	dataSendModel.SetProtocolCreatedCallback(func(p datasend.Protocol) {
		api.RegisterProtocol(p.AppName, p.Passkey, p.Description)
	})

	return Model{
		list:          l,
		help:          h,
		keys:          keys,
		view:          MenuView,
		config:        cfg,
		dataViewModel: dataview.NewModel(),
		dataSendModel: dataSendModel,
		settingsModel: settings.NewModel(cfg),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.help.Width = msg.Width
	}

	switch m.view {
	case DataViewView:
		return m.updateDataView(msg)
	case DataSendView:
		return m.updateDataSend(msg)
	case SettingsView:
		return m.updateSettings(msg)
	default:
		return m.updateMenu(msg)
	}
}

func (m Model) View() string {
	switch m.view {
	case DataViewView:
		return m.dataViewModel.View(m.width, m.height)
	case DataSendView:
		return m.dataSendModel.View(m.width, m.height)
	case SettingsView:
		return m.settingsModel.View(m.width, m.height)
	default:
		return m.viewMenu()
	}
}