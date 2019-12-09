package main

import (
	"fmt"
	"regexp"
	"time"
)

type Pattern struct {
	regex   *regexp.Regexp
	pattern string
}

func ParseDateTime(input string) (time.Time, error) {

	var patterns []Pattern

	buildPatterns := func() []Pattern {
		if len(patterns) == 0 {
			regexes := map[string]string{
				"^\\d{1,2}/\\d{1,2}/\\d{4}$":                                       "1/2/2006",
				"^\\d{1,2}/\\d{1,2}/\\d{2}$":                                       "1/2/06",
				"^\\d{1,2}/\\d{1,2}/\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2}$":            "1/2/06 15:4:5",
				"^\\d{1,2}/\\d{1,2}/\\d{4} \\d{1,2}:\\d{1,2}:\\d{1,2}$":            "1/2/2006 15:4:5",
				"^\\d{1,2}/\\d{1,2}/\\d{4} \\d{1,2}:\\d{1,2}$":                     "1/2/2006 15:4",
				"^\\d{1,2}/\\d{1,2}/\\d{4} \\d{1,2}:\\d{1,2} [aApP][mM]$":          "1/2/2006 3:4 PM",
				"^\\d{1,2}/\\d{1,2}/\\d{4} \\d{1,2}:\\d{1,2}:\\d{1,2} [aApP][mM]$": "1/2/2006 3:4:5 PM",

				// year first
				"^\\d{4}/\\d{1,2}/\\d{1,2}$":                                       "2006/1/2",
				"^\\d{4}/\\d{1,2}/\\d{1,2} \\d{1,2}:\\d{1,2}$":                     "2006/1/2 15:4",
				"^\\d{4}/\\d{1,2}/\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}$":            "2006/1/2 15:4:5",
				"^\\d{4}/\\d{1,2}/\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2} [aApP][mM]$": "2006/1/2 3:4:5 PM",

				"^\\d{4}-\\d{1,2}-\\d{1,2}$":                                       "2006-1-2",
				"^\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}$":                     "2006-1-2 15:4",
				"^\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}$":            "2006-1-2 15:4:5",
				"^\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2} [aApP][mM]$": "2006-1-2 3:4:5 PM",

				// yyyy.dd.mm
				"^\\d{4}\\.\\d{1,2}\\.\\d{1,2}$": "2006.1.2",

				//   dd.mm.yyyy
				"^\\d{1,2}\\.\\d{1,2}\\.\\d{4}$":                            "2.1.2006",
				"^\\d{1,2}\\.\\d{1,2}\\.\\d{4} \\d{1,2}:\\d{1,2}$":          "2.1.2006 15:4",
				"^\\d{1,2}\\.\\d{1,2}\\.\\d{4} \\d{1,2}:\\d{1,2}:\\d{1,2}$": "2.1.2006 15:4:5",

				//   dd.mm.yy
				"^\\d{1,2}\\.\\d{1,2}\\.\\d{2}$":                            "2.1.06",
				"^\\d{1,2}\\.\\d{1,2}\\.\\d{2} \\d{1,2}:\\d{1,2}$":          "2.1.06 15:4",
				"^\\d{1,2}\\.\\d{1,2}\\.\\d{2} \\d{1,2}:\\d{1,2}:\\d{1,2}$": "2.1.06 15:4:5",
			}

			for regex, p := range regexes {
				re := regexp.MustCompile(regex)
				pair := Pattern{regex: re, pattern: p}
				patterns = append(patterns, pair)
			}
			return patterns
		}

		return patterns
	}

	buildPatterns()
	format := ""
	for _, regex := range patterns {
		if regex.regex.MatchString(input) {
			format = regex.pattern
		}
	}
	fmt.Printf("Input %s -> Pattern: %s\n", input, format)
	if format == "" {
		panic("No pattern")
	}
	t, err := time.Parse(format, input)
	return t, err
}
