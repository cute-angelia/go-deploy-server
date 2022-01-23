package utils

import "os"

func IsWindows() bool {
	// return runtime.GOOS == "windows"
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
