package gwc

import "fmt"

func NewNodeEnvironment(nodes Nodes) *NodeEnvironment {
	nodes_map := NodesMap{}
	for _, node := range nodes {
		nodes_map[node.ID()] = node
	}

	var current NodeID
	if len(nodes) > 0 {
		current = nodes[0].ID()
	}

	return &NodeEnvironment{
		Current:      current,
		Nodes:        nodes,
		NodesMap:     nodes_map,
		CollapsedMap: NodeCollapsedMap{},
		StatesMap:    NodeStatesMap{},
	}
}

type NodeEnvironment struct {
	Nodes
	NodesMap
	Current      NodeID
	CollapsedMap NodeCollapsedMap
	StatesMap    NodeStatesMap
}

type NodeStates = []NodeState
type NodeStatesMap = map[NodeID]NodeState
type NodeCollapsedMap = map[NodeID]int
type NodeFilterFn = func(NodeID, NodeState) bool
type NodeIDsOrNodeFilterFn = interface{}

func (ne *NodeEnvironment) ToID(idx int) NodeID {
	if len(ne.Nodes) > idx {

		if node := ne.Nodes[idx]; node != nil {
			return node.ID()
		}
	}
	return ""
}

func (ne *NodeEnvironment) ToIndex(id NodeID) int {
	for idx, node := range ne.Nodes {
		if id == node.ID() {
			return idx
		}
	}
	return -1
}

func (ne *NodeEnvironment) States() NodeStates {
	states := make(NodeStates, len(ne.Nodes))
	for idx, node := range ne.Nodes {
		states[idx] = ne.StatesMap[node.ID()]
	}
	return states
}

func (ne *NodeEnvironment) Collapsed() NodeIDs {
	ids := make(NodeIDs, len(ne.CollapsedMap))
	for id, at := range ne.CollapsedMap {
		ids[at] = id
	}
	return ids
}

func (ne *NodeEnvironment) CollapsedNodes() Nodes {
	nodes := make(Nodes, len(ne.CollapsedMap))
	for id, at := range ne.CollapsedMap {
		nodes[at] = ne.NodesMap[id]
	}
	return nodes
}

func (ne *NodeEnvironment) IsNeighbour(other NodeID) bool {
	return ne.IsNeighbourOf(ne.Current, other)
}
func (ne *NodeEnvironment) IsNeighbourOf(id, other NodeID) bool {
	for _, ni := range ne.NodesMap[id].Neighbours() {
		if ni == other {
			return true
		}
	}
	return false
}

func (ne *NodeEnvironment) IsWithinRange(other NodeID, depth uint) bool {
	return ne.IsWithinRangeOf(ne.Current, other, depth)
}
func (ne *NodeEnvironment) IsWithinRangeOf(a, b NodeID, depth uint) bool {
	for _, ni := range ne.aggregateNodeNeighbours(NodeIDs{}, a, depth) {
		if ni == b {
			return true
		}
	}
	return false
}

func (ne *NodeEnvironment) NodesWithinRangeExcl(depth uint) NodeIDs {
	return ne.NodesWithinRangeOfExcl(ne.Current, depth)
}
func (ne *NodeEnvironment) NodesWithinRangeIncl(depth uint) NodeIDs {
	return ne.NodesWithinRangeOfIncl(ne.Current, depth)
}
func (ne *NodeEnvironment) NodesWithinRangeOfExcl(id NodeID, depth uint) NodeIDs {
	return ne.NodesWithinRangeOfIncl(id, depth)[1:]
}
func (ne *NodeEnvironment) NodesWithinRangeOfIncl(id NodeID, depth uint) NodeIDs {
	return ne.aggregateNodeNeighbours(NodeIDs{id}, id, depth)
}
func (ne *NodeEnvironment) aggregateNodeNeighbours(ids NodeIDs, id NodeID, depth uint) NodeIDs {
	// Stop if we're looking for a range smaller than the logical minimum.
	if depth < 1 {
		return ids
	}

	// Find all neighbours that aren't in the list yet.
	neighbours := NodeIDs{}
Outer:
	for _, ni := range ne.NodesMap[id].Neighbours() {
		for _, nc := range ids {
			if nc == ni {
				continue Outer
			}
		}
		neighbours = append(neighbours, ni)
	}

	// Add the found new neighbours to the list of nodes in range.
	ids = append(ids, neighbours...)

	// Continue aggregation with the neighbours that hadn't been in the list yet.
	for _, ni := range neighbours {
		ids = ne.aggregateNodeNeighbours(ids, ni, depth-1)
	}

	return ids
}

func (ne *NodeEnvironment) FilterNodes(fn_or_ids NodeIDsOrNodeFilterFn) NodeIDs {
	switch v := fn_or_ids.(type) {
	case []NodeID:
		return NodeIDs(v)
	case NodeIDs:
		return v
	case NodeFilterFn:
		filtered := NodeIDs{}
		for _, node := range ne.Nodes {
			id := node.ID()
			if v(id, node) {
				filtered = append(filtered, id)
			}
		}
		return filtered
	}

	panic(fmt.Sprintf("FilterNodes() cannot handle this type: %#v", fn_or_ids))
}

func (ne *NodeEnvironment) FilterNodesAnd(fn1 NodeIDsOrNodeFilterFn, fns ...NodeIDsOrNodeFilterFn) NodeIDs {
	filtered := ne.FilterNodes(fn1)
	for _, fn := range fns {
		filtered = filtered.And(ne.FilterNodes(fn))
	}
	return filtered
}

func (ne *NodeEnvironment) FilterNodesOr(fn1 NodeIDsOrNodeFilterFn, fns ...NodeIDsOrNodeFilterFn) NodeIDs {
	filtered := ne.FilterNodes(fn1)
	for _, fn := range fns {
		filtered = filtered.Or(ne.FilterNodes(fn))
	}
	return filtered
}
