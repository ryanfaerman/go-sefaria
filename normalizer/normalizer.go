package normalizer

import (
	"reflect"
)

// Normalizer is a function that takes a string and returns a normalized string.
// Normalizers are stateless functions that transform text in a consistent way.
type Normalizer func(string) string

// Apply recursively applies the given normalizers to all string fields in the provided value.
// The function uses reflection to traverse complex data structures and applies normalizers
// to string fields found within structs, slices, arrays, maps, and pointers.
//
// Apply modifies the input value in-place. If the value is nil, the function returns early.
// For maps, only string values are normalized directly; other value types are processed recursively.
//
// Example:
//
//	type Person struct {
//		Name string
//		Bio  string
//	}
//
//	p := &Person{Name: "John &amp; Jane", Bio: "Hello "world""}
//	normalizer.Apply(p, normalizer.HTMLUnescape, normalizer.Punctuation)
//	// p.Name is now "John & Jane", p.Bio is now "Hello "world""
func Apply(v any, normalizers ...Normalizer) {
	if v == nil {
		return
	}
	rv := reflect.ValueOf(v)
	applyValue(rv, normalizers)
}

func applyValue(rv reflect.Value, normalizers []Normalizer) {
	if !rv.IsValid() {
		return
	}

	switch rv.Kind() {
	case reflect.Ptr:
		if !rv.IsNil() {
			applyValue(rv.Elem(), normalizers)
		}
	case reflect.Interface:
		if !rv.IsNil() {
			elem := rv.Elem()
			if elem.Kind() == reflect.String {
				str := elem.String()
				for _, n := range normalizers {
					str = n(str)
				}
				rv.Set(reflect.ValueOf(str))
			} else {
				applyValue(elem, normalizers)
			}
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			if field.CanSet() {
				applyValue(field, normalizers)
			} else if field.Kind() == reflect.Struct || field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
				applyValue(field, normalizers)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			applyValue(rv.Index(i), normalizers)
		}
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			val := rv.MapIndex(key)
			if val.Kind() == reflect.String {
				str := val.String()
				for _, n := range normalizers {
					str = n(str)
				}
				rv.SetMapIndex(key, reflect.ValueOf(str))
			} else if val.Kind() == reflect.Interface {
				// Handle interface{} values that might contain strings
				if val.Elem().Kind() == reflect.String {
					str := val.Elem().String()
					for _, n := range normalizers {
						str = n(str)
					}
					rv.SetMapIndex(key, reflect.ValueOf(str))
				} else {
					applyValue(val, normalizers)
				}
			} else {
				applyValue(val, normalizers)
			}
		}
	case reflect.String:
		str := rv.String()
		for _, n := range normalizers {
			str = n(str)
		}
		rv.SetString(str)
	}
}
