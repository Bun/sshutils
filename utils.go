package sshutils

import (
	"log"
	"regexp"
	"strings"
)

var wildcardExpr = strings.NewReplacer(
	"\\*", ".*",
	"\\?", ".?",
)

// TODO: error reporting
func wildcards(wcs []string) (rs []*regexp.Regexp) {
	for _, wc := range wcs {
		wc := wildcardExpr.Replace(regexp.QuoteMeta(wc))
		if r, err := regexp.Compile(wc); err == nil {
			rs = append(rs, r)
		} else {
			log.Println("Wildcard failed:", err)
		}
	}
	return
}
