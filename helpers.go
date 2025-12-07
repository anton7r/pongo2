package pongo2

import (
	"unicode"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func stripWhitespace(s string, inQuote rune, lastChar rune) (string, rune, rune) {
	buf := getBuffer()
	defer putBuffer(buf)
	pendingSpace := false

	for _, r := range s {
		if inQuote != 0 {
			buf.WriteRune(r)
			if r == inQuote {
				inQuote = 0
			}
			lastChar = r
			continue
		}

		if r == '"' || r == '\'' {
			if pendingSpace {
				if lastChar != '=' {
					buf.WriteRune(' ')
				}
				pendingSpace = false
			}
			inQuote = r
			buf.WriteRune(r)
			lastChar = r
			continue
		}

		if unicode.IsSpace(r) {
			pendingSpace = true
			continue
		}

		if pendingSpace {
			shouldWrite := true
			if lastChar == '>' && r == '<' {
				shouldWrite = false
			}
			if lastChar == '=' || r == '=' {
				shouldWrite = false
			}
			if (lastChar == '"' || lastChar == '\'') && r == '>' {
				shouldWrite = false
			}
			if lastChar == 0 {
				shouldWrite = false
			}

			if shouldWrite {
				buf.WriteRune(' ')
			}
			pendingSpace = false
		}

		buf.WriteRune(r)
		lastChar = r
	}

	if pendingSpace {
		shouldWrite := true
		if lastChar == '>' || lastChar == '=' || lastChar == 0 {
			shouldWrite = false
		}
		if shouldWrite {
			buf.WriteRune(' ')
		}
	}

	return buf.String(), inQuote, lastChar
}
