//go:build windows
// +build windows

package main

import "syscall"

// Shared Windows DLL instances
var (
	user32Common = syscall.NewLazyDLL("user32.dll")
	shcore       = syscall.NewLazyDLL("shcore.dll")
	gdi32        = syscall.NewLazyDLL("gdi32.dll")
)
