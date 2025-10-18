//go:build !windows

package console

import "syscall"

// setNonblock temporarily sets stdin to non-blocking, only used on Unix
func setNonblock(fd int, nonblocking bool) error {
	return syscall.SetNonblock(fd, nonblocking)
}
