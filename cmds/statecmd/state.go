package statecmd

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

// State control the status command
type State struct {
	cmd *kingpin.CmdClause
}

const (
	StateItCmd = "status"
)

// Action execute the `check` command
func (c *State) Action([]string) {
	fmt.Printf("%s not currently defined\n", StateItCmd)
	os.Exit(5) // Function not defined

}

// Init initialize the check cli commands
func (c *State) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(StateItCmd, "Step to display release status")
}
