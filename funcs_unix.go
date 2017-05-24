// +build !windows

package main

import "strings"

func filter(path string) int8 {

	lowPath := strings.ToLower(path)

	suffixList := settings.EncSuffixList

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
