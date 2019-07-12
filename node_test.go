package gwc

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newDefaultTestNodes(super NodeSuperposition) Nodes {
	// The graph has the following form:
	//     6
	//     |
	//     4 - 5
	//     |   |
	// 1 - 2 - 3
	//     |
	//     0

	nodes := Nodes{
		NewSuperpositionNode("0", super, "2"),
		NewSuperpositionNode("1", super, "2"),
		NewSuperpositionNode("2", super, "0", "1", "3", "4"),
		NewSuperpositionNode("3", super, "2", "5"),
		NewSuperpositionNode("4", super, "2", "5", "6"),
		NewSuperpositionNode("5", super, "4", "3"),
		NewSuperpositionNode("6", super, "4"),
	}

	return nodes
}

func newLinearNodes(super NodeSuperposition) Nodes {
	// The graph has the following form:
	// 0 - 1 - 2 - 3

	nodes := Nodes{
		NewSuperpositionNode("0", super, "1"),
		NewSuperpositionNode("1", super, "0", "2"),
		NewSuperpositionNode("2", super, "1", "3"),
		NewSuperpositionNode("3", super, "2"),
	}

	return nodes
}

func collapse(rnd *rand.Rand, ne *NodeEnvironment) NodeState {
	return ne.Nodes[0].Collapse(rnd, *ne)
}

func Test_Collapse(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))

	// This superposition is empty and thus can only yield a nil state.
	empty_super_ne := newDefaultTestNodesEnvironment(NodeSuperposition{})
	assert.Equal(t, nil, collapse(rnd, empty_super_ne))

	// This superposition's only function always yields a nil state.
	nil_super_ne := newDefaultTestNodesEnvironment(NilSuperposition)
	assert.Equal(t, nil, collapse(rnd, nil_super_ne))

	// This superposition only has a single non-probable state, resulting in it still collapsing into that state.
	non_nil_state := "non_nil_state"
	non_probable_super_ne := newDefaultTestNodesEnvironment(NodeSuperposition{
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, non_nil_state
		},
	})
	assert.Equal(t, non_nil_state, collapse(rnd, non_probable_super_ne))

	// This superposition only has non-probable states, resulting in collapsing into a random state.
	non_nil_state_2 := "non_nil_state_2"
	non_nil_state_3 := "non_nil_state_3"
	multi_non_probable_super_ne := newDefaultTestNodesEnvironment(NodeSuperposition{
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, non_nil_state
		},
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, non_nil_state_2
		},
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, non_nil_state_3
		},
	})
	assert.Equal(t, non_nil_state_2, collapse(rnd, multi_non_probable_super_ne))

	// This superposition has two non-probable states, but one of them yields nil, resulting in collapsing into the nil state.
	non_probable_nil_super_ne := newDefaultTestNodesEnvironment(NodeSuperposition{
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, non_nil_state
		},
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, nil
		},
	})
	assert.Equal(t, nil, collapse(rnd, non_probable_nil_super_ne))
}

func Test_NewNodes(t *testing.T) {
	non_nil_super := NodeSuperposition{
		func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
			return 0, "non-nil-value"
		},
	}
	new_node := NewSuperpositionNode("A", non_nil_super)
	new_nil_node := NewSuperpositionNodeWithNil("B", non_nil_super)

	assert.NotNil(t, new_node)
	assert.NotNil(t, new_nil_node)
}

func Test_BaseNode(t *testing.T) {
	rnd := rand.New(rand.NewSource(42))
	env := *NewNodeEnvironment(Nodes{})

	node := new(BaseNode)
	state := node.Collapse(rnd, env)

	assert.Nil(t, state)
}

func Test_And_Or_Xor(t *testing.T) {
	as := NodeIDs{"0", "1", "2", "3"}
	bs := NodeIDs{"2", "3", "4", "5"}

	ands := as.And(bs)
	ors := as.Or(bs)
	xors := as.Xor(bs)

	assert.EqualValues(t, NodeIDs{"2", "3"}, ands)
	assert.EqualValues(t, NodeIDs{"0", "1", "2", "3", "4", "5"}, ors)
	assert.EqualValues(t, NodeIDs{"0", "1", "4", "5"}, xors)
}
