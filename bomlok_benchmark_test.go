package pongo2

import (
	"bytes"
	"reflect"
	"testing"
)

// BenchUser, BenchProduct, etc. are defined in benchmark_structs.go
// and have bomlok methods generated automatically.

// Structs without bomlok - for comparison (defined in test file so no generation)
type RegularUser struct {
	Name    string
	Email   string
	Age     int
	Active  bool
	Balance float64
}

type RegularProduct struct {
	ID          int
	Name        string
	Price       float64
	Description string
	InStock     bool
	Category    string
}

// BenchmarkFieldAccess_Bomlok benchmarks direct bomlok field access
func BenchmarkFieldAccess_Bomlok(b *testing.B) {
	user := &BenchUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	b.Run("SingleField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = user.Bomlok_GetValue("Name")
		}
	})

	b.Run("MultipleFields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = user.Bomlok_GetValue("Name")
			_ = user.Bomlok_GetValue("Email")
			_ = user.Bomlok_GetValue("Age")
			_ = user.Bomlok_GetValue("Active")
			_ = user.Bomlok_GetValue("Balance")
		}
	})

	b.Run("AllFields", func(b *testing.B) {
		fields := user.Bomlok_Fields()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, field := range fields {
				_ = user.Bomlok_GetValue(field)
			}
		}
	})
}

// BenchmarkFieldAccess_Reflection benchmarks reflection-based field access
func BenchmarkFieldAccess_Reflection(b *testing.B) {
	user := &RegularUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	b.Run("SingleField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := reflect.ValueOf(user).Elem()
			_ = rv.FieldByName("Name").Interface()
		}
	})

	b.Run("MultipleFields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := reflect.ValueOf(user).Elem()
			_ = rv.FieldByName("Name").Interface()
			_ = rv.FieldByName("Email").Interface()
			_ = rv.FieldByName("Age").Interface()
			_ = rv.FieldByName("Active").Interface()
			_ = rv.FieldByName("Balance").Interface()
		}
	})

	b.Run("AllFields", func(b *testing.B) {
		rv := reflect.ValueOf(user).Elem()
		rt := rv.Type()
		numFields := rt.NumField()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < numFields; j++ {
				fieldName := rt.Field(j).Name
				_ = rv.FieldByName(fieldName).Interface()
			}
		}
	})
}

func repeat[T any](slice []T, count int) []T {
	result := make([]T, 0, len(slice)*count)
	for i := 0; i < count; i++ {
		result = append(result, slice...)
	}
	return result
}

// BenchmarkTemplateRendering benchmarks template rendering with bomlok vs reflection
func BenchmarkTemplateRendering_Bomlok(b *testing.B) {
	user := &BenchUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	product := &BenchProduct{
		ID:          123,
		Name:        "Laptop",
		Price:       999.99,
		Description: "High performance laptop",
		InStock:     true,
		Category:    "Electronics",
	}

	b.Run("SimpleTemplate", func(b *testing.B) {
		tpl, _ := FromString("Hello {{ user.Name }}, your email is {{ user.Email }}")
		ctx := Context{"user": user}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})

	b.Run("ComplexTemplate", func(b *testing.B) {
		tpl, _ := FromString(`
			User: {{ user.Name }} ({{ user.Email }})
			Age: {{ user.Age }}
			Status: {% if user.Active %}Active{% else %}Inactive{% endif %}
			Balance: ${{ user.Balance }}
			
			Product: {{ product.Name }}
			Price: ${{ product.Price }}
			{% if product.InStock %}In Stock{% else %}Out of Stock{% endif %}
		`)
		ctx := Context{"user": user, "product": product}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})

	b.Run("LoopTemplate", func(b *testing.B) {
		users := repeat([]*BenchUser{user}, 500)
		tpl, _ := FromString(`
			{% for u in users %}
				{{ u.Name }} - {{ u.Email }} - {{ u.Age }} - {{ u.Balance }}
			{% endfor %}
		`)
		ctx := Context{"users": users}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})
}

func BenchmarkTemplateRendering_Reflection(b *testing.B) {
	user := &RegularUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	product := &RegularProduct{
		ID:          123,
		Name:        "Laptop",
		Price:       999.99,
		Description: "High performance laptop",
		InStock:     true,
		Category:    "Electronics",
	}

	b.Run("SimpleTemplate", func(b *testing.B) {
		tpl, _ := FromString("Hello {{ user.Name }}, your email is {{ user.Email }}")
		ctx := Context{"user": user}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})

	b.Run("ComplexTemplate", func(b *testing.B) {
		tpl, _ := FromString(`
			User: {{ user.Name }} ({{ user.Email }})
			Age: {{ user.Age }}
			Status: {% if user.Active %}Active{% else %}Inactive{% endif %}
			Balance: ${{ user.Balance }}
			
			Product: {{ product.Name }}
			Price: ${{ product.Price }}
			{% if product.InStock %}In Stock{% else %}Out of Stock{% endif %}
		`)
		ctx := Context{"user": user, "product": product}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})

	b.Run("LoopTemplate", func(b *testing.B) {
		users := repeat([]*RegularUser{user}, 500)
		tpl, _ := FromString(`
			{% for u in users %}
				{{ u.Name }} - {{ u.Email }} - {{ u.Age }} - {{ u.Balance }}
			{% endfor %}
		`)
		ctx := Context{"users": users}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})
}

