package draftcmd

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

// Draft control the draft-it command
type Draft struct {
	cmd *kingpin.CmdClause
}

const (
	DraftItCmd = "draft-it"
)

// Action execute the `check` command
func (c *Draft) Action([]string) {
	fmt.Printf("%s not currently defined\n", DraftItCmd)
	os.Exit(5) // Function not defined
}

// Init initialize the check cli commands
func (c *Draft) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(DraftItCmd, "Step to create a draft release")
}
