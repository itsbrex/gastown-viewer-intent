package model

// Column represents a status column in the board view.
type Column struct {
	Status Status         `json:"status"`
	Label  string         `json:"label"`
	Count  int            `json:"count"`
	Issues []IssueSummary `json:"issues"`
}

// Board represents the kanban board view with issues grouped by status.
type Board struct {
	Columns []Column `json:"columns"`
	Total   int      `json:"total"`
}

// NewBoard creates an empty board with standard columns.
func NewBoard() Board {
	return Board{
		Columns: []Column{
			{Status: StatusPending, Label: "Pending", Issues: []IssueSummary{}},
			{Status: StatusInProgress, Label: "In Progress", Issues: []IssueSummary{}},
			{Status: StatusDone, Label: "Done", Issues: []IssueSummary{}},
			{Status: StatusBlocked, Label: "Blocked", Issues: []IssueSummary{}},
		},
	}
}

// AddIssue adds an issue summary to the appropriate column.
func (b *Board) AddIssue(issue IssueSummary) {
	for i := range b.Columns {
		if b.Columns[i].Status == issue.Status {
			b.Columns[i].Issues = append(b.Columns[i].Issues, issue)
			b.Columns[i].Count++
			b.Total++
			return
		}
	}
}
