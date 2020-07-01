// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package msgraph

import (
	"io"
	"os"
)

func openConsole(path string) (io.WriteCloser, error) {
	if path == "" {
		path = "/dev/tty"
	}
	return os.OpenFile(path, os.O_RDWR, 0)
}
