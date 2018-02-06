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
	"os"
)

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func checkType(a string, b string) {

	var strList []string
	var mess string

	if b == "type" {
		strList = []string{"normal", "bold", "under", "background"}
		mess = "Possible value for type: normal, bold, under, background"
	} else if b == "color" {
		strList = []string{"black", "red", "green", "yellow", "blue", "purple", "cyan", "white"}
		mess = "Possible value for color: black, red, green, yellow, blue, purple, cyan, white"
	}

	if !stringInSlice(a, strList) {
		fmt.Println(mess)
		os.Exit(1)
	}

}
