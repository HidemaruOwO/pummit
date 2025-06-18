package prompt

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// エディター選択用のアイテム構造体
type EditorItem struct {
	name        string
	command     string
	description string
	available   bool
}

func (i EditorItem) FilterValue() string { return i.name }
func (i EditorItem) Title() string       { return i.name }
func (i EditorItem) Description() string { return i.description }

// エディター選択モデル
type EditorSelectorModel struct {
	list        list.Model
	choice      string
	customInput string
	inputMode   bool
	quitting    bool
	keys        editorKeyMap
	configPath  string
}

type editorKeyMap struct {
	selectItem key.Binding
	custom     key.Binding
	quit       key.Binding
}

var editorKeys = editorKeyMap{
	selectItem: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	custom: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "custom command"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// 利用可能なエディターのリストを生成
func getAvailableEditors() []list.Item {
	editors := []EditorItem{
		{"vim", "vim", "Vi/Vim editor", false},
		{"nvim", "nvim", "Neovim", false},
		{"nano", "nano", "Nano editor", false},
		{"emacs", "emacs", "Emacs", false},
		{"code", "code", "Visual Studio Code", false},
		{"gedit", "gedit", "GEdit", false},
		{"kate", "kate", "Kate", false},
		{"vi", "vi", "Vi editor", false},
		{"cat", "cat", "Display file content", false},
		{"bat", "bat", "Display file with syntax highlighting", false},
		{"gat", "gat", "Gat editor", false},
	}

	items := make([]list.Item, 0)

	// $EDITOR環境変数をチェックして最初に追加
	if editorEnv := os.Getenv("EDITOR"); editorEnv != "" {
		item := EditorItem{
			name:        fmt.Sprintf("$EDITOR (%s)", editorEnv),
			command:     editorEnv,
			description: "Editor from environment variable (recommended)",
			available:   true,
		}
		items = append(items, item)
	}

	// 各エディターの利用可能性をチェックして利用可能なもののみ追加
	for _, editor := range editors {
		_, err := exec.LookPath(editor.command)
		if err == nil {
			// 利用可能なエディターのみ追加
			editor.available = true
			items = append(items, editor)
		}
		// 利用不可のエディターは表示しない
	}

	return items
}

// エディター選択インターフェースの作成
func NewEditorSelector(configPath string) EditorSelectorModel {
	items := getAvailableEditors()

	l := list.New(items, list.NewDefaultDelegate(), 80, 14)
	l.Title = "Select an editor to open the TOML configuration file"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return EditorSelectorModel{
		list:       l,
		keys:       editorKeys,
		configPath: configPath,
	}
}

// スタイル定義
var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("86"))

	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4)

	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1)

	customInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Render
)

func (m EditorSelectorModel) Init() tea.Cmd {
	return nil
}

func (m EditorSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.inputMode {
			switch msg.String() {
			case "enter":
				// カスタムコマンドを実行
				if strings.TrimSpace(m.customInput) != "" {
					m.choice = m.customInput
					m.quitting = true
					return m, tea.Quit
				}
				m.inputMode = false
				m.customInput = ""
			case "ctrl+c", "esc":
				m.inputMode = false
				m.customInput = ""
			default:
				if msg.String() == "backspace" {
					if len(m.customInput) > 0 {
						m.customInput = m.customInput[:len(m.customInput)-1]
					}
				} else if len(msg.Runes) > 0 {
					m.customInput += string(msg.Runes)
				}
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.custom):
			m.inputMode = true
			m.customInput = ""
			return m, nil

		case key.Matches(msg, m.keys.selectItem):
			if item, ok := m.list.SelectedItem().(EditorItem); ok {
				if item.available {
					m.choice = item.command
					m.quitting = true
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m EditorSelectorModel) View() string {
	if m.quitting {
		return ""
	}

	view := m.list.View()

	if m.inputMode {
		view += "\n\n" + customInputStyle("Enter custom command:")
		view += "\n> " + m.customInput + "█"
		view += "\n\n" + customInputStyle("Enter: execute, Esc: cancel")
	} else {
		view += "\n\n" + customInputStyle("c: custom command, q: quit")
	}

	return view
}

// エディター選択の実行結果
type EditorResult struct {
	Command  string
	Canceled bool
}

// エディター選択インターフェースを実行
func RunEditorSelector(configPath string) (EditorResult, error) {
	m := NewEditorSelector(configPath)
	p := tea.NewProgram(m)
	model, err := p.Run()
	if err != nil {
		return EditorResult{}, fmt.Errorf("failed to run editor selector: %w", err)
	}

	finalModel, ok := model.(EditorSelectorModel)
	if !ok {
		return EditorResult{}, fmt.Errorf("internal error: could not assert editor selector model type")
	}

	if finalModel.choice == "" {
		return EditorResult{Canceled: true}, nil
	}

	return EditorResult{Command: finalModel.choice, Canceled: false}, nil
}

// エディターでファイルを開く
func OpenWithEditor(command, filePath string) error {
	// コマンドを空白で分割してargvを作成
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("invalid command: %s", command)
	}

	cmdName := parts[0]
	args := append(parts[1:], filePath)

	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
