package tagcmd

import "github.com/alecthomas/kingpin"

// Tag control the tag-it command
type Tag struct {
	cmd *kingpin.CmdClause
}

const (
	TagItCmd = "tag-it"
)

// Action execute the `check` command
func (c *Tag) Action([]string) {

}

// Init initialize the check cli commands
func (c *Tag) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(TagItCmd, "Step to tag the code, before build")
}
