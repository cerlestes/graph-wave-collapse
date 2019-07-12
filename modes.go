package gwc

import (
	"math/rand"
)

type CollapseModeFn func(*rand.Rand, NodeEnvironment) NodeID

func RandomCollapseMode(rnd *rand.Rand, env NodeEnvironment) NodeID {
FindRandomIndex:
	for _, idx := range rnd.Perm(len(env.Nodes)) {
		id := env.ToID(idx)
		for cid, _ := range env.CollapsedMap {
			if id == cid {
				continue FindRandomIndex
			}
		}
		return id
	}
	return ""
}

func NeighbourhoodCollapseMode(rnd *rand.Rand, env NodeEnvironment) NodeID {
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
		return RandomCollapseMode(rnd, env)
	}
	return ""
}

func AscendingCollapseMode(rnd *rand.Rand, env NodeEnvironment) NodeID {
	if len(env.CollapsedMap) < len(env.Nodes) {
		idx := len(env.CollapsedMap)
		node := env.Nodes[idx]
		if node != nil {
			return node.ID()
		}
	}
	return ""
}

func DescendingCollapseMode(rnd *rand.Rand, env NodeEnvironment) NodeID {
	if len(env.CollapsedMap) < len(env.Nodes) {
		idx := len(env.Nodes) - len(env.CollapsedMap) - 1
		node := env.Nodes[idx]
		if node != nil {
			return node.ID()
		}
	}
	return ""
}
