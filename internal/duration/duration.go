package duration

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var commaRegex = regexp.MustCompile(`(\d)[,_](\d)`)
var pluralRegex = regexp.MustCompile(`s$`)
var durationRegex = regexp.MustCompile(`(?i)(-?(?:\d+\.?\d*|\d*\.?\d+)(?:e[-+]?\d+)?)\s*([\p{L}]*)`)

var unitsTable = map[string]time.Duration{
	"nanosecond":  time.Nanosecond,
	"ns":          time.Nanosecond,
	"microsecond": time.Microsecond,
	"µs":          time.Microsecond,
	"us":          time.Microsecond,
	"millisecond": time.Millisecond,
	"ms":          time.Millisecond,
	"second":      time.Second,
	"sec":         time.Second,
	"s":           time.Second,
	"minute":      time.Minute,
	"min":         time.Minute,
	"m":           time.Minute,
	"hour":        time.Hour,
	"hr":          time.Hour,
	"h":           time.Hour,
	"day":         time.Hour * 24,
	"d":           time.Hour * 24,
}

const (
	Hours        Format = "h"
	Minutes             = "m"
	Seconds             = "s"
	Milliseconds        = "ms"
)

type Format string

// ParseDuration converts a human-readable duration to a duration
func ParseDuration(str string, format Format) float32 {
	// Implementation from here: https://github.com/jkroso/parse-duration/blob/master/index.js
	var result float32

	str = commaRegex.ReplaceAllString(str, "$1$2")
	r := durationRegex.FindAllStringSubmatch(str, -1)
	for _, match := range r {
		duration := match[1]
		unit := match[2]

		unitf := unitRatio(unit)

		durationf, err := strconv.ParseFloat(duration, 32)
		if err != nil {
			log.Fatalf("invalid duration specified: %s", match)
		}

		result += float32(durationf) * unitf
	}

	return result / unitRatio(string(format))
}
func unitRatio(str string) float32 {
	ratio, ok := unitsTable[str]
	if ok {
		return float32(ratio)
	} else {
		str = strings.ToLower(str)
		str = pluralRegex.ReplaceAllString(str, "")
		return float32(unitsTable[str])
	}
}
