package tui

const (
	// MinHeight is the minimum height of the TUI.
	MinHeight = 24
	// Height of prompt including borders
	PromptHeight = 3
	// FooterHeight is the height of the footer at the bottom of the TUI.
	FooterHeight = 1
	// Height of help widget, including borders
	HelpWidgetHeight = 12
	// MinContentHeight is the minimum height of content above the footer.
	MinContentHeight = MinHeight - FooterHeight
	// MinContentWidth is the minimum width of the content.
	MinContentWidth = 80
	// minimum height of each pane
	minPaneHeight = 4
	// minimum width of each pane
	minPaneWidth = 20
	// defaultTopRightPaneHeight is the default height of the top right pane.
	defaultTopRightPaneHeight = 15
	// defaultLeftPaneWidth is the default width of the left pane.
	defaultLeftPaneWidth = 40
)

func init() {
	if (minPaneHeight*2)+PromptHeight+HelpWidgetHeight > MinContentHeight {
		panic("mininum heights of panes, prompt, footer, and help cannot exceed overall minimum height")
	}
	if minPaneWidth*2 > MinContentWidth {
		panic("minimum width of panes must be no more than half of the minimum content width")
	}
	if minPaneHeight > defaultTopRightPaneHeight {
		panic("default top right pane height must not be lower than the overall minimum height")
	}
	if minPaneWidth > defaultLeftPaneWidth {
		panic("default left pane width must not be lower than the overall minimum width")
	}
}
