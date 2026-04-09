package history

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Print writes a human-readable table of history entries to w.
func Print(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No history entries found.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTIMESTAMP\tOLD FILE\tNEW FILE\t+\t-\t~\tFINDINGS")
	fmt.Fprintln(tw, "--\t---------\t--------\t--------\t-\t-\t-\t--------")

	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\n",
			shortID(e.ID),
			e.Timestamp.Format("2006-01-02 15:04"),
			baseName(e.OldFile),
			baseName(e.NewFile),
			e.Added,
			e.Removed,
			e.Modified,
			e.Findings,
		)
	}
	tw.Flush()
}

func shortID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

func baseName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}
