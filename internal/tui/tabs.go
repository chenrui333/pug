package tui

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/leg100/pug/internal/resource"
	"github.com/leg100/pug/internal/tui/keys"
)

const tabHeaderHeight = 2

// models implementing tabStatus can report a status that'll be rendered
// alongside the title in the tab header.
type tabStatus interface {
	TabStatus() string
}

// models implementing this can report info that'll be rendered on the opposite
// side away from the tab headers.
type tabSetInfo interface {
	TabSetInfo() string
}

// TabSet is a related set of tabs.
type TabSet struct {
	Tabs []Tab

	// Width and height of the content area
	width  int
	height int

	// The currently active tab
	active int

	info tabSetInfo
}

func NewTabSet(width, height int) TabSet {
	return TabSet{
		width:  width,
		height: height,
	}
}

func (m TabSet) WithTabSetInfo(i tabSetInfo) TabSet {
	m.info = i
	return m
}

// A tab is one of a set of tabs. A tab has a title, and an embedded model,
// which is responsible for the visible content under the tab.
type Tab struct {
	Model

	title string
}

// Init initializes the existing tabs in the collection.
func (m TabSet) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Tabs))
	for i, tab := range m.Tabs {
		cmds[i] = tab.Init()
	}
	return tea.Batch(cmds...)
}

var ErrDuplicateTab = errors.New("not allowed to create tabs with duplicate titles")

// AddTab adds a tab to the tab set, using the make and parent to construct the
// model associated with the tab. The title must be unique in the set. Upon
// success the associated model's Init() is returned for the caller to
// initialise the model.
func (m *TabSet) AddTab(maker Maker, parent resource.Resource, title string) (tea.Cmd, error) {
	for _, tab := range m.Tabs {
		if tab.title == title {
			return nil, ErrDuplicateTab
		}
	}

	model, err := maker.Make(parent, m.contentWidth(), m.contentHeight())
	if err != nil {
		return nil, err
	}
	m.Tabs = append(m.Tabs, Tab{Model: model, title: title})
	return model.Init(), nil
}

// SetActiveTab sets the currently active tab. If the tab index is outside the
// bounds of the current tab set then the active tab is automatically set to the
// last tab.
func (m *TabSet) SetActiveTab(tabIndex int) {
	if tabIndex < 0 || tabIndex > len(m.Tabs)-1 {
		m.active = len(m.Tabs) - 1
	} else {
		m.active = tabIndex
	}
}

// SetActiveTabWithTitle looks up a tab with a title and makes it the active
// tab. If no such tab exists no action is taken.
func (m *TabSet) SetActiveTabWithTitle(title string) {
	for i, tab := range m.Tabs {
		if tab.title == title {
			m.active = i
			return
		}
	}
}

func (m TabSet) Update(msg tea.Msg) (TabSet, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Navigation.TabNext):
			// Cycle tabs, going back to the first tab after the last tab
			if m.active == len(m.Tabs)-1 {
				m.active = 0
			} else {
				m.active = m.active + 1
			}
			return m, nil
		case key.Matches(msg, keys.Navigation.TabLast):
			m.active = max(m.active-1, 0)
			return m, nil
		}
		// Send other keys to active tab if there is one
		if len(m.Tabs) > 0 {
			cmd := m.updateTab(m.active, msg)
			return m, cmd
		}
	case BodyResizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Relay modified resize message onto each tab model
		m.updateTabs(BodyResizeMsg{
			Width:  m.contentWidth(),
			Height: m.contentHeight(),
		})
		return m, nil
	}

	// Updates each tab's respective model in-place.
	cmds = append(cmds, m.updateTabs(msg))

	return m, tea.Batch(cmds...)
}

func (m *TabSet) updateTabs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Tabs))
	for i := range m.Tabs {
		cmds[i] = m.updateTab(i, msg)
	}
	return tea.Batch(cmds...)
}

func (m *TabSet) updateTab(tabIndex int, msg tea.Msg) tea.Cmd {
	updated, cmd := m.Tabs[tabIndex].Update(msg)
	m.Tabs[tabIndex].Model = updated
	return cmd
}

var (
	activeTabStyle   = Bold.Copy().Foreground(lipgloss.Color("13"))
	inactiveTabStyle = Regular.Copy().Foreground(lipgloss.Color("250"))
)

func (m TabSet) View() string {
	var (
		tabHeaders       []string
		tabsHeadersWidth int
	)
	for i, t := range m.Tabs {
		var (
			headingStyle  lipgloss.Style
			underlineChar string
		)
		if i == m.active {
			headingStyle = activeTabStyle.Copy()
			underlineChar = "━"
		} else {
			headingStyle = inactiveTabStyle.Copy()
			underlineChar = "─"
		}
		heading := headingStyle.Copy().Padding(0, 1).Render(t.title)
		if status, ok := t.Model.(tabStatus); ok {
			heading += headingStyle.Copy().Bold(false).Padding(0, 1, 0, 0).Render(status.TabStatus())
		}
		underline := headingStyle.Render(strings.Repeat(underlineChar, Width(heading)))
		rendered := lipgloss.JoinVertical(lipgloss.Top, heading, underline)
		tabHeaders = append(tabHeaders, rendered)
		tabsHeadersWidth += Width(heading)
	}

	// Populate remaining space to the right of the tab headers with a faint
	// grey underline. If the tab set parent implements tabSetInfo then that'll
	// be called and its contents rendered above the underline.
	remainingWidth := max(0, m.width-tabsHeadersWidth)
	var rightSideInfo string
	if m.info != nil {
		rightSideInfo = Padded.Copy().
			Width(remainingWidth).
			Align(lipgloss.Right).
			Render(m.info.TabSetInfo())
	}
	tabHeadersFiller := lipgloss.JoinVertical(lipgloss.Top,
		rightSideInfo,
		inactiveTabStyle.Copy().Render(strings.Repeat("─", remainingWidth)),
	)
	tabHeaders = append(tabHeaders, tabHeadersFiller)

	// Join tab headers and filler together
	tabHeadersContainer := lipgloss.JoinHorizontal(lipgloss.Bottom, tabHeaders...)

	var tabContent string
	if len(m.Tabs) > 0 {
		tabContent = m.Tabs[m.active].View()
	}
	return lipgloss.JoinVertical(lipgloss.Top, tabHeadersContainer, tabContent)
}

// Width of the tab content area
func (m TabSet) contentWidth() int {
	return m.width
}

// Height of the tab content area
func (m TabSet) contentHeight() int {
	return m.height - tabHeaderHeight
}
