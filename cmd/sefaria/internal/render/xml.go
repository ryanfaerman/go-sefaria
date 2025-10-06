package render

import (
	"encoding/xml"
	"io"
)

type XMLRenderer struct {
	w io.Writer
}

func NewXMLRenderer(w io.Writer) *XMLRenderer {
	return &XMLRenderer{w: w}
}

func (r *XMLRenderer) Render(v any) error {
	data, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = r.w.Write(append(data, '\n'))
	return err
}

func (r *XMLRenderer) Flush() error { return nil }
