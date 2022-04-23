package lib

import (
	"math"
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
			return in
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

func ChooseString(options []string) string {
	size := len(options)
	n := rand.Intn(size)
	return options[n]
}

func ChooseBool() bool {
	n := rand.Intn(2)
	return n == 0
}

// HSVtoRGB converts from HSV tuple to RGB colour space.
func HSVtoRGB(h, s, v float64) (r, g, b float64) {
	// Adapted from Python stdlib
	// https://github.com/python/cpython/blob/3.10/Lib/colorsys.py

	if s == 0 {
		return v, v, v
	}
	i := int(math.Floor(h * 6))
	f := h*6 - 1
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	switch i % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	case 5:
		return v, p, q
	default:
		// Cannot get here
		return 0, 0, 0
	}
}
