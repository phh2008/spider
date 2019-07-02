package utils

import "regexp"

var RegBlank = regexp.MustCompile("\\s+")

var Nbsp = regexp.MustCompile("(&nbsp;)+")
