package draftcmd

import "github.com/alecthomas/kingpin"

// DraftCmd control the draft-it command
type DraftCmd struct {
}

// Action execute the `check` command
func (c *DraftCmd) Action([]string) {

}

// Init initialize the check cli commands
func (c *DraftCmd) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}

}
