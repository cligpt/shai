package metal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mgutz/ansi"
)

const (
	footerScrolled             = 100
	useHighPerformanceRenderer = false
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type MetalGPTModel struct {
	Content  string
	ready    bool
	viewport viewport.Model
}

func (m *MetalGPTModel) Init() tea.Cmd {
	return nil
}

func (m *MetalGPTModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.Content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MetalGPTModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m *MetalGPTModel) headerView() string {
	title := titleStyle.Render("Node Status")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *MetalGPTModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*footerScrolled))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func AddColor(s string) string {
	var output string

	jsonLines := strings.Split(s, "\n")

	for i := 0; i < len(jsonLines); i++ {
		if strings.Contains(jsonLines[i], "\":") {
			splitIndex := strings.Index(jsonLines[i], ":")
			key := jsonLines[i][:splitIndex+2]
			value := jsonLines[i][splitIndex+2:]
			key = ansi.Color(key, "cyan")
			if strings.Contains(value, `"`) {
				value = ansi.Color(value, "green")
			} else if strings.Contains(value, "[") {
				value = ansi.Color(value, "yellow")
			} else if strings.Contains(value, `true`) || strings.Contains(value, `false`) {
				value = ansi.Color(value, "blue")
			} else if strings.Contains(value, ".") {
				value = ansi.Color(value, "white")
			} else if strings.Contains(value, "{") {
				value = ansi.Color(value, "yellow")
			} else {
				value = ansi.Color(value, "green")
			}
			output = output + key + value + "\n"
		} else {
			if strings.Contains(jsonLines[i], "{") || strings.Contains(jsonLines[i], "}") ||
				strings.Contains(jsonLines[i], "[") || strings.Contains(jsonLines[i], "]") {
				output = output + ansi.Color(jsonLines[i], "yellow") + "\n"
			} else if strings.Contains(jsonLines[i], "\"") {
				output = output + ansi.Color(jsonLines[i], "green") + "\n"
			} else {
				output = output + ansi.Color(jsonLines[i], "white") + "\n"
			}
		}
	}

	return output
}
