package render

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type CSVRenderer struct {
	w              io.Writer
	cw             *csv.Writer
	headersWritten bool
}

func NewCSVRenderer(w io.Writer) *CSVRenderer {
	return &CSVRenderer{
		w:  w,
		cw: csv.NewWriter(w),
	}
}

func (r *CSVRenderer) Render(v any) error {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return nil
		}
		elem := val.Index(0)
		switch elem.Kind() {
		case reflect.Struct:
			if !r.headersWritten {
				headers := structFieldNames(elem.Type())
				if err := r.cw.Write(headers); err != nil {
					return err
				}
				r.headersWritten = true
			}
			for i := 0; i < val.Len(); i++ {
				row := structToStringSlice(val.Index(i))
				if err := r.cw.Write(row); err != nil {
					return err
				}
			}
		default:
			for i := 0; i < val.Len(); i++ {
				if err := r.cw.Write([]string{fmt.Sprint(val.Index(i).Interface())}); err != nil {
					return err
				}
			}
		}
	default:
		return r.cw.Write([]string{fmt.Sprint(v)})
	}
	return nil
}

func (r *CSVRenderer) Flush() error {
	r.cw.Flush()
	return r.cw.Error()
}

// Helpers
func structFieldNames(t reflect.Type) []string {
	var names []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("csv") == "-" {
			continue
		}
		if tag := f.Tag.Get("csv"); tag != "" {
			names = append(names, tag)
		} else if tag := f.Tag.Get("json"); tag != "" {
			names = append(names, strings.Split(tag, ",")[0])
		} else {
			names = append(names, f.Name)
		}
	}
	return names
}

func structToStringSlice(v reflect.Value) []string {
	t := v.Type()
	row := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("csv") == "-" {
			continue
		}
		row = append(row, fmt.Sprint(v.Field(i).Interface()))
	}
	return row
}
