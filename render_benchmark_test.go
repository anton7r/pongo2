package pongo2_test

import (
	"io"
	"testing"

	"github.com/anton7r/pongo2/v6"
)

// BenchmarkRenderSimple benchmarks rendering a simple template with variable substitution
func BenchmarkRenderSimple(b *testing.B) {
	tpl, err := pongo2.FromString("Hello {{ name }}! You have {{ count }} messages.")
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"name":  "World",
		"count": 42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderSimpleParallel benchmarks parallel rendering of a simple template
func BenchmarkRenderSimpleParallel(b *testing.B) {
	tpl, err := pongo2.FromString("Hello {{ name }}! You have {{ count }} messages.")
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"name":  "World",
		"count": 42,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkRenderWithFilters benchmarks rendering with multiple filters
func BenchmarkRenderWithFilters(b *testing.B) {
	tpl, err := pongo2.FromString(`{{ text|upper|truncatewords:5 }} - {{ number|add:10|floatformat:2 }}`)
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"text":   "this is a long sentence that will be truncated",
		"number": 3.14159,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderForLoop benchmarks rendering with a for loop
func BenchmarkRenderForLoop(b *testing.B) {
	tpl, err := pongo2.FromString(`
		{% for item in items %}
			{{ forloop.Counter }}. {{ item.name }} - {{ item.value }}
		{% endfor %}
	`)
	if err != nil {
		b.Fatal(err)
	}

	items := make([]map[string]any, 100)
	for i := 0; i < 100; i++ {
		items[i] = map[string]any{
			"name":  "Item",
			"value": i,
		}
	}

	ctx := pongo2.Context{
		"items": items,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderIfConditions benchmarks rendering with conditional logic
func BenchmarkRenderIfConditions(b *testing.B) {
	tpl, err := pongo2.FromString(`
		{% if user.is_admin %}
			Admin: {{ user.name }}
		{% elif user.is_moderator %}
			Moderator: {{ user.name }}
		{% else %}
			User: {{ user.name }}
		{% endif %}
	`)
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"user": map[string]any{
			"name":         "John",
			"is_admin":     false,
			"is_moderator": true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderComplex benchmarks rendering the complex.tpl template
func BenchmarkRenderComplex(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/complex.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(tplContext, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderComplexParallel benchmarks parallel rendering of complex.tpl
func BenchmarkRenderComplexParallel(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/complex.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := tpl.ExecuteWriterUnbuffered(tplContext, io.Discard)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkRenderInheritance benchmarks rendering with template inheritance
func BenchmarkRenderInheritance(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/extends.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(tplContext, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderInheritanceParallel benchmarks parallel rendering with inheritance
func BenchmarkRenderInheritanceParallel(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/extends.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := tpl.ExecuteWriterUnbuffered(tplContext, io.Discard)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkRenderMacro benchmarks rendering with macros
func BenchmarkRenderMacro(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/macro.tpl")
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"items": []string{"apple", "banana", "cherry"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderInclude benchmarks rendering with includes
func BenchmarkRenderInclude(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/includes.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(tplContext, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderExecuteBytes benchmarks rendering using ExecuteBytes
func BenchmarkRenderExecuteBytes(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/complex.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tpl.ExecuteBytes(tplContext)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderExecuteBytesParallel benchmarks parallel ExecuteBytes
func BenchmarkRenderExecuteBytesParallel(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/complex.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := tpl.ExecuteBytes(tplContext)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkRenderBlocks benchmarks rendering specific blocks
func BenchmarkRenderBlocks(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/block_render/block.tpl")
	if err != nil {
		b.Fatal(err)
	}

	blockNames := []string{"content", "more_content"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tpl.ExecuteBlocks(tplContext, blockNames)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderWithAutoescape benchmarks rendering with autoescaping
func BenchmarkRenderWithAutoescape(b *testing.B) {
	tpl, err := pongo2.FromFile("template_tests/autoescape.tpl")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(tplContext, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderNestedLoops benchmarks rendering with nested loops
func BenchmarkRenderNestedLoops(b *testing.B) {
	tpl, err := pongo2.FromString(`
		{% for i in outer %}
			{% for j in inner %}
				{{ i }}-{{ j }}
			{% endfor %}
		{% endfor %}
	`)
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"outer": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"inner": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderComplexFilters benchmarks rendering with chained filters
func BenchmarkRenderComplexFilters(b *testing.B) {
	tpl, err := pongo2.FromString(`{{ text|lower|truncatewords:10|title }}`)
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"text": "THIS IS A VERY LONG TEXT THAT WILL BE PROCESSED THROUGH MULTIPLE FILTERS TO TEST PERFORMANCE",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderWithFunction benchmarks rendering with function calls
func BenchmarkRenderWithFunction(b *testing.B) {
	tpl, err := pongo2.FromString(`{{ add(x, y) }} - {{ multiply(x, y) }}`)
	if err != nil {
		b.Fatal(err)
	}

	ctx := pongo2.Context{
		"x": 42,
		"y": 10,
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}