// BenchmarkVariableResolver benchmarks the variable resolver with different types
func BenchmarkVariableResolver_Bomlok(b *testing.B) {
	user := &BenchUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	tpl, _ := FromString("{{ user.Name }}")
	ctx := &ExecutionContext{
		Public:  Context{"user": user},
		Private: make(Context),
	}

	var buf bytes.Buffer
	writer := &templateWriter{w: &buf}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = tpl.root.Execute(ctx, writer)
	}
}

func BenchmarkVariableResolver_Reflection(b *testing.B) {
	user := &RegularUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	tpl, _ := FromString("{{ user.Name }}")
	ctx := &ExecutionContext{
		Public:  Context{"user": user},
		Private: make(Context),
	}

	var buf bytes.Buffer
	writer := &templateWriter{w: &buf}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = tpl.root.Execute(ctx, writer)
	}
}

// BenchmarkContainsMethod benchmarks the Contains method
func BenchmarkContainsMethod_Bomlok(b *testing.B) {
	user := &BenchUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	userValue := AsValue(user)
	fieldName := AsValue("Name")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userValue.Contains(fieldName)
	}
}

func BenchmarkContainsMethod_Reflection(b *testing.B) {
	user := &RegularUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	userValue := AsValue(user)
	fieldName := AsValue("Name")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userValue.Contains(fieldName)
	}
}

// BenchmarkNestedFieldAccess benchmarks nested struct field access
func BenchmarkNestedFieldAccess(b *testing.B) {
	user := &BenchUserWithProfile{
		Name: "John Doe",
		Profile: BenchProfile{
			Bio: "Software Engineer",
			Address: BenchAddress{
				Street:  "123 Main St",
				City:    "New York",
				ZipCode: "10001",
			},
		},
	}

	b.Run("Bomlok_NestedAccess", func(b *testing.B) {
		tpl, _ := FromString("{{ user.Profile.Address.City }}")
		ctx := Context{"user": user}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})

	type RegularAddress struct {
		Street  string
		City    string
		ZipCode string
	}

	type RegularProfile struct {
		Bio     string
		Address RegularAddress
	}

	type RegularUserWithProfile struct {
		Name    string
		Profile RegularProfile
	}

	regularUser := &RegularUserWithProfile{
		Name: "John Doe",
		Profile: RegularProfile{
			Bio: "Software Engineer",
			Address: RegularAddress{
				Street:  "123 Main St",
				City:    "New York",
				ZipCode: "10001",
			},
		},
	}

	b.Run("Reflection_NestedAccess", func(b *testing.B) {
		tpl, _ := FromString("{{ user.Profile.Address.City }}")
		ctx := Context{"user": regularUser}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = tpl.Execute(ctx)
		}
	})
}

// BenchmarkDifferentFieldTypes benchmarks access to different field types
func BenchmarkDifferentFieldTypes_Bomlok(b *testing.B) {
	user := &BenchUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	b.Run("StringField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = user.Bomlok_GetValue("Name")
		}
	})

	b.Run("IntField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = user.Bomlok_GetValue("Age")
		}
	})

	b.Run("BoolField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = user.Bomlok_GetValue("Active")
		}
	})

	b.Run("Float64Field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = user.Bomlok_GetValue("Balance")
		}
	})
}

func BenchmarkDifferentFieldTypes_Reflection(b *testing.B) {
	user := &RegularUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	b.Run("StringField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := reflect.ValueOf(user).Elem()
			_ = rv.FieldByName("Name").Interface()
		}
	})

	b.Run("IntField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := reflect.ValueOf(user).Elem()
			_ = rv.FieldByName("Age").Interface()
		}
	})

	b.Run("BoolField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := reflect.ValueOf(user).Elem()
			_ = rv.FieldByName("Active").Interface()
		}
	})

	b.Run("Float64Field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := reflect.ValueOf(user).Elem()
			_ = rv.FieldByName("Balance").Interface()
		}
	})
}

// BenchmarkCachedReflectionValue benchmarks with cached reflection value (best case for reflection)
func BenchmarkCachedReflectionValue(b *testing.B) {
	user := &RegularUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	rv := reflect.ValueOf(user).Elem()

	b.Run("CachedReflection_SingleField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = rv.FieldByName("Name").Interface()
		}
	})

	bomlokUser := &BenchUser{
		Name:    "John Doe",
		Email:   "john@example.com",
		Age:     30,
		Active:  true,
		Balance: 1500.50,
	}

	b.Run("Bomlok_SingleField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = bomlokUser.Bomlok_GetValue("Name")
		}
	})
}
