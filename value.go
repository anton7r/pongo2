package pongo2

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Value struct {
	val  any
	safe bool // used to indicate whether a Value needs explicit escaping in the template
}

// AsValue converts any given value to a pongo2.Value
// Usually being used within own functions passed to a template
// through a Context or within filter functions.
//
// Example:
//
//	AsValue("my string")
func AsValue(i any) *Value {
	return &Value{
		val: i,
	}
}

// AsSafeValue works like AsValue, but does not apply the 'escape' filter.
func AsSafeValue(i any) *Value {
	return &Value{
		val:  i,
		safe: true,
	}
}

func (v *Value) getResolvedValue() any {
	// Dereference pointer if needed
	if v.val == nil {
		return nil
	}

	// Handle common pointer types directly without reflection
	switch val := v.val.(type) {
	case *string:
		return *val
	case *int:
		return *val
	case *int8:
		return *val
	case *int16:
		return *val
	case *int32:
		return *val
	case *int64:
		return *val
	case *uint:
		return *val
	case *uint8:
		return *val
	case *uint16:
		return *val
	case *uint32:
		return *val
	case *uint64:
		return *val
	case *float32:
		return *val
	case *float64:
		return *val
	case *bool:
		return *val
	case *time.Time:
		return *val
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, time.Time:
		// Common non-pointer types - return as-is
		return val
	}

	// Use reflection only for uncommon pointer types
	rv := reflect.ValueOf(v.val)
	if rv.IsValid() && rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		return rv.Elem().Interface()
	}
	return v.val
}

// IsString checks whether the underlying value is a string
func (v *Value) IsString() bool {
	val := v.getResolvedValue()
	_, ok := val.(string)
	return ok
}

// IsBool checks whether the underlying value is a bool
func (v *Value) IsBool() bool {
	val := v.getResolvedValue()
	_, ok := val.(bool)
	return ok
}

// IsFloat checks whether the underlying value is a float
func (v *Value) IsFloat() bool {
	val := v.getResolvedValue()
	switch val.(type) {
	case float32, float64:
		return true
	}
	return false
}

// IsInteger checks whether the underlying value is an integer
func (v *Value) IsInteger() bool {
	val := v.getResolvedValue()
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

// IsNumber checks whether the underlying value is either an integer
// or a float.
func (v *Value) IsNumber() bool {
	return v.IsInteger() || v.IsFloat()
}

// IsTime checks whether the underlying value is a time.Time.
func (v *Value) IsTime() bool {
	_, ok := v.val.(time.Time)
	return ok
}

// IsNil checks whether the underlying value is NIL
func (v *Value) IsNil() bool {
	return v.val == nil
}

// String returns a string for the underlying value. If this value is not
// of type string, pongo2 tries to convert it. Currently the following
// types for underlying values are supported:
//
//  1. string
//  2. int/uint (any size)
//  3. float (any precision)
//  4. bool
//  5. time.Time
//  6. String() will be called on the underlying value if provided
//
// NIL values will lead to an empty string. Unsupported types are leading
// to their respective type name.
func (v *Value) String() string {
	if v.IsNil() {
		return ""
	}

	val := v.getResolvedValue()

	if t, ok := val.(fmt.Stringer); ok {
		return t.String()
	}

	switch val := val.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return fmt.Sprintf("%f", val)
	case float64:
		return fmt.Sprintf("%f", val)
	case bool:
		if val {
			return "True"
		}
		return "False"
	}

	logf("Value.String() not implemented for type: %T\n", val)
	return fmt.Sprintf("%v", val)
}

// Integer returns the underlying value as an integer (converts the underlying
// value, if necessary). If it's not possible to convert the underlying value,
// it will return 0.
func (v *Value) Integer() int {
	val := v.getResolvedValue()

	switch val := val.(type) {
	case int:
		return val
	case int8:
		return int(val)
	case int16:
		return int(val)
	case int32:
		return int(val)
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint8:
		return int(val)
	case uint16:
		return int(val)
	case uint32:
		return int(val)
	case uint64:
		return int(val)
	case float32:
		return int(val)
	case float64:
		return int(val)
	case string:
		// Try to convert from string to int (base 10)
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
		return int(f)
	default:
		logf("Value.Integer() not available for type: %T\n", val)
		return 0
	}
}

