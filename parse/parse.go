package parse

import (
	"regexp"
)

const (
	RegQuotedStr = `"[^"]*"`
	RegVarStr    = `(\w+)=(` + RegQuotedStr + `|[^ ]+)`
)

var (
	RegQuoted = regexp.MustCompile(RegQuotedStr)
	RegVar    = regexp.MustCompile(RegVarStr)
)

func compileRegexp(str string) *regexp.Regexp {
	reg, err := regexp.Compile(str)
	if err != nil {
		panic(err)
	}

	return reg
}

func ParseLine(str string) map[string]string {
	res := RegVar.FindAllStringSubmatch(str, -1)
	if res == nil {
		return nil
	}

	outMap := make(map[string]string)

	for _, r := range res {
		val := r[2]

		// Remove quotes
		if RegQuoted.MatchString(val) {
			val = val[1 : len(val)-1]
		}

		outMap[r[1]] = val
	}

	return outMap
}
