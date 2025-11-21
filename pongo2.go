package pongo2

//go:generate go run github.com/anton7r/bomlok/cmd/bomlok -include=. -exclude=template_tests -exclude=testdata

// Version string
const Version = "6.0.0"

// Must panics, if a Template couldn't successfully parsed. This is how you
// would use it:
//
//	var baseTemplate = pongo2.Must(pongo2.FromFile("templates/base.html"))
func Must(tpl *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return tpl
}
