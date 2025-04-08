package alias

import (
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/pkg/logger"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type confirm struct {
	question  string
	confirmed bool
	quitting  bool
}

type keyMap struct {
	yes key.Binding
	no  key.Binding
}

var keys = keyMap{
	yes: key.NewBinding(
		key.WithKeys("y", "Y"),
		key.WithHelp("y/Y", "yes"),
	),
	no: key.NewBinding(
		key.WithKeys("n", "N"),
		key.WithHelp("n/N", "no"),
	),
}

func (m confirm) Init() tea.Cmd {
	return nil
}

func (m confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.yes):
			m.confirmed = true
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.no):
			m.confirmed = false
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirm) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("%s (y/n) ", m.question)
}

var ResetCmd = &cobra.Command{
	Use:   "alias:reset",
	Short: "Reset all alias settings",
	// Short: "すべてのエイリアス設定をリセットします",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.New()

		m := confirm{question: "Are you sure you want to reset all alias settings?"}
		p := tea.NewProgram(m)
		model, err := p.Run()
		if err != nil {
			return err
		}

		final, ok := model.(confirm)
		if !ok {
			return fmt.Errorf("could not assert model type")
		}

		if !final.confirmed {
			log.Info("Reset canceled")
			return nil
		}

		if err := alias.Reset(); err != nil {
			return err
		}

		log.Info("All alias settings have been reset")
		return nil
	},
}
