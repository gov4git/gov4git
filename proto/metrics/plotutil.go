package metrics

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/gov4git/lib4git/form"
)

// utils

func writePyDateArray(w *bytes.Buffer, ds []time.Time) {
	w.WriteString("[")
	for _, date := range ds {
		fmt.Fprintf(w, "date(%d, %d, %d), ", date.Year(), date.Month(), date.Day())
	}
	w.WriteString("]")
}

func writePyIntArray(w *bytes.Buffer, vs []float64) {
	w.WriteString("[")
	for _, v := range vs {
		fmt.Fprintf(w, "%d, ", int(v))
	}
	w.WriteString("]")
}

func writePyFloatArray(w *bytes.Buffer, vs []float64) {
	w.WriteString("[")
	for _, v := range vs {
		fmt.Fprintf(w, "%f, ", v)
	}
	w.WriteString("]")
}

func generateRandomID() string {
	const w = 512 / 8 // 512 bits, measured in bytes
	buf := make([]byte, w)
	rand.Read(buf)
	return strings.ToLower(form.BytesHashForFilename(buf))
}
