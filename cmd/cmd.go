package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
)

const (
	defaultWidth = 20
	listHeight   = 14
)

// nolint:mnd
var (
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
)

type item string
type itemDelegate struct{}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func Run(_ context.Context) error {
	items := []list.Item{
		item("artifactgpt"),
		item("buildgpt"),
		item("codegpt"),
		item("gitgpt"),
		item("lintgpt"),
		item("metalgpt"),
	}

	l := list.New(items, &itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What kind of gpt would you need to use?👇"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		return errors.Wrap(err, "failed to new program")
	}

	return nil
}

func (i item) FilterValue() string {
	return ""
}

func (d *itemDelegate) Height() int {
	return 1
}

func (d *itemDelegate) Spacing() int {
	return 0
}

func (d *itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// nolint:gocritic
func (d *itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s selected!🚀", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("See you next time!💖")
	}
	return "\n" + m.list.View()
}
