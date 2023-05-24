//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

const (
	binPath  = "bin/httptoolkit-server.cmd"
	platform = "win32"
)

func hideWindow() {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func webkitFixes(_ unsafe.Pointer) {}
