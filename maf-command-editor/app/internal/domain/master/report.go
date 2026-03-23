package master

import (
	"sort"
	"strings"
)

type Counts struct {
	Items       int
	Grimoire    int
	Skills      int
	EnemySkills int
	Enemies     int
	SpawnTables int
	Treasures   int
	LootTables  int
}

type ValidationIssue struct {
	Entity  string
	ID      string
	Field   string
	Message string
}

type ValidationReport struct {
	OK     bool
	Counts Counts
	Issues []ValidationIssue
}

func (r ValidationReport) Sorted() ValidationReport {
	out := r
	out.Issues = append([]ValidationIssue{}, r.Issues...)
	sort.Slice(out.Issues, func(i, j int) bool {
		left := out.Issues[i]
		right := out.Issues[j]
		if left.Entity != right.Entity {
			return left.Entity < right.Entity
		}
		if left.ID != right.ID {
			return left.ID < right.ID
		}
		if left.Field != right.Field {
			return left.Field < right.Field
		}
		return left.Message < right.Message
	})
	return out
}

func (r ValidationReport) String() string {
	if r.OK {
		return "ok"
	}
	sorted := r.Sorted()
	lines := make([]string, 0, len(sorted.Issues))
	for _, issue := range sorted.Issues {
		label := issue.Entity
		if issue.ID != "" {
			label += "[" + issue.ID + "]"
		}
		if issue.Field != "" {
			label += "." + issue.Field
		}
		lines = append(lines, label+": "+issue.Message)
	}
	return strings.Join(lines, "\n")
}
