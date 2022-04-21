package lib

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func Atoi(s string) int {
	if num, err := strconv.Atoi(s); err != nil {
		return 0
	} else {
		return num
	}
}

func SplitOnce(s, delim string) (left, right string) {
	if !strings.Contains(s, delim) {
		return s, ""
	}
	if delim == "" {
		return "", s
	}
	pivot := strings.Index(s, delim)
	offset := pivot + len(delim)
	left = s[:pivot]
	right = s[offset:]
	return
}

func MustGetCallerDir() string {
	skip := 1 // skip current frame to get the caller's directory
	if _, file, _, ok := runtime.Caller(skip); !ok {
		panic("failed to get current runtime file")
	} else {
		return filepath.Dir(file)
	}
}
