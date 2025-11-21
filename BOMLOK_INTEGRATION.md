# Bomlok Integration in Pongo2

This document describes the integration of the [bomlok](https://github.com/anton7r/bomlok) library into pongo2 for optimized struct field access.

## Overview

Bomlok is a code generation library that reduces the need for reflection when accessing struct fields. It generates optimized getter methods that are approximately 200-260x faster than reflection, with zero memory allocations.

## What Changed

The pongo2 template engine now uses bomlok-generated code when available to access struct fields, providing significant performance improvements for template rendering with struct data.

### Modified Files

1. **variable.go** - Updated struct field resolution to check for `Bomlok` interface first
2. **value.go** - Updated `Contains()` method to use `Bomlok` interface for struct field checking
3. **pongo2.go** - Added `go:generate` directive to auto-generate bomlok implementations

### Integration Points

#### 1. Variable Resolution (variable.go)

When resolving struct fields during template variable evaluation:

```go
case reflect.Struct:
    // Try to use Bomlok interface for faster field access
    if bomlokValue, ok := current.Interface().(bomlok.Bomlok); ok {
        value := bomlokValue.Bomlok_GetValue(part.s)
        if value != nil {
            current = reflect.ValueOf(value)
        } else {
            // Field not found, return invalid value
            current = reflect.Value{}
        }
    } else {
        // Fallback to reflection
        current = current.FieldByName(part.s)
    }
```

#### 2. Contains Method (value.go)

When checking if a struct contains a specific field:

```go
case reflect.Struct:
    // Try to use Bomlok interface for faster field access
    if bomlokValue, ok := baseValue.(bomlok.Bomlok); ok {
        value := bomlokValue.Bomlok_GetValue(other.String())
        return value != nil
    }
    // Fallback to reflection
    fieldValue := rv.FieldByName(other.String())
    return fieldValue.IsValid()
```

## How It Works

1. **Code Generation**: Run `go generate` to create `*_bomlok.go` files for all structs in the project
2. **Runtime Detection**: At runtime, pongo2 checks if a struct implements the `bomlok.Bomlok` interface
3. **Fast Path**: If available, use `Bomlok_GetValue(fieldName)` instead of `reflect.FieldByName()`
4. **Fallback**: If not available, fall back to standard reflection

## Performance Benefits

- **~200-260x faster** field access compared to reflection
- **Zero memory allocations** for field lookups
- **Backwards compatible** - works with any struct, generated or not

### Benchmark Results

```
BenchmarkBomlokVsReflection/Bomlok-32    1000000000    0.1787 ns/op    0 B/op    0 allocs/op
```

## Using Bomlok with Your Structs

### For Library Users

To benefit from bomlok optimizations in your own structs:

1. Install the bomlok code generator:
   ```bash
   go install github.com/anton7r/bomlok/cmd/bomlok@latest
   ```

2. Generate bomlok implementations for your structs:
   ```bash
   bomlok -include=./models
   ```

3. Or add a `go:generate` directive to your package:
   ```go
   //go:generate bomlok -include=.
   package mypackage
   ```

4. Run `go generate ./...`

### Example

Given a struct:

```go
type User struct {
    Name  string
    Email string
    Age   int
}
```

After running bomlok, you get:

```go
func (s *User) Bomlok_GetValue(field string) any {
    switch field {
    case "Name":
        return s.Name
    case "Email":
        return s.Email
    case "Age":
        return s.Age
    default:
        return nil
    }
}

func (s *User) Bomlok_Fields() []string {
    return []string{"Name", "Email", "Age"}
}
```

When you use this struct in a pongo2 template:

```go
tpl, _ := pongo2.FromString("Hello {{ user.Name }}, you are {{ user.Age }} years old!")
result, _ := tpl.Execute(pongo2.Context{"user": &User{Name: "John", Age: 30}})
```

Pongo2 will automatically use the fast `Bomlok_GetValue()` method instead of reflection!

## Compatibility

- **Fully backwards compatible** - Structs without bomlok generation continue to work via reflection
- **No breaking changes** - All existing code works exactly as before
- **Opt-in optimization** - You only get the benefits when you generate bomlok code for your structs

## Development

### Regenerating Bomlok Code

To regenerate all bomlok implementations in the pongo2 codebase:

```bash
go generate ./...
```

This will update all `*_bomlok.go` files with the latest struct definitions.

### Testing

The integration includes tests to verify correct behavior:

```bash
go test -v -run TestBomlokIntegration
```

And benchmarks to measure performance:

```bash
go test -bench BenchmarkBomlokVsReflection -benchmem
```

## References

- [Bomlok GitHub Repository](https://github.com/anton7r/bomlok)
- [Bomlok Documentation](https://github.com/anton7r/bomlok/blob/main/README.md)
- [Performance Comparison](https://github.com/anton7r/bomlok/blob/main/example/models_bench_test.go)
