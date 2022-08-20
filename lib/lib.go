package lib

import (
	"math/rand"
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

func Clamper(floor, ceiling float64) func(float64) float64 {
	clampFunc := func(in float64) float64 {
		if in < floor {
			return floor
		}
		if in > ceiling {
			return ceiling
		}
		return in
	}
	return clampFunc
}

func UniformRange(floor, ceiling float64) float64 {
	if ceiling < floor {
		floor, ceiling = ceiling, floor
	}
	delta := ceiling - floor
	ratio := rand.Float64()
	return floor + ratio*delta
}

func CoinFlip() bool {
	n := rand.Intn(2)
	return n == 0
}

func IntDivmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
