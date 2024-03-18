package tui

import (
	"github.com/leg100/pug/internal/resource"
	"github.com/leg100/pug/internal/tui/table"
	"github.com/mattn/go-runewidth"
)

var (
	moduleColumn = table.Column{
		Title:          "MODULE",
		TruncationFunc: runewidth.TruncateLeft,
		FlexFactor:     3,
	}
	workspaceColumn = table.Column{
		Title:      "WORKSPACE",
		FlexFactor: 2,
	}
	runColumn = table.Column{
		Title:      "RUN",
		Width:      resource.IDEncodedMaxLen,
		FlexFactor: 1,
	}
	taskColumn = table.Column{
		Title:      "TASK",
		Width:      resource.IDEncodedMaxLen,
		FlexFactor: 1,
	}
)

// parentColumns returns appropriate columns for a table depending upon the parent kind.
func parentColumns(table modelKind, parentKind resource.Kind) (columns []table.Column) {
	switch table {
	case ModuleListKind:
		// Always render module path on modules table
		columns = append(columns, moduleColumn)
	case WorkspaceListKind:
		switch parentKind {
		case resource.Global:
			// Only show module path if global workspaces table.
			columns = append(columns, moduleColumn)
		}
		// Always render workspace name on workspaces table
		columns = append(columns, workspaceColumn)
	case RunListKind:
		switch parentKind {
		case resource.Global:
			// Show all parent columns on global runs table
			columns = append(columns, moduleColumn)
			fallthrough
		case resource.Module:
			// Show workspace and run columns on module runs table.
			columns = append(columns, workspaceColumn)
		}
		// Always render run ID on runs table
		columns = append(columns, runColumn)
	case TaskListKind:
		switch parentKind {
		case resource.Global:
			// Show all parent columns on global tasks table
			columns = append(columns, moduleColumn)
			fallthrough
		case resource.Module:
			// Show workspace, run, and task columns on module tasks table.
			columns = append(columns, workspaceColumn)
			fallthrough
		case resource.Workspace:
			// Show run and task columns on workspace tasks table.
			columns = append(columns, runColumn)
		}
		// Always render task ID on tasks table
		columns = append(columns, taskColumn)
	}
	return
}

// parentCells returns appropriate cells depending upon the table and the parent
// resource.
func parentCells(tbl modelKind, parentKind resource.Kind, res resource.Resource) (cells []table.Cell) {
	switch tbl {
	case ModuleListKind:
		// Always render module path on modules table
		cells = append(cells, table.Cell{Str: res.Module().String()})
	case WorkspaceListKind:
		switch parentKind {
		case resource.Global:
			// Only show module path if global workspaces table.
			cells = append(cells, table.Cell{Str: res.Module().String()})
		}
		// Always render workspace name on workspaces table
		cells = append(cells, table.Cell{Str: res.Workspace().String()})
	case RunListKind:
		switch parentKind {
		case resource.Global:
			// Show all parent cells on global runs table
			cells = append(cells, table.Cell{Str: res.Module().String()})
			fallthrough
		case resource.Module:
			// Show workspace and run cells on module runs table.
			cells = append(cells, table.Cell{Str: res.Workspace().String()})
		}
		// Always render run ID on runs table
		cells = append(cells, table.Cell{Str: res.Run().ID().String()})
	case TaskListKind:
		switch parentKind {
		case resource.Global:
			// Show all parent cells on global tasks table
			cells = append(cells, table.Cell{Str: res.Module().String()})
			fallthrough
		case resource.Module:
			// Show workspace and run cells on module tasks table. A task doesn't
			// always belong to a workspace however, so render a blank string
			// for workspace if it doesn't.
			if res.Workspace() != nil {
				cells = append(cells, table.Cell{Str: res.Workspace().String()})
			} else {
				cells = append(cells, table.Cell{})
			}
			fallthrough
		case resource.Workspace:
			// Show run and task cells on workspace tasks table. A task doesn't
			// always belong to a run however, so render a blank string for run
			// if it doesn't.
			if res.Run() != nil {
				cells = append(cells, table.Cell{Str: res.Run().ID().String()})
			} else {
				cells = append(cells, table.Cell{})
			}
		}
		// Always render task ID on tasks table
		cells = append(cells, table.Cell{Str: res.ID().String()})
	}
	return
}