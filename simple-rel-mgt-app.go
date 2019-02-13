package main

/*  - Simple release process management
    Copyright (C) 2019 clarsonneur@gmail.com

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.

	SEE LICENSE.txt for license details.
*/

import (
	"fmt"
	"simple-relmgt/cmds/checkcmd"
	"simple-relmgt/cmds/draftcmd"
	"simple-relmgt/cmds/releasecmd"
	"simple-relmgt/cmds/statecmd"
	"simple-relmgt/cmds/tagcmd"

	"github.com/alecthomas/kingpin"
	"github.com/forj-oss/forjj-modules/trace"
)

const (
	check     = "check"
	stateIt   = "status"
	draftIt   = "draft-it"
	tagIt     = "tag-it"
	releaseIt = "release-it"
)

var (
	build_branch string
	build_commit string
	build_date string
	build_tag string
)

type simpleRelMgtApp struct {
	app *kingpin.Application

	actionDispatch map[string]func([]string)

	check   checkcmd.CheckCmd
	state   statecmd.StateCmd
	draft   draftcmd.DraftCmd
	tag     tagcmd.TagCmd
	release releasecmd.ReleaseCmd
}

func (a *simpleRelMgtApp) init() {
	gotrace.SetInfo()

	a.app = kingpin.New("simple-relmgt", "Simple release management tool.")

	a.setVersion()

	a.actionDispatch = make(map[string]func([]string))
	a.actionDispatch[check] = a.check.Action
	a.actionDispatch[stateIt] = a.state.Action
	a.actionDispatch[draftIt] = a.draft.Action
	a.actionDispatch[tagIt] = a.tag.Action
	a.actionDispatch[releaseIt] = a.release.Action
}

// setVersion define the current jplugins version.
func (a *simpleRelMgtApp) setVersion() {
	version := "simple-relmgt"

	if PRERELEASE {
		version += " pre-release V" + VERSION
	} else if build_tag == "false" {
		version += " pre-version V" + VERSION
	} else {
		version += " V" + VERSION
	}

	if build_branch != "master" && build_branch != "HEAD" {
		version += fmt.Sprintf(" branch %s", build_branch)
	}
	if build_tag == "false" {
		version += fmt.Sprintf(" - %s - %s", build_date, build_commit)
	}

	a.app.Version(version).Author(AUTHOR)

}
