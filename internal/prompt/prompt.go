package prompt

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Question  string
	Confirmed bool
	Quitting  bool
	keys      keyMap
}

type Result struct {
	Confirmed bool
}

type keyMap struct {
	yes key.Binding
	no  key.Binding
}

var defaultKeys = keyMap{
	yes: key.NewBinding(
		key.WithKeys("y", "Y"),
		key.WithHelp("Y", "yes"),
	),
	no: key.NewBinding(
		key.WithKeys("n", "N"),
		key.WithHelp("n", "no"),
	),
}

func New(question string) Model {
	return Model{
		Question: question,
		keys:     defaultKeys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.yes):
			m.Confirmed = true
			m.Quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.no):
			m.Confirmed = false
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	return fmt.Sprintf("%s (%s/%s) ", m.Question, m.keys.yes.Help().Key, m.keys.no.Help().Key)
}

func Run(question string) (Result, error) {
	m := New(question)
	p := tea.NewProgram(m)
	model, err := p.Run()
	if err != nil {
		return Result{}, fmt.Errorf("failed to run prompt: %w", err)
	}

	finalModel, ok := model.(Model)
	if !ok {
		return Result{}, fmt.Errorf("internal error: could not assert prompt model type")
	}

	return Result{Confirmed: finalModel.Confirmed}, nil
} 
