package util

import (
	"regexp"
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

	for i, _ := range matches {
		matches[i].Group = groups[i]
	}

	return matches
}
