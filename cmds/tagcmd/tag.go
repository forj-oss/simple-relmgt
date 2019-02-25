package tagcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"simple-relmgt/core"

	"github.com/alecthomas/kingpin"
	version "github.com/hashicorp/go-version"
)

// Tag control the tag-it command
type Tag struct {
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
	TagItCmd = "tag-it"
)

// Action execute the `tag-it` command
func (c *Tag) Action([]string) {
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
		fmt.Printf("Release version file (%s) found, but version string has not been detected from '%s'.", c.versionFile, core.DefaultExtractVersion)
		os.Exit(2)
	}
	c.releaseVersion = result[1]
	fmt.Printf("Release version detected: %s (in %s)\n", c.releaseVersion, c.versionFile)

	c.release = core.NewRelease()
	_, err = c.release.CheckVersion(c.releaseVersion)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// Ready to tag it

	if err = c.git.CreateTag(c.releaseVersion); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// Check remote. Create it if missing.

	// Push it

	// create github release in draft mode
}

// Init initialize the check cli commands
func (c *Tag) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(TagItCmd, "Step to tag the code, before build")
	c.versionFile = core.DefaultVersionFile

	var err error
	c.extractVersionRE, err = regexp.Compile(fmt.Sprintf(core.DefaultExtractVersion, version.SemverRegexpRaw))
	kingpin.FatalIfError(err, "Unable to initialize check command")

}
