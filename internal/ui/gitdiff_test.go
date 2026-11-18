package ui

import (
	"strings"
	"testing"
)

// ── Diff parsing: content lines that look like file headers ──────────────────

func TestParseDiffKeepsContentLinesThatLookLikeHeaders(t *testing.T) {
	raw := strings.Join([]string{
		"diff --git a/q.sql b/q.sql",
		"index c95ce1d..16ed999 100644",
		"--- a/q.sql",
		"+++ b/q.sql",
		"@@ -1,4 +1,4 @@",
		" select 1;",
		"--- removed sql comment",
		"+++ added line starting with two plus signs",
		"-index removed line starting with index",
		" select 2;",
		"diff --git a/b.go b/b.go",
		"index 1111111..2222222 100644",
		"--- a/b.go",
		"+++ b/b.go",
		"@@ -10,2 +10,2 @@",
		"-old",
		"+new",
		"",
	}, "\n")

	got := parseDiff(raw)

	want := []DiffLine{
		{Type: LineHunkHeader, Content: "@@ -1,4 +1,4 @@", OldNum: -1, NewNum: -1},
		{Type: LineContext, Content: "select 1;", OldNum: 1, NewNum: 1},
		{Type: LineRemoved, Content: "-- removed sql comment", OldNum: 2, NewNum: -1},
		{Type: LineAdded, Content: "++ added line starting with two plus signs", OldNum: -1, NewNum: 2},
		{Type: LineRemoved, Content: "index removed line starting with index", OldNum: 3, NewNum: -1},
		{Type: LineContext, Content: "select 2;", OldNum: 4, NewNum: 3},
		// The second file's header lines are still skipped after a hunk.
		{Type: LineHunkHeader, Content: "@@ -10,2 +10,2 @@", OldNum: -1, NewNum: -1},
		{Type: LineRemoved, Content: "old", OldNum: 10, NewNum: -1},
		{Type: LineAdded, Content: "new", OldNum: -1, NewNum: 10},
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d diff lines, got %d: %+v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: expected %+v, got %+v", i, want[i], got[i])
		}
	}
}
