//go:build windows

package cli

import "fmt"

// Emph Color function for emphasising text.
var Emph = func(a ...any) string {
	return fmt.Sprint(a...)
}

var Warn = func(a ...any) string {
	return fmt.Sprint(a...)
}
