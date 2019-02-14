package statecmd

import "github.com/alecthomas/kingpin"

// State control the status command
type State struct {
	cmd *kingpin.CmdClause
}

const (
	StateItCmd = "status"
)

// Action execute the `check` command
func (c *State) Action([]string) {

}

// Init initialize the check cli commands
func (c *State) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(StateItCmd, "Step to display release status")
}
