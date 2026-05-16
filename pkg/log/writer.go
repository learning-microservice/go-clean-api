package log

import (
	"bytes"
	"encoding/json"
	"io"
)

type prettyJSONWriter struct {
	out    io.Writer
	indent string
}

func (w *prettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, p, "", w.indent); err != nil {
		return w.out.Write(p)
	}

	return w.out.Write(prettyJSON.Bytes())
}
