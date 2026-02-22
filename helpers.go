package main

import (
	"bytes"
	"unicode"
)

func FormatTitle(title string) string {
	var buf bytes.Buffer
	for _, v := range title {
		buf.WriteRune(unicode.ToUpper(v))
		if !unicode.IsSpace(v) {
			buf.WriteRune(' ')
		}
	}
	return buf.String()
}
