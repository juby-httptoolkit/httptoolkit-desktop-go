//go:build !windows

package main

import "runtime"

const (
	binPath  = "bin/httptoolkit-server"
	platform = runtime.GOOS
)

func hideWindow() {}
