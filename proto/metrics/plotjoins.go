package metrics

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/lib4git/must"
)

func plotDailyJoinsPNG(
	ctx context.Context,
	series *Series,

) []byte {

	// prepare py plotting program

	var w bytes.Buffer

	w.WriteString(`
import matplotlib.pyplot as plt
import numpy as np
from datetime import date

x = `)
	writePyDateArray(&w, series.DailyNumJoins.X)
	fmt.Fprintln(&w)

	w.WriteString(`y1 = np.array(`)
	writePyIntArray(&w, series.DailyNumJoins.Y)
	fmt.Fprintln(&w, ")")

	w.WriteString(
		`fig, ax = plt.subplots(figsize=(9, 5))
ax.bar(x, y1, color='#77eeaa')
`)

	w.WriteString(
		`ax.set_xlabel("Days")
ax.set_ylabel("Count")
ax.legend(["New members", ])
ax.set_title("Daily count of new community members")
`)

	n := series.DailyNumJoins.Len()
	fmt.Fprintf(&w, "ax.set_xticks(x[0::%d])\n", xTickSkipDates(n))

	fp := filepath.Join(os.TempDir(), generateRandomID()+".png")
	fmt.Fprintf(&w, "fig.savefig(%q, dpi=200, bbox_inches = 'tight')\n", fp)

	py := w.String()
	fmt.Println(py)

	// call python
	outerr, err := runtime.RunPython(ctx, py)
	fmt.Println(string(outerr))
	must.NoError(ctx, err)

	// retrieve png plot
	pngData, err := os.ReadFile(fp)
	must.NoError(ctx, err)

	return pngData
}
