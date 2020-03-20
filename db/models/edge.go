package models

// Edge represents an edge between two nodes in the neo4j database
type Edge struct {
	TimeSpent int64
}

// GenerateEdgeFromInt generates an edge struct from timeSpent integer
func GenerateEdgeFromInt(t int64) Edge {
	edge := Edge{TimeSpent: t}
	return edge
}