// Float returns the underlying value as a float (converts the underlying
// value, if necessary). If it's not possible to convert the underlying value,
// it will return 0.0.
func (v *Value) Float() float64 {
	val := v.getResolvedValue()

	switch val := val.(type) {
	case float32:
		return float64(val)
	case float64:
		return val
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case string:
		// Try to convert from string to float64 (base 10)
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0.0
		}
		return f
	default:
		logf("Value.Float() not available for type: %T\n", val)
		return 0.0
	}
}

// Bool returns the underlying value as bool. If the value is not bool, false
// will always be returned. If you're looking for true/false-evaluation of the
// underlying value, have a look on the IsTrue()-function.
func (v *Value) Bool() bool {
	val := v.getResolvedValue()
	if b, ok := val.(bool); ok {
		return b
	}
	logf("Value.Bool() not available for type: %T\n", val)
	return false
}

// Time returns the underlying value as time.Time.
// If the underlying value is not a time.Time, it returns the zero value of time.Time.
func (v *Value) Time() time.Time {
	tm, ok := v.val.(time.Time)
	if ok {
		return tm
	}
	return time.Time{}
}

// IsTrue tries to evaluate the underlying value the Pythonic-way:
//
// Returns TRUE in one the following cases:
//
//   - int != 0
//   - uint != 0
//   - float != 0.0
//   - len(array/chan/map/slice/string) > 0
//   - bool == true
//   - underlying value is a struct
//
// Otherwise returns always FALSE.
func (v *Value) IsTrue() bool {
	val := v.getResolvedValue()
	if val == nil {
		return false
	}

	switch val := val.(type) {
	case int:
		return val != 0
	case int8:
		return val != 0
	case int16:
		return val != 0
	case int32:
		return val != 0
	case int64:
		return val != 0
	case uint:
		return val != 0
	case uint8:
		return val != 0
	case uint16:
		return val != 0
	case uint32:
		return val != 0
	case uint64:
		return val != 0
	case float32:
		return val != 0
	case float64:
		return val != 0
	case bool:
		return val
	case string:
		return len(val) > 0
	default:
		// For complex types, use reflection
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
			return rv.Len() > 0
		case reflect.Struct:
			return true // struct instance is always true
		default:
			logf("Value.IsTrue() not available for type: %T\n", val)
			return false
		}
	}
}

// Negate tries to negate the underlying value. It's mainly used for
// the NOT-operator and in conjunction with a call to
// return_value.IsTrue() afterwards.
//
// Example:
//
//	AsValue(1).Negate().IsTrue() == false
func (v *Value) Negate() *Value {
	val := v.getResolvedValue()
	if val == nil {
		return AsValue(true)
	}

	switch val := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if v.Integer() != 0 {
			return AsValue(0)
		}
		return AsValue(1)
	case float32, float64:
		if v.Float() != 0.0 {
			return AsValue(float64(0.0))
		}
		return AsValue(float64(1.1))
	case bool:
		return AsValue(!val)
	case string:
		return AsValue(len(val) == 0)
	default:
		// For complex types, use reflection
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
			return AsValue(rv.Len() == 0)
		case reflect.Struct:
			return AsValue(false)
		default:
			logf("Value.IsTrue() not available for type: %T\n", val)
			return AsValue(true)
		}
	}
}

// Len returns the length for an array, chan, map, slice or string.
// Otherwise it will return 0.
func (v *Value) Len() int {
	val := v.getResolvedValue()
	if val == nil {
		return 0
	}

	if str, ok := val.(string); ok {
		return len([]rune(str))
	}

	// For arrays, slices, maps, and channels, use reflection
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return rv.Len()
	default:
		logf("Value.Len() not available for type: %T\n", val)
		return 0
	}
}

