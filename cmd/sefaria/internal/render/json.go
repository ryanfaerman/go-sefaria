package render

import (
	"encoding/json"
	"io"
)

type JSONLineRenderer struct {
	w   io.Writer
	enc *json.Encoder
}

func NewJSONLineRenderer(w io.Writer) *JSONLineRenderer {
	return &JSONLineRenderer{
		w:   w,
		enc: json.NewEncoder(w),
	}
}

func (r *JSONLineRenderer) Render(v any) error {
	return r.enc.Encode(v) // each object on a new line
}

func (r *JSONLineRenderer) Flush() error { return nil }

type JSONRenderer struct {
	w      io.Writer
	pretty bool
}

func NewJSONRenderer(w io.Writer, pretty bool) *JSONRenderer {
	return &JSONRenderer{w: w, pretty: pretty}
}

func (r *JSONRenderer) Render(v any) error {
	var data []byte
	var err error
	if r.pretty {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}
	_, err = r.w.Write(data)
	if err != nil {
		return err
	}
	_, _ = r.w.Write([]byte("\n"))
	return nil
}

func (r *JSONRenderer) Flush() error { return nil }
