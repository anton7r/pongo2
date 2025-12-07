package pongo2

import "strings"

var escapeReplacer = strings.NewReplacer(
	"&", "&amp;",
	">", "&gt;",
	"<", "&lt;",
	"\"", "&quot;",
	"'", "&#39;",
)
