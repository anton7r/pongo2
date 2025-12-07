# Loop Performance Optimization

## Problem Identified

Initial benchmarks showed that templates with loops were **slower** with bomlok than with reflection:

```
BenchmarkTemplateRendering_Bomlok/LoopTemplate        5082 ns/op    113 allocs/op
BenchmarkTemplateRendering_Reflection/LoopTemplate    4790 ns/op    113 allocs/op
```

Bomlok was 6% slower in loop scenarios.

## Root Cause Analysis

The issue was in the fast path optimization in `variable.go`. When items are iterated in a loop:

1. The loop iterator stores items as `*Value` objects in the private context
2. When accessing fields (e.g., `{{ u.Name }}`), the fast path checked for `bomlok.Bomlok` interface directly on the context value
3. But the actual struct was **wrapped inside** the `*Value` object
4. So the fast path failed and fell back to the slow reflection path

## Solution

Modified the fast path to **unwrap** `*Value` objects before checking for bomlok interface:

```go
// If val is a *Value (common in loops), unwrap it
if valuePtr, ok := val.(*Value); ok {
    val = valuePtr.val
    isSafe = valuePtr.safe
}

if bomlokValue, ok := val.(bomlok.Bomlok); ok {
    fieldValue := bomlokValue.Bomlok_GetValue(vr.parts[1].s)
    if fieldValue != nil {
        return &Value{val: fieldValue, safe: isSafe}, nil
    }
}
```

## Results

### Small Loop (5 items)
**Before fix:**
```
BenchmarkTemplateRendering_Bomlok/LoopTemplate        5082 ns/op    4718 B/op    113 allocs/op
BenchmarkTemplateRendering_Reflection/LoopTemplate    4824 ns/op    4718 B/op    113 allocs/op
```

**After fix:**
```
BenchmarkTemplateRendering_Bomlok/LoopTemplate        2926 ns/op    3394 B/op     88 allocs/op
BenchmarkTemplateRendering_Reflection/LoopTemplate    4824 ns/op    4718 B/op    113 allocs/op
```

**Improvement:**
- **1.74x faster** (5082 → 2926 ns/op)
- **28% less memory** (4718 → 3394 B/op)
- **22% fewer allocations** (113 → 88 allocs/op)
- Now **1.65x faster than reflection**!

### Large Loop (500 items)
**After fix:**
```
BenchmarkTemplateRendering_Bomlok/LoopTemplate        200,722 ns/op    252,578 B/op    7,031 allocs/op
BenchmarkTemplateRendering_Reflection/LoopTemplate    387,112 ns/op    385,001 B/op    9,531 allocs/op
```

**Improvement:**
- **1.93x faster** than reflection
- **34% less memory**
- **26% fewer allocations**

## Impact on Other Benchmarks

All other benchmarks remained excellent or improved:

| Benchmark | Bomlok | Reflection | Speedup |
|-----------|--------|------------|---------|
| SimpleTemplate | 680 ns/op (16 allocs) | 900 ns/op (18 allocs) | **1.32x faster** |
| ComplexTemplate | 1443 ns/op (31 allocs) | 2274 ns/op (42 allocs) | **1.58x faster** |
| VariableResolver | 54 ns/op (3 allocs) | 119 ns/op (4 allocs) | **2.20x faster** |
| FieldAccess | 0.18 ns/op (0 allocs) | 41.7 ns/op (1 alloc) | **232x faster** |

## Conclusion

By unwrapping `*Value` objects in the fast path, we:
1. ✅ Fixed the loop performance regression
2. ✅ Made loops **1.7-1.9x faster** with bomlok vs reflection
3. ✅ Reduced memory usage by 28-34% in loops
4. ✅ Reduced allocations by 22-26% in loops
5. ✅ Maintained excellent performance in all other scenarios

Bomlok is now consistently faster than reflection **across all template rendering scenarios**, including loops!
