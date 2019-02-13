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
	"os"

	"github.com/alecthomas/kingpin"
)

// App is the application core struct
var App simpleRelMgtApp

func main() {
	App.init()

	switch kingpin.MustParse(App.app.Parse(os.Args[1:])) {
	// Register user
	case App.draftIt.cmd.FullCommand():
		App.draftIt.DoDraftIt()
	case App.check.cmd.FullCommand():
		App.check.DoCheck()
	case App.tagIt.cmd.FullCommand():
		App.tagIt.DoTagIt()
	case App.releaseIt.cmd.FullCommand():
		App.releaseIt.DoReleaseIt()
	}
}
