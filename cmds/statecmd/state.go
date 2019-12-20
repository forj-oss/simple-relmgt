package statecmd

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"simple-relmgt/core"

	"github.com/alecthomas/kingpin"
	version "github.com/hashicorp/go-version"
)

// State control the status command
type State struct {
	cmd *kingpin.CmdClause

	config *core.Config
	github *core.Github
	git    *core.Git

	versionFile      string
	extractVersionRE *regexp.Regexp

	releaseVersion string

	release *core.Release
}

const (
	StateItCmd = "status"
)

// Action execute the `check` command
func (c *State) Action([]string) (code int) {
	c.config = core.NewConfig("release-mgt.yaml")
	
	_, err := c.config.Load()
	kingpin.FatalIfError(err, "Unable to load %s properly.", c.config.Filename())

	c.github = core.NewGithub()

	err = c.github.CheckGithub()
	kingpin.FatalIfError(err, "Unable to get github-release")

	c.git = core.NewGit()

	err = c.git.OpenRepo()
	kingpin.FatalIfError(err, "Unable to open the local repository.")

	var data []byte
	data, err = ioutil.ReadFile(c.versionFile)
	if err != nil {
		fmt.Printf("Unable to read release version file %s. %s", c.versionFile, err)
		return
	}

	result := c.extractVersionRE.FindStringSubmatch(string(data))
	if result == nil {
		fmt.Printf("Release version file (%s) found, but version string has not been detected from '%s'.", c.versionFile, core.DefaultExtractVersion)
		return
	}
	c.releaseVersion = result[1]
	fmt.Printf("Release version detected: %s (in %s)\n", c.releaseVersion, c.versionFile)

	c.release = core.NewRelease()
	_, err = c.release.CheckVersion(c.releaseVersion)
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("Release %s ready to be published.\n", c.releaseVersion)
	}
	return
}

// Init initialize the check cli commands
func (c *State) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(StateItCmd, "Step to display release status")

	c.versionFile = core.DefaultVersionFile

	var err error
	c.extractVersionRE, err = regexp.Compile(fmt.Sprintf(core.DefaultExtractVersion, version.SemverRegexpRaw))
	kingpin.FatalIfError(err, "Unable to initialize check command")

}
