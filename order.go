package gwc

import (
	"math/rand"
)

// CollapseOrderFn is a function that takes the current NodeEnvironment and returns the NodeID of the next Node to be collapsed.
type CollapseOrderFn func(*rand.Rand, NodeEnvironment) NodeID

var RandomCollapseOrder CollapseOrderFn = func(rnd *rand.Rand, env NodeEnvironment) NodeID {
FindRandomIndex:
	for _, idx := range rnd.Perm(len(env.Nodes)) {
		id := env.GetID(idx)
		for cid, _ := range env.CollapsedMap {
			if id == cid {
				continue FindRandomIndex
			}
		}
		return id
	}
	return ""
}

var NeighbourhoodCollapseOrder CollapseOrderFn = func(rnd *rand.Rand, env NodeEnvironment) NodeID {
	// If this is not the first run, search for uncollapsed neighbours of the current node.
	if env.Current != "" {
		neighbours := env.NodesMap[env.Current].Neighbours()
	FindRandomNeighbour:
		for _, idx := range rnd.Perm(len(neighbours)) {
			nid := neighbours[idx]
			for cid, _ := range env.CollapsedMap {
				if nid == cid {
					continue FindRandomNeighbour
				}
			}
			return nid
		}
	}
	// If this is the first run, or the current node has run out of neighbours, pick a new random node to continue with.
	if len(env.Nodes) > 0 {
		return RandomCollapseOrder(rnd, env)
	}
	return ""
}

var AscendingCollapseOrder CollapseOrderFn = func(rnd *rand.Rand, env NodeEnvironment) NodeID {
	if len(env.CollapsedMap) < len(env.Nodes) {
		idx := len(env.CollapsedMap)
		node := env.Nodes[idx]
		if node != nil {
			return node.ID()
		}
	}
	return ""
}

var DescendingCollapseOrder CollapseOrderFn = func(rnd *rand.Rand, env NodeEnvironment) NodeID {
	if len(env.CollapsedMap) < len(env.Nodes) {
		idx := len(env.Nodes) - len(env.CollapsedMap) - 1
		node := env.Nodes[idx]
		if node != nil {
			return node.ID()
		}
	}
	return ""
}
