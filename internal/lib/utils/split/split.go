package split

import "strings"

func SplitMsg(s string) (l []string) {
	l = make([]string, 0)
	for _, el := range strings.Split(s, ",") {
		l = append(l, strings.TrimSpace(el))
	}
	return
}
