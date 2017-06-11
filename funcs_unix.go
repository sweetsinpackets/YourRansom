// +build !windows

package main

import "strings"

func filter(path string) int8 {

	lowPath := strings.ToLower(path)

	suffixList := settings.EncSuffixList
	prefixList := []string{"/dev"}

	for _, prefix := range prefixList {
		if strings.HasPrefix(lowPath, prefix) {
			return 0
		}
	}

	for _, suffix := range suffixList {
		if strings.HasSuffix(lowPath, suffix) {
			return 1
		}
	}

	return 2
}

func fileFilter(path string) int8 {
	return filter(path)
}
