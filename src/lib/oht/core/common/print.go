//
// Superior - Toolkit for the Go programming language
// Available at http://github.com/liamka/Printm
//
// Copyright © Kirill Kotikov <liamka@me.com>.
// Superior is open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT)
// See README.md for details.
//

package common

import (
	"fmt"
)

const (
	txtblk = "0;30m" // Black - Regular
	txtred = "0;31m" // Red
	txtgrn = "0;32m" // Green
	txtylw = "0;33m" // Yellow
	txtblu = "0;34m" // Blue
	txtpur = "0;35m" // Purple
	txtcyn = "0;36m" // Cyan
	txtwht = "0;37m" // White

	bldblk = "1;30m" // Black - Bold
	bldred = "1;31m" // Red
	bldgrn = "1;32m" // Green
	bldylw = "1;33m" // Yellow
	bldblu = "1;34m" // Blue
	bldpur = "1;35m" // Purple
	bldcyn = "1;36m" // Cyan
	bldwht = "1;37m" // White

	undblk = "4;30m" // Black - Underline
	undred = "4;31m" // Red
	undgrn = "4;32m" // Green
	undylw = "4;33m" // Yellow
	undblu = "4;34m" // Blue
	undpur = "4;35m" // Purple
	undcyn = "4;36m" // Cyan
	undwht = "4;37m" // White

	bakblk = "40m" // Black - Background
	bakred = "41m" // Red
	bakgrn = "42m" // Green
	bakylw = "43m" // Yellow
	bakblu = "44m" // Blue
	bakpur = "45m" // Purple
	bakcyn = "46m" // Cyan
	bakwht = "47m" // White

	txtrst = "0m" // Reset
)

func Print(str string, t string, c string) {
	x := "\x1b["

	if t == "" {
		t = "normal"
	}

	if c == "" {
		c = "white"
	}

	checkType(t, "type")
	checkType(c, "color")

	if t == "normal" {
		switch c {
		case "black":
			x = x + txtblk
		case "red":
			x = x + txtred
		case "green":
			x = x + txtgrn
		case "yellow":
			x = x + txtylw
		case "blue":
			x = x + txtblu
		case "purple":
			x = x + txtpur
		case "cyan":
			x = x + txtcyn
		case "white":
			x = x + txtwht
		}
	} else if t == "bold" {
		switch c {
		case "black":
			x = x + bldblk
		case "red":
			x = x + bldred
		case "green":
			x = x + bldgrn
		case "yellow":
			x = x + bldylw
		case "blue":
			x = x + bldblu
		case "purple":
			x = x + bldpur
		case "cyan":
			x = x + bldcyn
		case "white":
			x = x + bldwht
		}
	} else if t == "under" {
		switch c {
		case "black":
			x = x + undblk
		case "red":
			x = x + undred
		case "green":
			x = x + undgrn
		case "yellow":
			x = x + undylw
		case "blue":
			x = x + undblu
		case "purple":
			x = x + undpur
		case "cyan":
			x = x + undcyn
		case "white":
			x = x + undwht
		}
	} else if t == "background" {
		switch c {
		case "black":
			x = x + bakblk
		case "red":
			x = x + bakred
		case "green":
			x = x + bakgrn
		case "yellow":
			x = x + bakylw
		case "blue":
			x = x + bakblu
		case "purple":
			x = x + bakpur
		case "cyan":
			x = x + bakcyn
		case "white":
			x = x + bakwht
		}
	}
	x = x + "%s\x1b[" + txtrst + "\n"
	fmt.Printf(x, str) // return
}
