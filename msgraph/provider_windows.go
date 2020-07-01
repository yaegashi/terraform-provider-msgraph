package msgraph

import (
	"io"
	"os"
)

func openConsole(path string) (io.WriteCloser, error) {
	if path == "" {
		// XXX: I don't know if this actually works...
		path = "CON"
	}
	return os.OpenFile(path, os.O_RDWR, 0)
}
