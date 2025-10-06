package render

import (
	"io"
	"strings"
)

type Renderer interface {
	Render(v any) error
	Flush() error
}

func NewRenderer(format string, w io.Writer) Renderer {
	switch strings.ToLower(format) {
	case "json", "jsonl", "json-lines", "json-compact":
		return NewJSONLineRenderer(w)
	case "json-pretty":
		return NewJSONRenderer(w, true)
	case "yaml", "yml":
		return NewYAMLRenderer(w)
	case "xml":
		return NewXMLRenderer(w)
	case "csv":
		return NewCSVRenderer(w)
	case "text", "pretty", "human":
		return NewTextRenderer(w)
	case "plain", "shell":
		return NewLineRenderer(w)
	default:
		return NewTextRenderer(w)
	}
}
