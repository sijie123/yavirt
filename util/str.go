package util

import "strings"

func MergeStrings(a, b []string) []string {
	var dic = make(map[string]struct{})
	var strs = make([]string, 0, len(a)+len(b))

	var append = func(ar []string) {
		for _, k := range ar {
			if _, exists := dic[k]; exists {
				continue
			}
			dic[k] = struct{}{}
			strs = append(strs, k)
		}
	}

	append(a)
	append(b)

	return strs
}

func PartRight(str, substr string) (string, string) {
	switch i := strings.LastIndex(str, substr); {
	case i < 0:
		return "", str
	case i >= len(str)-1:
		return str[:i], ""
	default:
		return str[:i], str[i+1:]
	}
}