// Slice slices an array, slice or string. Otherwise it will
// return an empty []int.
func (v *Value) Slice(i, j int) *Value {
	val := v.getResolvedValue()
	if val == nil {
		return AsValue([]int{})
	}

	if str, ok := val.(string); ok {
		runes := []rune(str)
		return AsValue(string(runes[i:j]))
	}

	// For arrays and slices, use reflection
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		return AsValue(rv.Slice(i, j).Interface())
	default:
		logf("Value.Slice() not available for type: %T\n", val)
		return AsValue([]int{})
	}
}

// Index gets the i-th item of an array, slice or string. Otherwise
// it will return NIL.
func (v *Value) Index(i int) *Value {
	val := v.getResolvedValue()
	if val == nil {
		return AsValue(nil)
	}

	if str, ok := val.(string); ok {
		runes := []rune(str)
		if i < len(runes) {
			return AsValue(string(runes[i]))
		}
		return AsValue("")
	}

	// For arrays and slices, use reflection
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		if i >= rv.Len() {
			return AsValue(nil)
		}
		return AsValue(rv.Index(i).Interface())
	default:
		logf("Value.Slice() not available for type: %T\n", val)
		return AsValue([]int{})
	}
}

// Contains checks whether the underlying value (which must be of type struct, map,
// string, array or slice) contains of another Value (e. g. used to check
// whether a struct contains of a specific field or a map contains a specific key).
//
// Example:
//
//	AsValue("Hello, World!").Contains(AsValue("World")) == true
func (v *Value) Contains(other *Value) bool {
	baseValue := v.getResolvedValue()
	if baseValue == nil {
		return false
	}

	// Handle string case directly
	if str, ok := baseValue.(string); ok {
		return strings.Contains(str, other.String())
	}

	// For other types, use reflection
	rv := reflect.ValueOf(baseValue)
	switch rv.Kind() {
	case reflect.Struct:
		fieldValue := rv.FieldByName(other.String())
		return fieldValue.IsValid()
	case reflect.Map:
		otherVal := other.getResolvedValue()
		if otherVal == nil {
			return false
		}
		otherRV := reflect.ValueOf(otherVal)
		if !otherRV.IsValid() {
			return false
		}
		// Ensure that map key type is equal to other's type.
		if rv.Type().Key() != otherRV.Type() {
			return false
		}

		mapValue := rv.MapIndex(otherRV)
		return mapValue.IsValid()
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i)
			if other.EqualValueTo(AsValue(item.Interface())) {
				return true
			}
		}
		return false

	default:
		logf("Value.Contains() not available for type: %T\n", baseValue)
		return false
	}
}

// CanSlice checks whether the underlying value is of type array, slice or string.
// You normally would use CanSlice() before using the Slice() operation.
func (v *Value) CanSlice() bool {
	val := v.getResolvedValue()
	if val == nil {
		return false
	}

	if _, ok := val.(string); ok {
		return true
	}

	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		return true
	}
	return false
}

// Iterate iterates over a map, array, slice or a string. It calls the
// function's first argument for every value with the following arguments:
//
//	idx      current 0-index
//	count    total amount of items
//	key      *Value for the key or item
//	value    *Value (only for maps, the respective value for a specific key)
//
// If the underlying value has no items or is not one of the types above,
// the empty function (function's second argument) will be called.
func (v *Value) Iterate(fn func(idx, count int, key, value *Value) bool, empty func()) {
	v.IterateOrder(fn, empty, false, false)
}

