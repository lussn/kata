// Package diameter contains functionality to measure the diameter of a graph
// of connected nodes http://mathworld.wolfram.com/GraphDiameter.html
// The diameter is the length of the longest shortest path in the graph.
// Two points in a graph have a shortest path between them.
// One pair of points will have a longer shortest path than others.
// The length of this path is the diameter.
// Any pair of points will be at most this far apart.
package diameter

import (
	"container/list"
)

// nodeID is an unique identifier for each node
type nodeID int32

// nodeName is the name of the node looked up by id from the symbol table.
type nodeName string

// symbolTable contains the mapping from id to name.
type symbolTable map[nodeName]nodeID

// getID returns the id of the node with name if it exists, otherwise it adds
// the name to the table and returns it.
func (s symbolTable) getID(name nodeName) nodeID {
	id, ok := s[name]
	if !ok {
		id = nodeID(len(s))
		s[name] = id
	}
	return id
}

// Graph is the complete graph containing the lookup table for node names and
// the actual nodes graph.
type Graph struct {
	symbolTable
	nodes
}

// New returns a new graph.
func New() *Graph {
	return &Graph{
		symbolTable: make(symbolTable),
		nodes:       make(nodes),
	}
}

// addEdge adds a connection between node a and b identified by their name.
// It retrieves the nodes from the lookup table to get ids.
func (g *Graph) addEdge(a, b nodeName) {
	aid := g.symbolTable.getID(a)
	bid := g.symbolTable.getID(b)

	g.nodes.addEdge(aid, bid)
}

// node represents one node in the graph, identified by it's id.
// A node knows about all adjacent nodes.
type node struct {
	id nodeID

	// adjacent edges
	adj map[nodeID]*node
}

// add adds an adjacent neighbor node for the node.
func (n *node) add(adjNode *node) {
	n.adj[adjNode.id] = adjNode
}

// nodes represents the graph of nodes.
type nodes map[nodeID]*node

// get retrieves one node by it's id.
// if the id is not present in the graph it is added and returned.
func (nodes nodes) get(id nodeID) *node {
	n, ok := nodes[id]
	if !ok {
		n = &node{
			id:  id,
			adj: make(map[nodeID]*node),
		}
		nodes[id] = n
	}
	return n
}

// addEdge adds a connection between node a and b identified by their id.
// it adds retrieves/adds the nodes and makes the connection between them, i.e.
// adding them as adjacent nodes.
func (nodes *nodes) addEdge(a, b nodeID) {
	an := nodes.get(a)
	bn := nodes.get(b)

	an.add(bn)
	bn.add(an)
}

// diameter returns the maximum length of a shortest path in the graph.
func (nodes nodes) diameter() int {
	var diameter int
	for id := range nodes {
		df := nodes.longestShortestPath(id)
		if df > diameter {
			diameter = df
		}
	}
	return diameter
}

// bfsNode is used to keep track of nodes in the Breadth First Search.
type bfsNode struct {
	parent *node
	depth  int
}

// longestShortestPath executes the BFS from the start node identified by id.
// Returns the depth of the BFS which is the longest minimum distance between
// nodes in the graph.
func (nodes nodes) longestShortestPath(start nodeID) int {
	q := list.New()

	bfsData := make(map[nodeID]bfsNode, len(nodes))

	n := nodes.get(start)
	bfsData[n.id] = bfsNode{parent: n, depth: 0}
	q.PushBack(n)

	for {
		elt := q.Front()
		if elt == nil {
			break
		}
		n = q.Remove(elt).(*node)

		for id, m := range n.adj {
			bm := bfsData[id]
			if bm.parent == nil {
				bfsData[id] = bfsNode{parent: n, depth: bfsData[n.id].depth + 1}
				q.PushBack(m)
			}
		}
	}

	return bfsData[n.id].depth
}
