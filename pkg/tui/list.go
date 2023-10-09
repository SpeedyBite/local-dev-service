package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Option interface {
	Title() string
	Description() string
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return fmt.Sprintf("%s %s", i.title, i.desc) }

type model struct {
	list           list.Model
	selectedChanel chan<- string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) processTeaMessage(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		close(m.selectedChanel)
		return tea.Quit
	case "enter":
		if m.list.SelectedItem() != nil {
			title := m.list.SelectedItem().(item).Title()
			go func() {
				defer close(m.selectedChanel)
				m.selectedChanel <- title
			}()
			return tea.Quit
		}
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd := m.processTeaMessage(msg)
		if cmd != nil {
			return m, cmd
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func SelectList(
	options []Option,
	title string) (int, error) {
	items := make([]list.Item, len(options))

	for i, option := range options {
		items[i] = item{
			title: option.Title(),
			desc:  option.Description(),
		}
	}
	ch := make(chan string)

	selectedTitle := ""
	m := model{
		list:           list.New(items, list.NewDefaultDelegate(), 0, 0),
		selectedChanel: ch,
	}
	m.list.Title = title

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	selectedTitle, ok := <-ch
	if !ok {
		return -1, errors.Errorf("no option selected")
	}
	selectedIndex := 0
	for i, option := range options {
		if option.Title() == selectedTitle {
			selectedIndex = i
			break
		}
	}
	return selectedIndex, nil
}