// IterateOrder behaves like Value.Iterate, but can iterate through an array/slice/string in reverse. Does
// not affect the iteration through a map because maps don't have any particular order.
// However, you can force an order using the `sorted` keyword (and even use `reversed sorted`).
func (v *Value) IterateOrder(fn func(idx, count int, key, value *Value) bool, empty func(), reverse bool, sorted bool) {
	val := v.getResolvedValue()
	if val == nil {
		empty()
		return
	}

	// Handle string case directly
	if str, ok := val.(string); ok {
		rs := []rune(str)
		charCount := len(rs)

		if charCount > 0 {
			if sorted {
				sort.SliceStable(rs, func(i, j int) bool {
					return rs[i] < rs[j]
				})
			}

			if reverse {
				for i, j := 0, charCount-1; i < j; i, j = i+1, j-1 {
					rs[i], rs[j] = rs[j], rs[i]
				}
			}

			for i := 0; i < charCount; i++ {
				if !fn(i, charCount, &Value{val: string(rs[i])}, nil) {
					return
				}
			}
		} else {
			empty()
		}
		return // done
	}

	// For other types, use reflection
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Map:
		keys := sortedKeys(rv.MapKeys())
		if sorted {
			if reverse {
				sort.Sort(sort.Reverse(keys))
			} else {
				sort.Sort(keys)
			}
		}
		keyLen := len(keys)
		for idx, key := range keys {
			value := rv.MapIndex(key)
			var keyVal, valueVal any
			if key.IsValid() {
				keyVal = key.Interface()
			}
			if value.IsValid() {
				valueVal = value.Interface()
			}
			if !fn(idx, keyLen, &Value{val: keyVal}, &Value{val: valueVal}) {
				return
			}
		}
		if keyLen == 0 {
			empty()
		}
		return // done
	case reflect.Array, reflect.Slice:
		var items valuesList

		itemCount := rv.Len()
		for i := 0; i < itemCount; i++ {
			itemRV := rv.Index(i)
			var itemVal any
			if itemRV.IsValid() {
				itemVal = itemRV.Interface()
			}
			items = append(items, &Value{val: itemVal})
		}

		if sorted {
			if reverse {
				sort.Sort(sort.Reverse(items))
			} else {
				sort.Sort(items)
			}
		} else {
			if reverse {
				for i := 0; i < itemCount/2; i++ {
					items[i], items[itemCount-1-i] = items[itemCount-1-i], items[i]
				}
			}
		}

		if len(items) > 0 {
			for idx, item := range items {
				if !fn(idx, itemCount, item, nil) {
					return
				}
			}
		} else {
			empty()
		}
		return // done
	default:
		logf("Value.Iterate() not available for type: %T\n", val)
	}
	empty()
}

// Interface gives you access to the underlying value.
func (v *Value) Interface() any {
	return v.val
}

// EqualValueTo checks whether two values are containing the same value or object (if comparable).
func (v *Value) EqualValueTo(other *Value) bool {
	// comparison of uint with int fails using .Interface()-comparison (see issue #64)
	if v.IsInteger() && other.IsInteger() {
		return v.Integer() == other.Integer()
	}
	if v.IsTime() && other.IsTime() {
		return v.Time().Equal(other.Time())
	}
	if v.val == nil || other.val == nil {
		return v.val == other.val
	}

	// For simple comparable types, try direct comparison
	defer func() {
		// If comparison panics (uncomparable types like maps/slices), just return false
		recover()
	}()
	return v.val == other.val
}

type sortedKeys []reflect.Value

func (sk sortedKeys) Len() int {
	return len(sk)
}

func (sk sortedKeys) Less(i, j int) bool {
	vi := &Value{val: sk[i]}
	vj := &Value{val: sk[j]}
	switch {
	case vi.IsInteger() && vj.IsInteger():
		return vi.Integer() < vj.Integer()
	case vi.IsFloat() && vj.IsFloat():
		return vi.Float() < vj.Float()
	default:
		return vi.String() < vj.String()
	}
}

func (sk sortedKeys) Swap(i, j int) {
	sk[i], sk[j] = sk[j], sk[i]
}

type valuesList []*Value

func (vl valuesList) Len() int {
	return len(vl)
}

func (vl valuesList) Less(i, j int) bool {
	vi := vl[i]
	vj := vl[j]
	switch {
	case vi.IsInteger() && vj.IsInteger():
		return vi.Integer() < vj.Integer()
	case vi.IsFloat() && vj.IsFloat():
		return vi.Float() < vj.Float()
	default:
		return vi.String() < vj.String()
	}
}

func (vl valuesList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}
