package tagcmd

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"simple-relmgt/core"

	"github.com/alecthomas/kingpin"
	version "github.com/hashicorp/go-version"
)

// Tag control the tag-it command
type Tag struct {
	cmd          *kingpin.CmdClause
	proto        *string
	user         *string
	pass         *string
	remoteName   *string
	removeRemote *bool

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
func (c *Tag) Action([]string) int {
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
		return 1
	}

	result := c.extractVersionRE.FindStringSubmatch(string(data))
	if result == nil {
		fmt.Printf("Release version file (%s) found, but version string has not been detected from '%s'.", c.versionFile, core.DefaultExtractVersion)
		return 1
	}
	c.releaseVersion = result[1]
	fmt.Printf("Release version detected: %s (in %s)\n", c.releaseVersion, c.versionFile)

	c.release = core.NewRelease()
	_, err = c.release.CheckVersion(c.releaseVersion)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	// Ready to tag it

	if err = c.git.CreateTag(c.releaseVersion); err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	// Check remote. Create it if missing.
	c.git.SetRemote(c.config.Yaml.Upstream.Name, c.config.Yaml.Upstream.Protocol, c.config.Yaml.Upstream.Server, c.config.Yaml.Upstream.RepoPath)

	options := make(core.GitRemoteConfig)

	options.Set("username", *c.user)
	options.Set("password", *c.pass)
	options.SetIfNotContains("remote-name", *c.remoteName, "ci-upstream")
	options.SetIf("auto-remove-remote", "false", !*c.removeRemote)
	options.Set("protocol", *c.proto)

	if err = c.git.CreateRemote(options); err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	defer c.git.CleanRemote()

	// Push it
	if err = c.git.PushTag(); err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	// create github release in draft mode
	/*	if err = c.github.CreateDraft(c.releaseVersion) ; err != nil  {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}*/
	return 0
}

// Init initialize the check cli commands
func (c *Tag) Init(app *kingpin.Application) {
	if c == nil || app == nil {
		return
	}
	c.cmd = app.Command(TagItCmd, "Step to tag the code, and push it, before build")
	c.proto = c.cmd.Flag("protocol", "url protocol to use. by default, uses https. Supported ones are https/http/ssh").Envar("https").String()
	c.user = c.cmd.Flag("user", "User name to define in the upstream url (https/https/ssh). Detect GIT_USER").Envar("GIT_USER").String()
	c.pass = c.cmd.Flag("password", "Password to provide, like a github token. (https/http). Detect GIT_PASSWORD").Envar("GIT_PASSWORD").String()
	c.remoteName = c.cmd.Flag("remote-name", "Remote name to manage. By default, uses 'ci-upstream'").Default("ci-upstream").String()
	c.removeRemote = c.cmd.Flag("auto-remove", "Remove the created remote, when done.").Default("true").Bool()

	c.versionFile = core.DefaultVersionFile

	var err error
	c.extractVersionRE, err = regexp.Compile(fmt.Sprintf(core.DefaultExtractVersion, version.SemverRegexpRaw))
	kingpin.FatalIfError(err, "Unable to initialize check command")

}
