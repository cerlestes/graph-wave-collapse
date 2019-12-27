package gwc

import (
	"math/rand"
)

// CollapseOrderFn is a function that takes the current NodeEnvironment and returns the NodeID of the next Node to be collapsed.
type CollapseOrderFn func(*rand.Rand, NodeEnvironment) NodeID

// Collapses the Nodes in totally random order.
var RandomCollapseOrder CollapseOrderFn = func(rnd *rand.Rand, env NodeEnvironment) NodeID {
	for _, idx := range rnd.Perm(len(env.Nodes)) {
		id := env.GetID(idx)
		if _, collapsed := env.CollapsedMap[id]; collapsed == false {
			return id
		}
	}
	return ""
}

// Collapses the Nodes by choosing a random Node and then continuining with a random neighbour of the latest Node until running out of neighbours.
var RandomStreakCollapseOrder CollapseOrderFn = func(rnd *rand.Rand, env NodeEnvironment) NodeID {
	// If this is not the first run, search for uncollapsed neighbours of the current node.
	if env.Current != "" {
		neighbours := env.NodesMap[env.Current].Neighbours()
		for _, idx := range rnd.Perm(len(neighbours)) {
			id := neighbours[idx]
			if _, collapsed := env.CollapsedMap[id]; collapsed == false {
				return id
			}
		}
	}
	// If this is the first run, or the current node has run out of neighbours, pick a new random node to continue with.
	if len(env.Nodes) > 0 {
		return RandomCollapseOrder(rnd, env)
	}
	return ""
}

// Collapses the Nodes in ascending order.
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

// Collapses the Nodes in descending order.
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

// Produces a CollapseOrderFn that collapses the Nodes in the provided order.
func FixedCollapseOrder(order []NodeID) CollapseOrderFn {
	return func(rnd *rand.Rand, env NodeEnvironment) NodeID {
		if len(env.CollapsedMap) < len(env.Nodes) {
			return order[len(env.CollapsedMap)]
		}
		return ""
	}
}
