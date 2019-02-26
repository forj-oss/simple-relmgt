package core

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/forj-oss/forjj-modules/trace"
	git "gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Git struct {
	repo *git.Repository

	remoteName   string
	removeRemote bool

	protocol string
	host     string
	repoPath string

	
}

const (
	defaultRepo = "."
)

// NewGit creates the internal GIT object
func NewGit() (ret *Git) {
	ret = new(Git)

	ret.repoPath = defaultRepo

	return
}

func (g *Git) SetRemote(name, protocol, host, repoPath string) {
	if g == nil {
		return
	}
	g.remoteName = name
	g.protocol = protocol
	g.host = host
	g.repoPath = repoPath
}

// OpenRepo open the GIT repo
func (g *Git) OpenRepo() (err error) {

	g.repo, err = git.PlainOpen(g.repoPath)
	if err != nil {
		return fmt.Errorf("%s is not a valid GIT repository. %s", g.repoPath, err)
	}

	if _, err := g.repo.Worktree(); err != nil {
		return fmt.Errorf("Unable to open %s. %s", g.repoPath, err)
	}
	return
}

// CheckTag verify if the tag `name` exist in the default remote.
func (g *Git) CheckTag(name string) (found bool, _ error) {
	if g == nil {
		return false, errors.New("git object is nil")
	}

	var fetchOptions git.FetchOptions
	fetchOptions.Validate()
	g.repo.Fetch(&fetchOptions)

	tagrefs, err := g.repo.Tags()
	if err != nil {
		return false, err
	}

	err = tagrefs.ForEach(func(t *plumbing.Reference) (_ error) {
		if t.Name().String() == "refs/tags/"+name {
			found = true
		}
		return
	})
	return
}

// CreateTag create the requested tag name.
func (g *Git) CreateTag(name string) (err error) {
	if g == nil {
		return errors.New("git object is nil")
	}

	// create/update the tag
	var fetchOptions git.FetchOptions
	fetchOptions.Validate()
	g.repo.Fetch(&fetchOptions)

	tagrefs, err := g.repo.Tags()
	if err != nil {
		return err
	}

	commitID := plumbing.Hash{}
	if id, err := g.repo.Head(); err != nil {
		return err
	} else {
		commitID = id.Hash()
	}

	gotrace.Trace("HEAD is %s", commitID)

	err = tagrefs.ForEach(func(t *plumbing.Reference) (_ error) {
		if t.Name().String() == "refs/tags/"+name {
			if t.Hash() != commitID {
				g.repo.DeleteTag(name)
			}
		}
		return
	})

	g.repo.CreateTag(name, commitID, nil)
	return
}

// CreateRemote create or update a remote.
func (g *Git) CreateRemote(options GitRemoteConfig) (_ error) {
	g.remoteName = options.Get("remote-name", g.remoteName, "ci-upstream")

	if g.host == "" || g.repoPath == "" {
		return errors.New("Unable to create a remote without host/repo-path setup in 'release-mgt.yaml'")
	}

	if r, err := g.repo.Remote(g.remoteName); r == nil && err.Error() != git.ErrRemoteNotFound.Error() {
		return err
	} else if r != nil {
		g.removeRemote = false
		gotrace.Trace("Remote '%s' already exist. Not recreated.", g.remoteName)
		return nil
	}

	if v, found := options["auto-remove-remote"]; !found || v == "true" {
		g.removeRemote = true
	}

	remoteConfig := gitconfig.RemoteConfig{
		Name: g.remoteName,
	}

	switch protocol := options.Get("protocol", g.protocol, "https"); protocol {
	case "https", "http":
		var user *url.Userinfo
		if u, foundUser := options["user"]; foundUser {
			if p, foundPassword := options["password"]; foundPassword {
				user = url.UserPassword(u, p)
			} else {
				user = url.User(u)
			}
		}

		gitURL := url.URL{
			Scheme: protocol,
			Host:   g.host,
			Path:   g.repoPath,
			User:   user,
		}
		remoteConfig.URLs = []string{gitURL.String()}
	case "ssh":
		var user *url.Userinfo
		if u, found := options["user"]; found {
			user = url.User(u)
		}

		gitURL := url.URL{
			Host: g.host,
			Path: g.repoPath,
			User: user,
		}
		// By default, define format user@server:repoPath
		if gitURL.Port() != "" {
			// define format ssh://user@server:port/repoPath
			gitURL.Scheme = "ssh"
		}
		remoteConfig.URLs = []string{gitURL.String()}
	default:
		return errors.New("invalid protocol " + protocol)
	}

	if _, err := g.repo.CreateRemote(&remoteConfig); err != nil {
		return err
	}
	gotrace.Trace("Remote '%s' created.", g.remoteName)

	return
}

// CleanRemote remove a remote previously created.
func (g *Git) CleanRemote() (_ error) {
	if g == nil {
		return errors.New("Git object is nil")
	}
	if !g.removeRemote {
		return
	}

	r, err := g.repo.Remote(g.remoteName)
	if r != nil {
		g.repo.DeleteRemote(g.remoteName)
		gotrace.Trace("Remote '%s' removed.", g.remoteName)
	}

	return err
}

// PushTag push a tag to the remote repository
func (g *Git) PushTag() (err error) {


	return
}