// +build !windows

package ansi

import "os"

func EnableVT(f *os.File) (ok bool, cleanup func()) {
	return true, func() {}
}
