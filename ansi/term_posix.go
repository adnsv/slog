// +build !windows

package ansi

import "os"

func EnableVT(f *os.File) (bool, func()) {
	return true, func() {}
}
