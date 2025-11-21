package pongo2

// Structs used for benchmarking bomlok performance
// These will have bomlok methods generated via go:generate

type BenchUser struct {
	Name    string
	Email   string
	Age     int
	Active  bool
	Balance float64
}

type BenchProduct struct {
	ID          int
	Name        string
	Price       float64
	Description string
	InStock     bool
	Category    string
}

type BenchOrder struct {
	OrderID    string
	CustomerID int
	Total      float64
	Items      []string
	Completed  bool
}

type BenchAddress struct {
	Street  string
	City    string
	ZipCode string
}

type BenchProfile struct {
	Bio     string
	Address BenchAddress
}

type BenchUserWithProfile struct {
	Name    string
	Profile BenchProfile
}
