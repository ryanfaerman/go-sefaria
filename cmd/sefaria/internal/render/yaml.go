package render

import (
	"io"

	"gopkg.in/yaml.v3"
)

type YAMLRenderer struct {
	w io.Writer
}

func NewYAMLRenderer(w io.Writer) *YAMLRenderer {
	return &YAMLRenderer{w: w}
}

func (r *YAMLRenderer) Render(v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = r.w.Write(data)
	return err
}

func (r *YAMLRenderer) Flush() error { return nil }
