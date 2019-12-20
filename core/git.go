package core

import (
	"errors"
	"fmt"
	"net/url"

	gotrace "github.com/forj-oss/forjj-modules/trace"
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

	tag    *plumbing.Reference
	remote *git.Remote
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

	// Fetch repo
	var fetchOptions git.FetchOptions
	fetchOptions.Validate()
	g.repo.Fetch(&fetchOptions)

	// Head commit
	commitID := plumbing.Hash{}
	if id, err := g.repo.Head(); err != nil {
		return err
	} else {
		commitID = id.Hash()
	}
	gotrace.Trace("HEAD is %s", commitID)

	// Get tag if exist
	g.tag, err = g.repo.Tag(name)
	if err != nil && err.Error() != git.ErrTagNotFound.Error() {
		return err
	} else if g.tag != nil {
		if g.tag.Hash() == commitID {
			gotrace.Trace("Tag %s already linked to %s", name, commitID)
			return
		}
		g.repo.DeleteTag(name)
		gotrace.Trace("Updating tag %s to %s", name, commitID)
	} else {
		gotrace.Trace("Creating tag %s to %s", name, commitID)
	}

	g.tag, err = g.repo.CreateTag(name, commitID, nil)
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
		g.remote = r
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
		var user string
		if u, found := options["user"]; found {
			user = u
		} else {
			user = "git"
		}

		URL := ""
		if user != "" {
			URL = user + "@"
		}
		URL += g.host + ":" + g.repoPath
		gitURL := url.URL{
			Host: g.host,
			Path: g.repoPath,
			User: url.User(user),
		}
		// By default, define format user@server:repoPath
		if gitURL.Port() != "" {
			// define format ssh://user@server:port/repoPath
			gitURL.Scheme = "ssh"
			remoteConfig.URLs = []string{gitURL.String()}
		} else {
			remoteConfig.URLs = []string{URL}
		}

	default:
		return errors.New("invalid protocol " + protocol)
	}

	if r, err := g.repo.CreateRemote(&remoteConfig); err != nil {
		return err
	} else {
		g.remote = r
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
	if g == nil {
		return errors.New("git object is nil")
	}

	if g.remote == nil {
		return errors.New("git.remote object is nil")
	}

	if g.tag == nil {
		return errors.New("git.tag object is nil")
	}

	ref := gitconfig.RefSpec(g.tag.Name() + ":" + g.tag.Name())
	pushOptions := git.PushOptions{
		RemoteName: g.remoteName,
		RefSpecs: []gitconfig.RefSpec{
			ref,
		},
	}

	return g.remote.Push(&pushOptions)
}
