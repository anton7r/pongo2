package pongo2

import (
	"testing"
)

func TestTrimWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Basic whitespace stripping",
			input: `<div>
		<span>  foo  </span>
	</div>`,
			expected: `<div><span> foo </span></div>`,
		},
		{
			name:     "Attributes preservation (double quotes)",
			input:    `<div class="  foo  bar  "></div>`,
			expected: `<div class="  foo  bar  "></div>`,
		},
		{
			name: "Mixed attributes and tags",
			input: `<div  class="  foo  "  >
		<span> bar </span>
	</div>`,
			expected: `<div class="  foo  "><span> bar </span></div>`,
		},
		{
			name:     "With template tags inside attributes",
			input:    `<div class="{% if true %}  foo  {% endif %}">`,
			expected: `<div class="  foo  ">`,
		},
		{
			name:     "Single quotes",
			input:    `<div class='  foo  '>`,
			expected: `<div class='  foo  '>`,
		},
		{
			name:     "Variable spacing",
			input:    `foo {{ "bar" }} baz`,
			expected: `foo bar baz`,
		},
		{
			name:     "Nested quotes (HTML style)",
			input:    `<div data-val=" ' foo ' ">`,
			expected: `<div data-val=" ' foo ' ">`,
		},
		{
			name:     "Strip around equals",
			input:    `<div class = "foo">`,
			expected: `<div class="foo">`,
		},
		{
			name:     "Strip between tags",
			input:    `<div>   </div>`,
			expected: `<div></div>`,
		},
		{
			name:     "Strip leading",
			input:    `   <div>`,
			expected: `<div>`,
		},
		{
			name:     "Strip trailing",
			input:    `</div>   `,
			expected: `</div>`,
		},
	}

	set := NewSet("test_trim", MustNewLocalFileSystemLoader(""))
	set.Options.TrimWhitespace = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpl, err := set.FromString(tt.input)
			if err != nil {
				t.Fatalf("Error parsing template: %v", err)
			}
			out, err := tpl.Execute(Context{})
			if err != nil {
				t.Fatalf("Error executing template: %v", err)
			}
			if out != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, out)
			}
		})
	}
}
