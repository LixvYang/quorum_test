// Copyright 2019 Google LLC
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd
package crypto

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/term"
)


func PassphrasePromptForUnlock() (string, error) {
	pass, err := readPassphrase("Enter passphrase:")
	if err != nil {
		return "", fmt.Errorf("could not read passphrase:%v",err)
	}

	p := string(pass)
	if p == "" {
		return "",fmt.Errorf("Passphrase can't be blank")
	}
	return p, nil
}

func PassphrasePromptForEncryption() (string, error) {
	pass, err := readPassphrase("Enter passphrase (leave empty to autogenerate a secure one):")
	if err != nil {
		return "", fmt.Errorf("could not read passphrase: %v", err)
	}
	p := string(pass)
	if p == "" {
		var words []string
		for i := 0; i < 10; i++ {
			words = append(words, randomWord())
		}
		p = strings.Join(words,"-")
		// TODO:consider printing this to the terminal, instead of stderr.
		fmt.Fprintf(os.Stderr,"Using the autogenerate passphrase:%q.\n",p)
	} else {
		confirm, err := readPassphrase("Confirm passphrase:")
		if err != nil {
		  return "", fmt.Errorf("could not read passphrase: %v", err)
		}
		if string(confirm) != p {
			return "", fmt.Errorf("passphrase didn't match")
		}
	}
	return p, nil
}

func readPassphrase(prompt string) ([]byte, error) {
	var in,out *os.File
	if runtime.GOOS == "windows" {
		var err error
		in, err = os.OpenFile("CONIN$", os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		defer in.Close()
		out, err = os.OpenFile("CONOUT$", os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
		defer out.Close()
	} else if _, err := os.Stat("/dev/tty"); err == nil {
		tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		defer tty.Close()
		in, out = tty, tty
	} else {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return nil, fmt.Errorf("standard input is not a terminal, and /dev/tty is not available: %v", err)
		}
		in, out = os.Stdin, os.Stderr
	}
	fmt.Fprintf(out, "%s ", prompt)
	// Use CRLF to work around an apparent bug in WSL2's handling of CONOUT$.
	// Only when running a Windows binary from WSL2, the cursor would not go
	// back to the start of the line with a simple LF. Honestly, it's impressive
	// CONIN$ and CONOUT$ even work at all inside WSL2.
	defer fmt.Fprintf(out, "\r\n")
	return term.ReadPassword(int(in.Fd()))
}