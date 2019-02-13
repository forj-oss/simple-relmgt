package tagcmd

import "github.com/alecthomas/kingpin"

// TagCmd control the tag-it command
type TagCmd struct {
}

// Action execute the `check` command
func (c *TagCmd) Action([]string) {

}

// Init initialize the check cli commands
func (c *TagCmd) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}

}
