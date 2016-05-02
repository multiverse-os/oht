// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// Package utils contains internal helper functions for go-ethereum commands.
package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"./../../common"
	"github.com/peterh/liner"
)

var (
	interruptCallbacks = []func(os.Signal){}
)

func openLogFile(Datadir string, filename string) *os.File {
	path := common.AbsolutePath(Datadir, filename)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening log file '%s': %v", filename, err))
	}
	return file
}

func PromptConfirm(prompt string) (bool, error) {
	var (
		input string
		err   error
	)
	prompt = prompt + " [y/N] "

	// if liner.TerminalSupported() {
	// 	fmt.Println("term")
	// 	lr := liner.NewLiner()
	// 	defer lr.Close()
	// 	input, err = lr.Prompt(prompt)
	// } else {
	fmt.Print(prompt)
	input, err = bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println()
	// }

	if len(input) > 0 && strings.ToUpper(input[:1]) == "Y" {
		return true, nil
	} else {
		return false, nil
	}

	return false, err
}

func PromptPassword(prompt string, warnTerm bool) (string, error) {
	if liner.TerminalSupported() {
		lr := liner.NewLiner()
		defer lr.Close()
		return lr.PasswordPrompt(prompt)
	}
	if warnTerm {
		fmt.Println("!! Unsupported terminal, password will be echoed.")
	}
	fmt.Print(prompt)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	input = strings.TrimRight(input, "\r\n")
	fmt.Println()
	return input, err
}

// Fatalf formats a message to standard error and exits the program.
// The message is also printed to standard output if standard error
// is redirected to a different file.
func Fatalf(format string, args ...interface{}) {
	w := io.MultiWriter(os.Stdout, os.Stderr)
	outf, _ := os.Stdout.Stat()
	errf, _ := os.Stderr.Stat()
	if outf != nil && errf != nil && os.SameFile(outf, errf) {
		w = os.Stderr
	}
	fmt.Fprintf(w, "Fatal: "+format+"\n", args...)
	logger.Flush()
	os.Exit(1)
}

func StartOht(oht *oht.OnionHashTable) {
	log.Println("Starting", oht.Name())
	if err := oht.Start(); err != nil {
		Fatalf("Error starting Oht: %v", err)
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		defer signal.Stop(sigc)
		<-sigc
		log.Println("Got interrupt, shutting down...")
		go oht.Stop()
		logger.Flush()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.Println("Already shutting down, please be patient.")
				log.Println("Interrupt", i-1, "more times to induce panic.")
			}
		}
		log.Println("Force quitting: this might not end so well.")
		panic("boom")
	}()
}
