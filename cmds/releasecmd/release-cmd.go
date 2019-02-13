package releasecmd

import "github.com/alecthomas/kingpin"

// ReleaseCmd control the release-it command
type ReleaseCmd struct {
}

// Action execute the `check` command
func (c *ReleaseCmd) Action([]string) {

}

// Init initialize the check cli commands
func (c *ReleaseCmd) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}

}
