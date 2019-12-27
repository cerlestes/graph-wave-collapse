package gwc

import (
	"math/rand"
)

func New(rnd *rand.Rand, mode CollapseOrderFn, nodes Nodes) *GraphWaveCollapse {
	return &GraphWaveCollapse{
		rnd:   rnd,
		mode:  mode,
		nodes: nodes,
	}
}

type GraphWaveCollapse struct {
	rnd   *rand.Rand
	mode  CollapseOrderFn
	nodes Nodes
}

func (gwc *GraphWaveCollapse) Collapse() NodeEnvironment {
	env := *NewNodeEnvironment(gwc.nodes)

	for {
		// Retrieve next NodeIndex according to mode.
		next := gwc.mode(gwc.rnd, env)
		if _, exists := env.NodesMap[next]; !exists {
			break
		}

		// Collapse the chosen Node and mark it as such.
		env.Current = next
		env.StatesMap[next] = env.NodesMap[next].Collapse(gwc.rnd, env)
		env.CollapsedMap[next] = len(env.CollapsedMap)
	}

	return env
}
