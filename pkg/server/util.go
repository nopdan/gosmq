package server

import (
	"os/exec"
	"runtime"

	"github.com/nopdan/gosmq/pkg/util"
)

var logger = util.Logger

func openBrowser(url string) {
	var name string
	switch runtime.GOOS {
	case "windows":
		name = "explorer"
	case "linux":
		name = "xdg-open"
	default:
		name = "open"
	}
	cmd := exec.Command(name, url)
	cmd.Start()
}
