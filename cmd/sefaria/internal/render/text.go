package render

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/mattn/go-runewidth"
)

type TextRenderer struct {
	w io.Writer
}

func NewTextRenderer(w io.Writer) *TextRenderer {
	return &TextRenderer{w: w}
}

func (r *TextRenderer) Render(v any) error {
	return r.renderValue(reflect.ValueOf(v), 0)
}

func (r *TextRenderer) Flush() error { return nil }

func (r *TextRenderer) renderValue(v reflect.Value, level int) error {
	indent := strings.Repeat("  ", level)

	if !v.IsValid() {
		fmt.Fprintf(r.w, "%s<nil>\n", indent)
		return nil
	}

	// unwrap pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			fmt.Fprintf(r.w, "%s<nil>\n", indent)
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" { // unexported
				continue
			}
			val := v.Field(i)
			if val.Kind() == reflect.Slice && val.Len() > 0 && val.Index(0).Kind() == reflect.Struct {
				// Render as table if slice of struct
				fmt.Fprintf(r.w, "%s%s:\n", indent, f.Name)
				r.renderStructSliceTable(val, level+1)
				continue
			}
			if val.Kind() == reflect.Struct || val.Kind() == reflect.Slice || val.Kind() == reflect.Map || val.Kind() == reflect.Ptr {
				fmt.Fprintf(r.w, "%s%s:\n", indent, f.Name)
				r.renderValue(val, level+1)
			} else if val.Kind() == reflect.String && val.String() == "" {
				continue
			} else {
				fmt.Fprintf(r.w, "%s%s: %v\n", indent, f.Name, val.Interface())
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() == reflect.Struct {
				fmt.Fprintf(r.w, "%s- \n", indent)
				r.renderValue(elem, level+1)
			} else {
				fmt.Fprintf(r.w, "%s- %v\n", indent, elem.Interface())
			}
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			fmt.Fprintf(r.w, "%s%v: ", indent, key.Interface())
			if val.Kind() == reflect.Struct || val.Kind() == reflect.Slice || val.Kind() == reflect.Map {
				fmt.Fprintln(r.w)
				r.renderValue(val, level+1)
			} else {
				fmt.Fprintf(r.w, "%v\n", val.Interface())
			}
		}

	default:
		fmt.Fprintf(r.w, "%s%v\n", indent, v.Interface())
	}

	return nil
}

// renderStructSliceTable renders a slice of structs as an aligned table using table tags
func (r *TextRenderer) renderStructSliceTable(slice reflect.Value, level int) {
	if slice.Len() == 0 {
		return
	}
	indent := strings.Repeat("  ", level)

	t := slice.Index(0).Type()
	headers := []string{}
	headerNames := []string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("table")
		if tag == "-" || tag == "" {
			continue
		}
		headers = append(headers, f.Name)
		headerNames = append(headerNames, tag)
	}

	if len(headers) == 0 {
		return
	}

	// Determine column widths
	colWidths := make([]int, len(headers))
	for i, name := range headerNames {
		colWidths[i] = runewidth.StringWidth(name)
	}
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		for j, field := range headers {
			val := elem.FieldByName(field)
			str := fmt.Sprint(val.Interface())
			l := runewidth.StringWidth(str)
			if l > colWidths[j] {
				colWidths[j] = l
			}
		}
	}

	// Print headers using ANSI positioning
	fmt.Fprint(r.w, indent)
	currentPos := len(indent)
	for i, name := range headerNames {
		fmt.Fprint(r.w, name)
		currentPos += runewidth.StringWidth(name)
		// Move to next column position
		nextPos := currentPos + colWidths[i] - runewidth.StringWidth(name) + 2
		if i < len(headerNames)-1 {
			fmt.Fprintf(r.w, "\033[%dC", nextPos-currentPos)
			currentPos = nextPos
		}
	}
	fmt.Fprintln(r.w)

	// Print separator using ANSI positioning
	fmt.Fprint(r.w, indent)
	currentPos = len(indent)
	for i, w := range colWidths {
		fmt.Fprint(r.w, strings.Repeat("-", w))
		currentPos += w
		// Move to next column position
		nextPos := currentPos + 2
		if i < len(colWidths)-1 {
			fmt.Fprintf(r.w, "\033[%dC", nextPos-currentPos)
			currentPos = nextPos
		}
	}
	fmt.Fprintln(r.w)

	// Print rows using ANSI positioning
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		fmt.Fprint(r.w, indent)
		currentPos := len(indent)
		for j, field := range headers {
			val := elem.FieldByName(field)
			cellContent := fmt.Sprint(val.Interface())
			fmt.Fprint(r.w, cellContent)
			currentPos += runewidth.StringWidth(cellContent)
			// Move to next column position
			nextPos := currentPos + colWidths[j] - runewidth.StringWidth(cellContent) + 2
			if j < len(headers)-1 {
				fmt.Fprintf(r.w, "\033[%dC", nextPos-currentPos)
				currentPos = nextPos
			}
		}
		fmt.Fprintln(r.w)
	}
}
