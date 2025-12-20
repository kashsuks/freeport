package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c"{
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				selected := selectedItem.(item)
				switch selected.title {
				case "View Data":
					m.view = DataViewView
					return m, nil
				case "Send Data":
					m.view = DataSendView
					return m, nil
				case "Settings":
					m.view = SettingsView
					return m, nil
				case "Exit":
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) viewMenu() string {
	return docStyle.Render(m.list.View() + "\n" + m.help.View(m.keys))
}