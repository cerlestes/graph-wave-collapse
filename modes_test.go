package gwc

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newAbcdNodeSuperposition() NodeSuperposition {
	return NodeSuperposition{
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 5, "A"
		},
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 2, "B"
		},
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 3, "C"
		},
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, "D"
		},
	}
}

func Test_RandomCollapseMode(t *testing.T) {
	super := newAbcdNodeSuperposition()
	nodes := newDefaultTestNodes(super)
	rnd := rand.New(rand.NewSource(42))

	sim := New(rnd, RandomCollapseMode, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeStates{"A", "C", "C", "A", "A", "B", "C"}, collapsed.States())
	assert.EqualValues(t, NodeIDs{"0", "1", "6", "4", "3", "5", "2"}, collapsed.Collapsed())
}

func Test_NeighbourhoodCollapseMode(t *testing.T) {
	super := newAbcdNodeSuperposition()
	nodes := newDefaultTestNodes(super)
	rnd := rand.New(rand.NewSource(420))

	sim := New(rnd, NeighbourhoodCollapseMode, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeStates{"B", "A", "A", "A", "C", "A", "C"}, collapsed.States())
	assert.EqualValues(t, NodeIDs{"2", "1", "0", "4", "5", "3", "6"}, collapsed.Collapsed())
}

func Test_AscendingCollapseMode(t *testing.T) {
	nodes := newDefaultTestNodes(NilSuperposition)
	rnd := rand.New(rand.NewSource(1337))

	sim := New(rnd, AscendingCollapseMode, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeIDs{"0", "1", "2", "3", "4", "5", "6"}, collapsed.Collapsed())
}

func Test_DescendingCollapseMode(t *testing.T) {
	nodes := newDefaultTestNodes(NilSuperposition)
	rnd := rand.New(rand.NewSource(3141))

	sim := New(rnd, DescendingCollapseMode, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeIDs{"6", "5", "4", "3", "2", "1", "0"}, collapsed.Collapsed())
}

func Test_EmptyNodes(t *testing.T) {
	rnd := rand.New(rand.NewSource(1337))
	env := NewNodeEnvironment(Nodes{})

	modes := []CollapseModeFn{
		RandomCollapseMode,
		NeighbourhoodCollapseMode,
		AscendingCollapseMode,
		DescendingCollapseMode,
	}
	for _, mode := range modes {
		next := mode(rnd, *env)
		assert.Empty(t, next)
	}
}

func Test_LinearNodes(t *testing.T) {
	nodes := newLinearNodes(NilSuperposition)

	modes := []CollapseModeFn{
		RandomCollapseMode,
		NeighbourhoodCollapseMode,
		AscendingCollapseMode,
		DescendingCollapseMode,
	}
	orders := []NodeIDs{
		NodeIDs{"1", "0", "3", "2"},
		NodeIDs{"1", "2", "3", "0"},
		NodeIDs{"0", "1", "2", "3"},
		NodeIDs{"3", "2", "1", "0"},
	}
	for i, mode := range modes {
		rnd := rand.New(rand.NewSource(1337))
		sim := New(rnd, mode, nodes)
		collapsed := sim.Collapse()

		assert.EqualValues(t, orders[i], collapsed.Collapsed())
	}
}
