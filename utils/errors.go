/**
 * Copyright © 2021 Mirco Veltri <github@mircoveltri.me>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package utils ...
package utils

import (
	"fmt"

	jww "github.com/spf13/jwalterweatherman"
)

// ExitIfError panics on os.Exit(1) if error.
func ExitIfError(err error) {
	if err == nil {
		return
	}
	jww.FATAL.Fatalf("\x1b[31;1m✘ %s\x1b[0m\n", fmt.Sprintf("error: %s", err))
}

// IsError returns true if error is not nil.
// If showMessage is true it prints out a warning with the error message.
func IsError(err error, showMessage bool) bool {
	if err != nil && showMessage {
		jww.FEEDBACK.Printf("\x1b[33;1;33m %s\x1b[0m\n", fmt.Sprintf("warning: %s", err))
		return true
	} else if err != nil {
		return true
	}
	return false
}
