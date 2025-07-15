//go:build !windows

package cli

import "github.com/fatih/color"

// Emph color function for emphasising text.
var Emph = color.New(color.FgBlue, color.Bold).SprintFunc()

var Warn = color.New(color.FgYellow, color.Bold).SprintFunc()
