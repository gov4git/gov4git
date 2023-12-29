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

func plotDailyCreditsPNG(
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
	writePyDateArray(&w, series.DailyCreditsIssued.X)
	fmt.Fprintln(&w)

	w.WriteString(`y1 = np.array(`)
	writePyIntArray(&w, series.DailyCreditsIssued.Y)
	fmt.Fprintln(&w, ")")

	w.WriteString(`y2 = np.array(`)
	writePyIntArray(&w, series.DailyCreditsBurned.Y)
	fmt.Fprintln(&w, ")")

	w.WriteString(`y3 = np.array(`)
	writePyIntArray(&w, series.DailyCreditsTransferred.Y)
	fmt.Fprintln(&w, ")")

	w.WriteString(
		`fig, ax = plt.subplots(figsize=(9, 5))
ax.bar(x, y1, color='#5599cc')
ax.bar(x, y2, bottom=y1, color='#cc5599')
ax.bar(x, y3, bottom=y1+y2, color='#aabbcc')
`)

	w.WriteString(
		`ax.set_xlabel("Days")
ax.set_ylabel("Credits")
ax.legend(["Issued", "Burned", "Transferred", ])
ax.set_title("Daily credits issued/burned/transferred")
`)

	n := series.DailyCreditsIssued.Len()
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
