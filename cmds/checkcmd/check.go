package checkcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"simple-relmgt/core"

	"github.com/alecthomas/kingpin"
	version "github.com/hashicorp/go-version"
)

// Check control the check command
type Check struct {
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
	CheckCmd              = "check"
	defaultVersionFile    = "version.go"
	defaultExtractVersion = ` *VERSION *= *[\"'](%s)["']`
)

// Action execute the `check` command
func (c *Check) Action([]string) {
	c.config = core.NewConfig("release-mgt.yaml")

	c.github = core.NewGithub()

	err := c.github.CheckGithub()
	kingpin.FatalIfError(err, "Unable to get github-release")

	c.git = core.NewGit()

	err = c.git.OpenRepo()
	kingpin.FatalIfError(err, "Unable to open the local repository.")

	var data []byte
	data, err = ioutil.ReadFile(c.versionFile)
	if err != nil {
		fmt.Printf("Unable to read release version file %s. %s", c.versionFile, err)
		os.Exit(3)
	}

	result := c.extractVersionRE.FindStringSubmatch(string(data))
	if result == nil {
		fmt.Printf("Release version file (%s) found, but version string has not been detected from '%s'.", c.versionFile, defaultExtractVersion)
		os.Exit(2)
	}
	c.releaseVersion = result[1]
	fmt.Printf("Release version detected: %s (in %s)\n", c.releaseVersion, c.versionFile)

	c.release = core.NewRelease()
	var code int
	code, err = c.release.CheckVersion(c.releaseVersion)
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("Release %s ready to be published.\n", c.releaseVersion)
	}
	os.Exit(code)
}

// Init initialize the check cli commands
func (c *Check) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(CheckCmd, "Provide a return code on the release status")

	c.versionFile = defaultVersionFile

	var err error
	c.extractVersionRE, err = regexp.Compile(fmt.Sprintf(defaultExtractVersion, version.SemverRegexpRaw))
	kingpin.FatalIfError(err, "Unable to initialize check command")

}
