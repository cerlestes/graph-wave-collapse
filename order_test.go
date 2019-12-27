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

func Test_RandomCollapseOrder(t *testing.T) {
	super := newAbcdNodeSuperposition()
	nodes := newDefaultTestNodes(super)
	rnd := rand.New(rand.NewSource(42))

	sim := New(rnd, RandomCollapseOrder, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeStates{"A", "C", "C", "A", "A", "B", "C"}, collapsed.States())
	assert.EqualValues(t, NodeIDs{"0", "1", "6", "4", "3", "5", "2"}, collapsed.Collapsed())
}

func Test_NeighbourhoodCollapseOrder(t *testing.T) {
	super := newAbcdNodeSuperposition()
	nodes := newDefaultTestNodes(super)
	rnd := rand.New(rand.NewSource(420))

	sim := New(rnd, NeighbourhoodCollapseOrder, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeStates{"B", "A", "A", "A", "C", "A", "C"}, collapsed.States())
	assert.EqualValues(t, NodeIDs{"2", "1", "0", "4", "5", "3", "6"}, collapsed.Collapsed())
}

func Test_AscendingCollapseOrder(t *testing.T) {
	nodes := newDefaultTestNodes(NilSuperposition)
	rnd := rand.New(rand.NewSource(1337))

	sim := New(rnd, AscendingCollapseOrder, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeIDs{"0", "1", "2", "3", "4", "5", "6"}, collapsed.Collapsed())
}

func Test_DescendingCollapseOrder(t *testing.T) {
	nodes := newDefaultTestNodes(NilSuperposition)
	rnd := rand.New(rand.NewSource(3141))

	sim := New(rnd, DescendingCollapseOrder, nodes)
	collapsed := sim.Collapse()

	assert.EqualValues(t, NodeIDs{"6", "5", "4", "3", "2", "1", "0"}, collapsed.Collapsed())
}

func Test_EmptyNodes(t *testing.T) {
	rnd := rand.New(rand.NewSource(1337))
	env := NewNodeEnvironment(Nodes{})

	orders := []CollapseOrderFn{
		RandomCollapseOrder,
		NeighbourhoodCollapseOrder,
		AscendingCollapseOrder,
		DescendingCollapseOrder,
	}
	for _, order := range orders {
		next := order(rnd, *env)
		assert.Empty(t, next)
	}
}

func Test_LinearNodes(t *testing.T) {
	nodes := newLinearNodes(NilSuperposition)

	orders := []CollapseOrderFn{
		RandomCollapseOrder,
		NeighbourhoodCollapseOrder,
		AscendingCollapseOrder,
		DescendingCollapseOrder,
	}
	expected := []NodeIDs{
		NodeIDs{"1", "0", "3", "2"},
		NodeIDs{"1", "2", "3", "0"},
		NodeIDs{"0", "1", "2", "3"},
		NodeIDs{"3", "2", "1", "0"},
	}
	for i, order := range orders {
		rnd := rand.New(rand.NewSource(1337))
		sim := New(rnd, order, nodes)
		collapsed := sim.Collapse()

		assert.EqualValues(t, expected[i], collapsed.Collapsed())
	}
}
