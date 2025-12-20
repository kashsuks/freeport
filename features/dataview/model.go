package dataview

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type keyMap struct {
	Query key.Binding
	Back  key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Query: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "query battery"),
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
	return []key.Binding{k.Query, k.Back, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Query, k.Back, k.Quit},
	}
}

type batteryDataMsg struct {
	time    string
	battery string
	appName string
	err     error
}

type Model struct {
	Table       table.Model
	Help        help.Model
	Keys        keyMap
	Input       textinput.Model
	loading     bool
	lastQueried string
	errorMsg    string
}

func NewModel() *Model {
	columns := []table.Column{
		{Title: "Field", Width: 20},
		{Title: "Value", Width: 40},
	}

	rows := []table.Row{
		{"Status", "Ready - Press 'q' to query battery"},
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

	ti := textinput.New()
	ti.Placeholder = "Press 'q' to query battery data"

	return &Model{
		Table:   t,
		Help:    h,
		Keys:    keys,
		Input:   ti,
		loading: false,
	}
}

func queryBatteryData() tea.Msg {
	resp, err := http.Get("http://localhost:6767/system/battery")
	if err != nil {
		return batteryDataMsg{err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return batteryDataMsg{err: err}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return batteryDataMsg{err: err}
	}

	return batteryDataMsg{
		time:    data["time"].(string),
		battery: fmt.Sprintf("%.0f%%", data["battery"].(float64)),
		appName: data["app_name"].(string),
		err:     nil,
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" && !m.loading {
			m.loading = true
			m.errorMsg = ""
			return m, queryBatteryData
		}
	case batteryDataMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMsg = fmt.Sprintf("Error: %v", msg.err)
			m.Table.SetRows([]table.Row{
				{"Status", "Error querying battery"},
				{"Error", msg.err.Error()},
			})
		} else {
			m.lastQueried = time.Now().Format("15:04:05")
			m.errorMsg = ""
			m.Table.SetRows([]table.Row{
				{"Time", msg.time},
				{"Battery", msg.battery},
				{"App Name", msg.appName},
				{"Last Queried", m.lastQueried},
			})
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View(width, height int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1, 0)

	title := titleStyle.Render(" View Data - Battery Info")
	
	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Render("\nQuery battery information from the API bus.\n")

	status := ""
	if m.loading {
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("yellow")).
			Render("\n Loading battery data...\n")
	} else if m.errorMsg != "" {
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("red")).
			Render("\n"+ m.errorMsg + "\n")
	} else if m.lastQueried != "" {
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("green")).
			Render(fmt.Sprintf("\n Last updated: %s\n", m.lastQueried))
	}

	helpView := m.Help.View(m.Keys)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(title + info + status + "\n" + baseStyle.Render(m.Table.View()) + "\n\n" + helpView)
}