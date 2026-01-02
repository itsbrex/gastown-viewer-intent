package model

// EdgeType represents the type of dependency relationship.
type EdgeType string

const (
	EdgeTypeBlocks EdgeType = "blocks"
	EdgeTypeParent EdgeType = "parent"
)

// GraphNode represents a node in the dependency graph.
type GraphNode struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Status   Status   `json:"status"`
	Priority Priority `json:"priority"`
}

// GraphEdge represents a directed edge in the dependency graph.
type GraphEdge struct {
	From string   `json:"from"`
	To   string   `json:"to"`
	Type EdgeType `json:"type"`
}

// GraphStats contains statistics about the graph.
type GraphStats struct {
	NodeCount int `json:"node_count"`
	EdgeCount int `json:"edge_count"`
	MaxDepth  int `json:"max_depth"`
}

// Graph represents the full dependency graph.
type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
	Stats GraphStats  `json:"stats"`
}

// NewGraph creates an empty graph.
func NewGraph() Graph {
	return Graph{
		Nodes: []GraphNode{},
		Edges: []GraphEdge{},
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(node GraphNode) {
	g.Nodes = append(g.Nodes, node)
	g.Stats.NodeCount++
}

// AddEdge adds an edge to the graph.
func (g *Graph) AddEdge(edge GraphEdge) {
	g.Edges = append(g.Edges, edge)
	g.Stats.EdgeCount++
}

// GraphFormat specifies output format for the graph endpoint.
type GraphFormat string

const (
	GraphFormatJSON GraphFormat = "json"
	GraphFormatDOT  GraphFormat = "dot"
)
