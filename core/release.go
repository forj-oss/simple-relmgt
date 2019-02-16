package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

// Release reprensent a release file
type Release struct {
	file          string
	fileType      string
	publishDateRE *regexp.Regexp
	date          time.Time
}

const (
	defaultReleaseFile   = "releases/release-%s.md"
	defaultPublishDateRE = "^ *date *: *"
	defaultDateRE        = `[0-9]{4}/[0-9]{2}/[0-9]{2}`
	timeLayout           = "2006/01/02"
)

// NewRelease create a release file object
func NewRelease() (ret *Release) {
	ret = new(Release)

	ret.fileType = defaultReleaseFile
	// (?m) => Match per line
	ret.publishDateRE, _ = regexp.Compile(`(?mi)` + defaultPublishDateRE + "(" + defaultDateRE + ")")

	return
}

// CheckVersion return a status/error related to the release file name and expected content.
func (r *Release) CheckVersion(version string) (_ int, _ error) {
	if r == nil {
		return
	}

	r.file = fmt.Sprintf(r.fileType, version)

	// status 3: File existence
	if fi, err := os.Stat(r.file); err != nil {
		return 3, fmt.Errorf("Issue accessing %s. %s", r.file, err)
	} else if !fi.Mode().IsRegular() {
		return 3, fmt.Errorf("%s is not a regular file. %s", r.file, err)
	}

	data, err := ioutil.ReadFile(r.file)

	if err != nil {
		return 3, fmt.Errorf("Unable to read %s. %s", r.file, err)
	}

	// Status 2: No date found
	found := r.publishDateRE.FindStringSubmatch(string(data))
	if found == nil {
		return 2, fmt.Errorf("Unable to find the publish date at tag '%s' from %s", defaultPublishDateRE, r.file)
	}

	// Status 1: Date is newer
	r.date, err = time.Parse(timeLayout, found[1])
	if time.Since(r.date).Truncate(time.Hour) == 0 {
		return 1, fmt.Errorf("The release %s is currently planned for %s. Not ready now", version, r.date)
	}

	// Status 0: Ready
	return
}
