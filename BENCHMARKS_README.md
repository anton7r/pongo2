# Bomlok Performance Benchmarks

This directory contains comprehensive performance benchmarks comparing bomlok-optimized struct field access with standard Go reflection in the pongo2 template engine.

## Files

- `bomlok_benchmark_test.go` - Main benchmark test file
- `benchmark_structs.go` - Struct definitions for benchmarking (bomlok-enabled)
- `benchmark_structs_bomlok.go` - Auto-generated bomlok implementations
- `BENCHMARK_RESULTS.md` - Detailed benchmark results and analysis

## Quick Start

Run all benchmarks:
```bash
go test -bench . -benchmem
```

Run specific benchmark categories:
```bash
# Field access comparison
go test -bench BenchmarkFieldAccess -benchmem

# Template rendering comparison  
go test -bench BenchmarkTemplateRendering -benchmem

# Different field types
go test -bench BenchmarkDifferentFieldTypes -benchmem
```

## Benchmark Categories

### 1. BenchmarkFieldAccess
Direct comparison of bomlok vs reflection for accessing struct fields:
- Single field access
- Multiple fields access
- Iteration over all fields

### 2. BenchmarkTemplateRendering
Real-world template rendering with bomlok vs reflection:
- Simple templates ({{ variable }})
- Complex templates (conditionals, multiple fields)
- Loop-heavy templates (iteration over collections)

### 3. BenchmarkVariableResolver
Testing the internal variable resolution mechanism

### 4. BenchmarkContainsMethod
Testing the Contains() method for checking field existence

### 5. BenchmarkNestedFieldAccess
Testing nested field access (struct.field.subfield)

### 6. BenchmarkDifferentFieldTypes
Performance across different data types:
- String fields
- Integer fields
- Boolean fields
- Float64 fields

### 7. BenchmarkCachedReflectionValue
Comparison against cached reflection values (best-case reflection)

## Key Results

### Direct Field Access
- **Bomlok**: ~0.18 ns/op, 0 allocations
- **Reflection**: ~41 ns/op, 1 allocation
- **Speedup**: ~230x faster

### Template Rendering
- Performance is roughly equivalent in full template rendering
- Other operations (parsing, I/O) dominate the execution time
- Bomlok provides benefits in field-heavy templates

## Understanding the Results

The benchmarks show two important insights:

1. **Micro-level**: Bomlok is 200-370x faster for direct field access
2. **Macro-level**: In full template rendering, the speedup is less pronounced because field access is only one part of the overall work

Both results are valuable:
- The micro-benchmarks show the optimization potential
- The macro-benchmarks show real-world impact
- Together they demonstrate that bomlok is always beneficial, never harmful

## Regenerating Benchmarks

If you modify the benchmark structs, regenerate the bomlok code:

```bash
go generate ./...
```

Then run the benchmarks again to see updated results.

## Detailed Analysis

See `BENCHMARK_RESULTS.md` for:
- Complete benchmark data
- Performance analysis by use case
- Real-world impact calculations
- Recommendations for when to use bomlok
