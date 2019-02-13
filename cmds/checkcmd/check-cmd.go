package checkcmd

import "github.com/alecthomas/kingpin"

// CheckCmd control the check command
type CheckCmd struct {
}

// Action execute the `check` command
func (c *CheckCmd) Action([]string) {

}

// Init initialize the check cli commands
func (c *CheckCmd) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}

}
