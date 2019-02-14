package checkcmd

import "github.com/alecthomas/kingpin"

// CheckCmd control the check command
type Check struct {
	cmd *kingpin.CmdClause
}

const (
	CheckCmd     = "check"
)

// Action execute the `check` command
func (c *Check) Action([]string) {

}

// Init initialize the check cli commands
func (c *Check) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(CheckCmd, "Provide a return code on the release status")
}
