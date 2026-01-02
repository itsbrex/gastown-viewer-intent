// Package beads provides integration with the Beads issue tracker via the bd CLI.
package beads

import (
	"context"

	"github.com/intent-solutions-io/gastown-viewer-intent/internal/model"
)

// Adapter defines the interface for interacting with Beads.
// All methods shell out to the bd CLI and parse JSON output.
type Adapter interface {
	// ListIssues returns all issues matching the optional filter.
	ListIssues(ctx context.Context, filter model.IssueFilter) ([]model.Issue, error)

	// GetIssue returns a single issue by ID with full details.
	GetIssue(ctx context.Context, id string) (*model.Issue, error)

	// Board returns issues grouped by status for board view.
	Board(ctx context.Context) (*model.Board, error)

	// Graph returns the dependency graph.
	Graph(ctx context.Context) (*model.Graph, error)

	// IsInitialized checks if beads is initialized in the current directory.
	IsInitialized(ctx context.Context) (bool, error)

	// Version returns the bd CLI version.
	Version(ctx context.Context) (string, error)
}

// CLIAdapter implements Adapter by shelling out to the bd CLI.
type CLIAdapter struct {
	executor Executor
	workDir  string
}

// NewCLIAdapter creates a new CLI-based adapter.
// If workDir is empty, uses the current directory.
func NewCLIAdapter(workDir string) *CLIAdapter {
	return &CLIAdapter{
		executor: &DefaultExecutor{},
		workDir:  workDir,
	}
}

// NewCLIAdapterWithExecutor creates an adapter with a custom executor (for testing).
func NewCLIAdapterWithExecutor(workDir string, executor Executor) *CLIAdapter {
	return &CLIAdapter{
		executor: executor,
		workDir:  workDir,
	}
}

// ListIssues implements Adapter.ListIssues.
func (a *CLIAdapter) ListIssues(ctx context.Context, filter model.IssueFilter) ([]model.Issue, error) {
	args := []string{"list", "--json"}

	if filter.Status != "" {
		args = append(args, "--status", filter.Status)
	}

	output, err := a.executor.Execute(ctx, a.workDir, args...)
	if err != nil {
		return nil, err
	}

	bdIssues, err := ParseIssueList(output)
	if err != nil {
		return nil, &ParseError{Command: "list", Err: err}
	}

	issues := make([]model.Issue, 0, len(bdIssues))
	for _, bi := range bdIssues {
		issues = append(issues, bi.ToModelIssue())
	}

	return issues, nil
}

// GetIssue implements Adapter.GetIssue.
func (a *CLIAdapter) GetIssue(ctx context.Context, id string) (*model.Issue, error) {
	output, err := a.executor.Execute(ctx, a.workDir, "show", id, "--json")
	if err != nil {
		if IsNotFoundError(err) {
			return nil, &NotFoundError{ID: id}
		}
		return nil, err
	}

	bdIssues, err := ParseIssueList(output)
	if err != nil {
		return nil, &ParseError{Command: "show", Err: err}
	}

	if len(bdIssues) == 0 {
		return nil, &NotFoundError{ID: id}
	}

	issue := bdIssues[0].ToModelIssue()
	return &issue, nil
}

// Board implements Adapter.Board.
func (a *CLIAdapter) Board(ctx context.Context) (*model.Board, error) {
	issues, err := a.ListIssues(ctx, model.NewIssueFilter())
	if err != nil {
		return nil, err
	}

	board := model.NewBoard()
	for _, issue := range issues {
		board.AddIssue(model.IssueSummary{
			ID:       issue.ID,
			Title:    issue.Title,
			Status:   issue.Status,
			Priority: issue.Priority,
		})
	}

	return &board, nil
}

// Graph implements Adapter.Graph.
func (a *CLIAdapter) Graph(ctx context.Context) (*model.Graph, error) {
	// Get all issues for nodes
	issues, err := a.ListIssues(ctx, model.NewIssueFilter())
	if err != nil {
		return nil, err
	}

	graph := model.NewGraph()

	// Build node map and add nodes
	nodeMap := make(map[string]bool)
	for _, issue := range issues {
		graph.AddNode(model.GraphNode{
			ID:       issue.ID,
			Title:    issue.Title,
			Status:   issue.Status,
			Priority: issue.Priority,
		})
		nodeMap[issue.ID] = true
	}

	// Get blocked info for edges
	output, err := a.executor.Execute(ctx, a.workDir, "blocked", "--json")
	if err != nil {
		// blocked command may fail if no blocked issues - that's ok
		return &graph, nil
	}

	blockedIssues, err := ParseBlockedList(output)
	if err != nil {
		// Graceful degradation - return graph without edges
		return &graph, nil
	}

	// Add edges from blocked relationships
	for _, bi := range blockedIssues {
		for _, blockerID := range bi.BlockedBy {
			if nodeMap[blockerID] && nodeMap[bi.ID] {
				graph.AddEdge(model.GraphEdge{
					From: blockerID,
					To:   bi.ID,
					Type: model.EdgeTypeBlocks,
				})
			}
		}
	}

	return &graph, nil
}

// IsInitialized implements Adapter.IsInitialized.
func (a *CLIAdapter) IsInitialized(ctx context.Context) (bool, error) {
	_, err := a.executor.Execute(ctx, a.workDir, "status")
	if err != nil {
		if IsNotInitializedError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Version implements Adapter.Version.
func (a *CLIAdapter) Version(ctx context.Context) (string, error) {
	output, err := a.executor.Execute(ctx, a.workDir, "--version")
	if err != nil {
		return "", err
	}
	return ParseVersion(output), nil
}
