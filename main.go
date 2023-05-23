package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"github.com/webview/webview"
)

var (
	cmd          *exec.Cmd
	serverKilled = false
)

func main() {
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("HTTP Toolkit")
	w.SetSize(1366, 768, webview.HintNone)
	w.SetHtml("Setting up server..<script>window.onload=()=>_onload()</script>")

	w.Bind("_onload", func() {
		go func() {
			if err := setupServer(); err != nil {
				w.Dispatch(func() {
					w.SetHtml("error while setting up server: <pre>" + err.Error() + "</pre>")
				})
			} else {
				w.Dispatch(func() {
					w.Navigate("https://juby-httptoolkit.github.io")
				})

				out, err := startServer()
				serverKilled = true
				w.Dispatch(func() {
					if err != nil {
						w.SetHtml("server error: <pre>" + err.Error() + "\n\n" + string(out) + "</pre>")
					} else {
						w.SetHtml("server stopped <pre>" + string(out) + "</pre>")
					}
				})
			}
		}()
	})

	//goland:noinspection ALL
	if platform == "win32" {
		w.Bind("prompt", func(_ string) string {
			result, err := cfdutil.ShowOpenFileDialog(cfd.DialogConfig{})
			if err != cfd.ErrorCancelled {
				fmt.Println(err)
			}
			return result
		})
	}

	w.Run()

	if !serverKilled && cmd != nil {
		req, err := http.NewRequest(
			"POST",
			"http://127.0.0.1:45457",
			strings.NewReader(`{"operationName":"Shutdown","query":"mutation Shutdown { shutdown }","variables":{}}`),
		)
		if err == nil {
			req.Header.Set("content-type", "application/json")
			req.Header.Set("origin", "https://app.httptoolkit.tech")
			c := &http.Client{Timeout: 3 * time.Second}
			_, err = c.Do(req)
			if err == nil {
				fmt.Println("soft shutdown")
			}
		}
	}
}

func startServer() ([]byte, error) {
	cmd = exec.Command("bin/node", "bin/run", "start")
	cmd.Dir = "httptoolkit-server"
	cmd.Env = append(
		os.Environ(),
		"HTTPTOOLKIT_SERVER_BINPATH="+binPath,
		"NODE_SKIP_PLATFORM_CHECK=1",
		`NODE_OPTIONS="--max-http-header-size=102400 --insecure-http-parser"`,
	)
	hideWindow()
	return cmd.CombinedOutput()
}
