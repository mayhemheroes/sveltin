/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package utils contains utility function for errors, node package manager, text, progressbar, github.
package utils

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/sveltinio/sveltin/common"
)

// ExitIfExists panics on os.Exit(1) if the given path exists.
func ExitIfExists(fs afero.Fs, path string) {
	if common.DirExists(fs, path) {
		name := filepath.Base(path)
		err := fmt.Errorf("%s already exists in %s. Remove it to re-create", name, path)
		log.Fatalf("\x1b[31;1m✘ %s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	}
}

// ExitIfError panics on os.Exit(1) if error.
func ExitIfError(err error) {
	if err == nil {
		return
	}
	log.Fatalf("\x1b[31;1m✘ %s\x1b[0m\n", fmt.Sprintf("error: %s", err))
}

// IsError returns true if error is not nil.
// If showMessage is true it prints out a warning with the error message.
func IsError(err error, showMessage bool) bool {
	if err != nil && showMessage {
		log.Printf("\x1b[33;1;33m %s\x1b[0m\n", fmt.Sprintf("warning: %s", err))
		return true
	} else if err != nil {
		return true
	}
	return false
}
