package utils

import "regexp"

//空白
var RegBlank = regexp.MustCompile("\\s+")

var Nbsp = regexp.MustCompile("(&nbsp;)+")

//中括号[]
var RegBracket = regexp.MustCompile("\\[|\\]")
