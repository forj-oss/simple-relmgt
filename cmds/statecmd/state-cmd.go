package statecmd

import "github.com/alecthomas/kingpin"

// StateCmd control the status command
type StateCmd struct {
}

// Action execute the `check` command
func (c *StateCmd) Action([]string) {

}

// Init initialize the check cli commands
func (c *StateCmd) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}

}
