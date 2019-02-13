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

import "github.com/alecthomas/kingpin"

type simpleRelMgtApp struct {
	app *kingpin.Application
}

func (a *simpleRelMgtApp) init() {
	gotrace.SetInfo()

	a.app = kingpin.New("simple-relmgt", "Simple release management tool.")

	a.setVersion()
}

// setVersion define the current jplugins version.
func (a *simpleRelMgtApp) setVersion() {
	version := "jplugins"

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
