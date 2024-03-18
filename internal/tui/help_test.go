package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/stretchr/testify/assert"
)

func Test_render(t *testing.T) {
	tests := []struct {
		name     string
		bindings []key.Binding
		want     string
	}{
		{
			"single column",
			[]key.Binding{
				key.NewBinding(key.WithHelp("a", "aaa")),
				key.NewBinding(key.WithHelp("b", "bbb")),
				key.NewBinding(key.WithHelp("c", "ccc")),
			},
			"a aaa   \nb bbb   \nc ccc   ",
		},
		{
			"two columns",
			[]key.Binding{
				key.NewBinding(key.WithHelp("a", "aaa")),
				key.NewBinding(key.WithHelp("b", "bbb")),
				key.NewBinding(key.WithHelp("c", "ccc")),
				key.NewBinding(key.WithHelp("d", "ddd")),
			},
			"a aaa   d ddd   \nb bbb           \nc ccc           ",
		},
		{
			"three columns",
			[]key.Binding{
				key.NewBinding(key.WithHelp("a", "aaa")),
				key.NewBinding(key.WithHelp("b", "bbb")),
				key.NewBinding(key.WithHelp("c", "ccc")),
				key.NewBinding(key.WithHelp("d", "ddd")),
				key.NewBinding(key.WithHelp("e", "eee")),
				key.NewBinding(key.WithHelp("f", "fff")),
				key.NewBinding(key.WithHelp("g", "ggg")),
			},
			"a aaa   d ddd   g ggg   \nb bbb   e eee           \nc ccc   f fff           ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortHelpView(tt.bindings, 30)
			assert.Equal(t, tt.want, got)
		})
	}
}