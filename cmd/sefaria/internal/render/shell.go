package render

import (
	"fmt"
	"io"
	"reflect"
)

type LineRenderer struct {
	w io.Writer
}

func NewLineRenderer(w io.Writer) *LineRenderer {
	return &LineRenderer{w: w}
}

func (r *LineRenderer) Render(v any) error {
	return r.renderValue(reflect.ValueOf(v))
}

func (r *LineRenderer) Flush() error { return nil }

func (r *LineRenderer) renderValue(v reflect.Value) error {
	if !v.IsValid() {
		return nil
	}

	// unwrap pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() == reflect.String {
				fmt.Fprintln(r.w, elem.String())
			} else {
				r.renderValue(elem)
			}
		}

	case reflect.Struct:
		// optional: extract a specific field if needed
		// e.g., a struct with Title, Key, etc.
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			val := v.Field(i)
			// if you want just Title
			if f.Name == "Title" && val.Kind() == reflect.String {
				fmt.Fprintln(r.w, val.String())
			}
		}

	case reflect.String:
		fmt.Fprintln(r.w, v.String())
	}

	return nil
}
