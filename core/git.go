package core

import (
	"errors"
	"fmt"

	"github.com/forj-oss/forjj-modules/trace"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Git struct {
	repoPath string
	repo     *git.Repository
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
