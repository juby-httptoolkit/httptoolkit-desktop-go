//go:build !windows

package main

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

void webkit_fixes(void *window) {
	GtkWidget *wv = gtk_bin_get_child((GtkBin *) window);
	WebKitWebContext *ctx = webkit_web_view_get_context(WEBKIT_WEB_VIEW(wv));
	WebKitSecurityManager *sm = webkit_web_context_get_security_manager(ctx);
	webkit_security_manager_register_uri_scheme_as_secure(sm, "http");
	webkit_security_manager_register_uri_scheme_as_secure(sm, "ws");
}
*/
import "C"

import (
	"runtime"
	"unsafe"
)

const (
	binPath  = "bin/httptoolkit-server"
	platform = runtime.GOOS
)

func hideWindow() {}

func webkitFixes(window unsafe.Pointer) {
	C.webkit_fixes(window)
}
