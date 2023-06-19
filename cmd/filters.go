package cmd

import (
	"regexp"

	"github.com/IGLOU-EU/go-wildcard"
)

type filterFunc func(s string) (excluded bool)

type filters []filterFunc

func (f filters) excluded(s string) (excluded bool) {
	for _, ff := range f {
		if ff(s) {
			return true
		}
	}
	return false
}

func matchFilter(patterns []string) filterFunc {
	return func(s string) (excluded bool) {
		for _, p := range patterns {
			if wildcard.MatchSimple(p, s) {
				return false
			}
		}
		return true
	}
}

func excludeFilter(patterns []string) filterFunc {
	return func(s string) (excluded bool) {
		for _, p := range patterns {
			if wildcard.MatchSimple(p, s) {
				return true
			}
		}
		return false
	}
}

func excludeRegexpFilter(patterns []*regexp.Regexp) filterFunc {
	return func(s string) (excluded bool) {
		for _, rx := range patterns {
			if rx.MatchString(s) {
				return true
			}
		}
		return false
	}
}
