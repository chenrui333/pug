package table

import (
	"github.com/leg100/pug/internal/resource"
	"github.com/leg100/pug/internal/run"
)

var (
	IDColumn = Column{
		Key:   "id",
		Title: "ID", Width: resource.IDEncodedMaxLen,
	}

	ModuleColumn = Column{
		Key:            "module",
		Title:          "MODULE",
		TruncationFunc: TruncateLeft,
		FlexFactor:     3,
	}
	WorkspaceColumn = Column{
		Key:        "workspace",
		Title:      "WORKSPACE",
		FlexFactor: 2,
	}
	RunColumn = Column{
		Key:        "run",
		Title:      "RUN",
		Width:      resource.IDEncodedMaxLen,
		FlexFactor: 1,
	}
	TaskColumn = Column{
		Key:        "task",
		Title:      "TASK",
		Width:      resource.IDEncodedMaxLen,
		FlexFactor: 1,
	}
	RunStatusColumn = Column{
		Key:   "run_status",
		Title: "STATUS",
		Width: run.MaxStatusLen,
	}
	RunChangesColumn = Column{
		Key:   "run_changes",
		Title: "CHANGES",
		Width: 10,
	}
)
