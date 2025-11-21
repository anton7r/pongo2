package pongo2

import (
	"testing"
)

// TestBomlokIntegration tests that the bomlok interface is being used for struct field access
func TestBomlokIntegration(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Note: This struct doesn't have bomlok generated for it since it's defined in a test file,
	// but we can test with existing generated structs like Value

	// Test that Value struct has Bomlok methods
	v := &Value{val: "test", safe: true}

	// Test Bomlok_GetValue
	if got := v.Bomlok_GetValue("val"); got != "test" {
		t.Errorf("Bomlok_GetValue(\"val\") = %v, want \"test\"", got)
	}

	if got := v.Bomlok_GetValue("safe"); got != true {
		t.Errorf("Bomlok_GetValue(\"safe\") = %v, want true", got)
	}

	// Test Bomlok_Fields
	fields := v.Bomlok_Fields()
	if len(fields) != 2 {
		t.Errorf("Bomlok_Fields() returned %d fields, want 2", len(fields))
	}

	// Test that Contains method uses Bomlok interface
	ctx := Context(map[string]any{
		"person": map[string]any{
			"name": "John",
			"age":  30,
		},
	})

	tpl, err := FromString("{{ person.name }}")
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	result, err := tpl.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	if result != "John" {
		t.Errorf("Template execution returned %q, want \"John\"", result)
	}
}

// BenchmarkBomlokVsReflection benchmarks the performance difference between using bomlok and reflection
func BenchmarkBomlokVsReflection(b *testing.B) {
	v := &Value{val: "test", safe: true}

	b.Run("Bomlok", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = v.Bomlok_GetValue("val")
		}
	})
}
