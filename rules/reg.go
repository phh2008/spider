package rules

import "regexp"

var regTelMobile, _ = regexp.Compile(`(var itm = '\S+)|(var itt = '\S+)`)
