package core

import (
	"archive/tar"
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/forj-oss/forjj-modules/trace"
)

// This file manage the relationship between simple-relmgt and github through github-release (https://github.com/aktau/github-release)

// It download the binary from https://github.com/aktau/github-release/releases/download/v{version}/linux-amd64-github-release.tar.bz2

// Github represents the github-release command used by simple-relmgt
type Github struct {
	url            string
	file           string
	version        string
	untarCmd       string
	packageName    string
	packageExtract string
}

const (
	defaultVersion        = "v0.7.2"
	defaultURLPath        = "https://github.com/aktau/github-release/releases/download/%s/%s"
	defaultFilePath       = "github-release"
	defaultFileName       = "linux-amd64-github-release.tar.bz2"
	defaultPackageExtract = "tar -xvjf -"
)

// NewGithub creates the Github object
func NewGithub() (ret *Github) {
	ret = new(Github)

	ret.url = defaultURLPath
	ret.file = defaultFilePath
	ret.version = defaultVersion
	ret.packageName = defaultFileName
	ret.packageExtract = defaultPackageExtract

	return
}

// SetAppVersion define the version of github-release to use.
func (g *Github) SetAppVersion(version string) {
	if g == nil {
		return
	}
	g.version = version
}

// SetURLPath define the URL path where versioned package are stored.
func (g *Github) SetURLPath(urlPath string) (err error) {
	if g == nil {
		return errors.New("Github object is nil")
	}

	if found, _ := regexp.Match("%s.*%s", []byte(urlPath)); !found {
		return fmt.Errorf("%s is an invalid package URL base. It must contains '%%s' twice. The first one will get package version, the second will get package file name", urlPath)
	}
	finalURL := fmt.Sprintf(urlPath, g.version, g.packageName)

	urlTest := new(url.URL)
	if urlTest == nil {
		return fmt.Errorf("Cannot test the URL %s. Unable to allocate url.URL", urlPath)
	} else if _, err := urlTest.Parse(finalURL); err != nil {
		return fmt.Errorf("The URL '%s' is invalid. %s", finalURL, err)
	}

	g.url = urlPath
	return
}

// CheckGithub verify the binary existence and its version.
func (g *Github) CheckGithub() error {
	if g == nil {
		return errors.New("Github object is nil")
	}
	if ok, _ := g.checkGithub(); !ok {
		return g.download()
	}
	return nil
}

// Internal github-release check
// - file found and executable
// - returning version requested.
func (g *Github) checkGithub() (bool, error) {
	if g == nil {
		return false, errors.New("Github object is nil")
	}

	gotrace.Trace("Checking %s ...", g.file)
	info, err := os.Stat(g.file)
	if err != nil {
		return false, err
	}

	mode := info.Mode().Perm()
	if (mode & 0100) == 0 {
		return false, fmt.Errorf("%s is not executable", g.file)
	}

	// Have a relative path to the file
	cmd := path.Clean(g.file)
	if filepath.Dir(cmd) == "." {
		cmd = "./" + g.file
	}

	command := exec.Command(cmd, "--version")
	output, err := command.Output()
	if err != nil {
		return false, err
	}

	if !strings.Contains(string(output), g.version) {
		return false, fmt.Errorf("Version %s not detected. Got %s", g.version, string(output))
	}
	gotrace.Info("OK: Found %s version %s", g.file, g.version)
	return true, nil
}

// Download the github-release file
func (g *Github) download() (err error) {
	if g == nil {
		return errors.New("Github object is nil")
	}

	finalURL := fmt.Sprintf(g.url, g.version, g.packageName)
	gotrace.Trace("Downloading %s...", finalURL)

	// Get the data
	resp, err := http.Get(finalURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unable to download '%s'. bad status: %s", finalURL, resp.Status)
	}

	if err = g.UntarFile(g.file, resp.Body); err != nil {
		return fmt.Errorf("Unable to uncompress tar ball. %s", err)
	}

	if info, err := os.Stat(g.file); err != nil {
		return fmt.Errorf("Unable to find '%s' from '%s'. %s", g.file, g.packageName, err)
	} else {
		mode := info.Mode().Perm()
		if (mode & 0100) == 0 {
			if err = os.Chmod(g.file, 0755); err != nil {
				return fmt.Errorf("Unable to set '%s' as executable. %s", g.file, err)
			}
		}
	}

	if _, err := g.checkGithub(); err != nil {
		return fmt.Errorf("Something is wrong. Expected to have an executable %s version %s. %s", g.file, g.version, err)
	}

	return
}

// UntarFile the wanted file
func (g *Github) UntarFile(file string, r io.Reader) error {

	gzr := bzip2.NewReader(r)

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		if filepath.Base(header.Name) != file {
			continue
		}

		// check the file type
		switch header.Typeflag {

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
