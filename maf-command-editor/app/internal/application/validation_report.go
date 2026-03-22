package application

import "strings"

func (r ValidationReport) String() string {
	if r.OK {
		return "ok"
	}
	lines := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
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
