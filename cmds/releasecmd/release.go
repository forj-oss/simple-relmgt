package releasecmd

import (
	"fmt"

	"github.com/alecthomas/kingpin"
)

// Release control the release-it command
type Release struct {
	cmd *kingpin.CmdClause
}

const (
	ReleaseItCmd = "release-it"
)

// Action execute the `check` command
func (c *Release) Action([]string) (code int) {
	fmt.Printf("%s not currently defined\n", ReleaseItCmd)
	return 5 // Function not defined

}

// Init initialize the check cli commands
func (c *Release) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(ReleaseItCmd, "Step to release the code after build success.")
	return
}
