package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"freeport/features/settings"
)

func (m Model) updateDataView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "b":
			m.view = MenuView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.dataViewModel.Table, cmd = m.dataViewModel.Table.Update(msg)
	return m, cmd
}

func (m Model) updateDataSend(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "b":
			m.view = MenuView
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.settingsModel.Mode == settings.EditMode {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				if m.settingsModel.Input.Value() != "" {
					m.config.WelcomeMessage = m.settingsModel.Input.Value()
					m.config.Save()
					m.list.Title = m.config.WelcomeMessage
					m.settingsModel.StatusMsg = "✓ Welcome message saved!"
					m.settingsModel.Input.SetValue("")
					m.settingsModel.Input.Blur()
					m.settingsModel.Mode = settings.NormalMode
					m.settingsModel.Keys = settings.NormalKeys
				}
				return m, nil
			case "esc":
				m.settingsModel.Input.SetValue("")
				m.settingsModel.Input.Blur()
				m.settingsModel.Mode = settings.NormalMode
				m.settingsModel.Keys = settings.NormalKeys
				m.settingsModel.StatusMsg = ""
				return m, nil
			}
			m.settingsModel.Input, cmd = m.settingsModel.Input.Update(msg)
			return m, cmd
		} else {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc", "b":
				m.view = MenuView
				m.settingsModel.StatusMsg = ""
				return m, nil
			case "e":
				m.settingsModel.Mode = settings.EditMode
				m.settingsModel.Keys = settings.EditKeys
				m.settingsModel.Input.SetValue(m.config.WelcomeMessage)
				m.settingsModel.Input.Focus()
				m.settingsModel.StatusMsg = ""
				return m, nil
			}
		}
	}
	return m, nil
}