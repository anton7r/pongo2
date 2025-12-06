package pongo2_test

import (
	"io"
	"testing"
	"time"

	"github.com/anton7r/pongo2/v6"
)

// Struct types for benchmarking struct field access
type BenchUser struct {
	ID        int
	Username  string
	Email     string
	FirstName string
	LastName  string
	Age       int
	IsActive  bool
	CreatedAt time.Time
}

type BenchProduct struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    string
	Tags        []string
}

type BenchOrder struct {
	OrderID     string
	User        BenchUser
	Products    []BenchProduct
	Total       float64
	Status      string
	OrderDate   time.Time
	ShippedDate *time.Time
}

// BenchmarkRenderStructFieldAccess benchmarks accessing struct fields
func BenchmarkRenderStructFieldAccess(b *testing.B) {
	tpl, err := pongo2.FromString(`
		User: {{ user.Username }} ({{ user.Email }})
		Name: {{ user.FirstName }} {{ user.LastName }}
		Age: {{ user.Age }}
		Active: {{ user.IsActive }}
	`)
	if err != nil {
		b.Fatal(err)
	}

	user := BenchUser{
		ID:        1,
		Username:  "johndoe",
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	ctx := pongo2.Context{
		"user": user,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderStructFieldAccessParallel benchmarks parallel struct field access
func BenchmarkRenderStructFieldAccessParallel(b *testing.B) {
	tpl, err := pongo2.FromString(`
		User: {{ user.Username }} ({{ user.Email }})
		Name: {{ user.FirstName }} {{ user.LastName }}
		Age: {{ user.Age }}
		Active: {{ user.IsActive }}
	`)
	if err != nil {
		b.Fatal(err)
	}

	user := BenchUser{
		ID:        1,
		Username:  "johndoe",
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	ctx := pongo2.Context{
		"user": user,
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

// BenchmarkRenderNestedStructAccess benchmarks accessing nested struct fields
func BenchmarkRenderNestedStructAccess(b *testing.B) {
	tpl, err := pongo2.FromString(`
		Order: {{ order.OrderID }}
		Customer: {{ order.User.FirstName }} {{ order.User.LastName }}
		Email: {{ order.User.Email }}
		Total: {{ order.Total }}
		Status: {{ order.Status }}
	`)
	if err != nil {
		b.Fatal(err)
	}

	order := BenchOrder{
		OrderID: "ORD-12345",
		User: BenchUser{
			ID:        1,
			Username:  "johndoe",
			Email:     "john@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Age:       30,
			IsActive:  true,
			CreatedAt: time.Now(),
		},
		Total:     299.99,
		Status:    "shipped",
		OrderDate: time.Now(),
	}

	ctx := pongo2.Context{
		"order": order,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderStructSliceIteration benchmarks iterating over a slice of structs
func BenchmarkRenderStructSliceIteration(b *testing.B) {
	tpl, err := pongo2.FromString(`
		{% for product in products %}
			{{ product.ID }}. {{ product.Name }}
			Price: ${{ product.Price }}
			Stock: {{ product.Stock }}
			Category: {{ product.Category }}
		{% endfor %}
	`)
	if err != nil {
		b.Fatal(err)
	}

	products := make([]BenchProduct, 50)
	for i := 0; i < 50; i++ {
		products[i] = BenchProduct{
			ID:          i + 1,
			Name:        "Product " + string(rune(i)),
			Description: "Description for product",
			Price:       19.99 + float64(i),
			Stock:       100 - i,
			Category:    "Category",
			Tags:        []string{"tag1", "tag2"},
		}
	}

	ctx := pongo2.Context{
		"products": products,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderComplexStructTemplate benchmarks a complex template with nested structs
func BenchmarkRenderComplexStructTemplate(b *testing.B) {
	tpl, err := pongo2.FromString(`
		<h1>Order {{ order.OrderID }}</h1>
		<div class="customer">
			<h2>Customer Information</h2>
			<p>Name: {{ order.User.FirstName }} {{ order.User.LastName }}</p>
			<p>Email: {{ order.User.Email }}</p>
			<p>Username: {{ order.User.Username }}</p>
			<p>Account Status: {{ order.User.IsActive|yesno:"Active,Inactive" }}</p>
		</div>
		
		<div class="products">
			<h2>Products</h2>
			{% for product in order.Products %}
				<div class="product">
					<h3>{{ product.Name }}</h3>
					<p>{{ product.Description }}</p>
					<p>Price: ${{ product.Price|floatformat:2 }}</p>
					<p>Stock: {{ product.Stock }} units</p>
					<p>Category: {{ product.Category }}</p>
					<p>Tags: {% for tag in product.Tags %}{{ tag }}{% if not forloop.Last %}, {% endif %}{% endfor %}</p>
				</div>
			{% endfor %}
		</div>
		
		<div class="summary">
			<h2>Order Summary</h2>
			<p>Total: ${{ order.Total|floatformat:2 }}</p>
			<p>Status: {{ order.Status|upper }}</p>
			<p>Order Date: {{ order.OrderDate }}</p>
		</div>
	`)
	if err != nil {
		b.Fatal(err)
	}

	products := []BenchProduct{
		{
			ID:          1,
			Name:        "Laptop",
			Description: "High-performance laptop",
			Price:       999.99,
			Stock:       10,
			Category:    "Electronics",
			Tags:        []string{"computers", "electronics", "best-seller"},
		},
		{
			ID:          2,
			Name:        "Mouse",
			Description: "Wireless mouse",
			Price:       29.99,
			Stock:       50,
			Category:    "Accessories",
			Tags:        []string{"accessories", "peripherals"},
		},
		{
			ID:          3,
			Name:        "Keyboard",
			Description: "Mechanical keyboard",
			Price:       79.99,
			Stock:       30,
			Category:    "Accessories",
			Tags:        []string{"accessories", "peripherals", "gaming"},
		},
	}

	order := BenchOrder{
		OrderID: "ORD-98765",
		User: BenchUser{
			ID:        42,
			Username:  "janedoe",
			Email:     "jane@example.com",
			FirstName: "Jane",
			LastName:  "Doe",
			Age:       28,
			IsActive:  true,
			CreatedAt: time.Now().Add(-365 * 24 * time.Hour),
		},
		Products:  products,
		Total:     1109.97,
		Status:    "processing",
		OrderDate: time.Now().Add(-2 * 24 * time.Hour),
	}

	ctx := pongo2.Context{
		"order": order,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderStructVsMap benchmarks struct access vs map access
func BenchmarkRenderStructVsMap(b *testing.B) {
	tplStr := `{{ obj.Field1 }} {{ obj.Field2 }} {{ obj.Field3 }} {{ obj.Field4 }} {{ obj.Field5 }}`

	b.Run("Struct", func(b *testing.B) {
		tpl, err := pongo2.FromString(tplStr)
		if err != nil {
			b.Fatal(err)
		}

		type TestStruct struct {
			Field1 string
			Field2 int
			Field3 float64
			Field4 bool
			Field5 string
		}

		obj := TestStruct{
			Field1: "value1",
			Field2: 42,
			Field3: 3.14,
			Field4: true,
			Field5: "value5",
		}

		ctx := pongo2.Context{"obj": obj}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Map", func(b *testing.B) {
		tpl, err := pongo2.FromString(tplStr)
		if err != nil {
			b.Fatal(err)
		}

		obj := map[string]any{
			"Field1": "value1",
			"Field2": 42,
			"Field3": 3.14,
			"Field4": true,
			"Field5": "value5",
		}

		ctx := pongo2.Context{"obj": obj}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = tpl.ExecuteWriterUnbuffered(ctx, io.Discard)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
