package util

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

type RegexMatch struct {
	Index   int
	Group   map[string]string
	Input   string
	Context string
}

func FindNamedMatches(regex *regexp.Regexp, str string) []RegexMatch {
	matchIndices := regex.FindAllStringSubmatchIndex(str, -1)
	stringMatches := regex.FindAllStringSubmatch(str, -1)

	var matches []RegexMatch
	var groups []map[string]string

	for _, stringMatch := range stringMatches {
		group := map[string]string{}
		for i, name := range stringMatch {
			group[regex.SubexpNames()[i]] = name
		}
		groups = append(groups, group)
	}

	for _, matchIndex := range matchIndices {
		matches = append(matches, RegexMatch{
			Index:   matchIndex[0],
			Group:   nil,
			Input:   str,
			Context: str[utf8.RuneCountInString(str[:matchIndex[0]]):utf8.RuneCountInString(str[:matchIndex[1]])],
		})
	}

	for i := range matches {
		matches[i].Group = groups[i]
	}

	return matches
}

var commaRegex = regexp.MustCompile(`(\d)[,_](\d)`)
var pluralRegex = regexp.MustCompile(`s$`)
var durationRegex = regexp.MustCompile(`(?i)(-?(?:\d+\.?\d*|\d*\.?\d+)(?:e[-+]?\d+)?)\s*([\p{L}]*)`)

var unitsTable = map[string]float32{
	"nanosecond":  1 / 1e6,
	"ns":          1 / 1e6,
	"microsecond": 1 / 1e3,
	"µs":          1 / 1e3,
	"us":          1 / 1e3,
	"millisecond": 1,
	"ms":          1,
	"second":      1000,
	"sec":         1000,
	"s":           1000,
	"minute":      1000 * 60,
	"min":         1000 * 60,
	"m":           1000 * 60,
	"hour":        1000 * 60 * 60,
	"hr":          1000 * 60 * 60,
	"h":           1000 * 60 * 60,
	"day":         1000 * 60 * 60 * 24,
	"d":           1000 * 60 * 60 * 24,
}

// ParseDuration converts a human-readable duration to ms
func ParseDuration(str, format string) float32 {
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

	return result / unitRatio(format)
}
func unitRatio(str string) float32 {
	ratio, ok := unitsTable[str]
	if ok {
		return ratio
	} else {
		str = strings.ToLower(str)
		str = pluralRegex.ReplaceAllString(str, "")
		return unitsTable[str]
	}
}
