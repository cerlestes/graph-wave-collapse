package gwc

import (
	"math"
	"math/rand"
)

// Builds a Node from the provided superposition and neighbours.
func NewSuperpositionNode(id NodeID, super NodeSuperposition, neighbours ...NodeID) Node {
	return &SuperpositionNode{BaseNode{id, neighbours}, super}
}

// Builds a Node from the provided superposition and neighbours, which is allowed to collapse into a nil result.
// This is required when no state was probable enough for the Node to collapse into.
func NewSuperpositionNodeWithNil(id NodeID, super NodeSuperposition, neighbours ...NodeID) Node {
	super = append(NodeSuperposition{}, super...)
	super = append(super, NilStateFn)
	return NewSuperpositionNode(id, super, neighbours...)
}

type (
	NodesMap map[NodeID]Node
	Nodes    []Node
	Node     interface {
		ID() NodeID
		Neighbours() NodeIDs
		Collapse(*rand.Rand, NodeEnvironment) NodeState
	}

	NodeIDs []NodeID
	NodeID  = string

	NodeProbability   = float64
	NodeState         = interface{}
	NodeStateFn       = func(*rand.Rand, NodeEnvironment) (NodeProbability, NodeState)
	NodeSuperposition = []NodeStateFn
)

var (
	// A nil state allows Nodes to collapse into a nil result, instead of forcing one of the states available.
	NilSuperposition NodeSuperposition = NodeSuperposition{NilStateFn}
	NilStateFn       NodeStateFn       = func(_ *rand.Rand, _ NodeEnvironment) (NodeProbability, NodeState) {
		return 0, nil
	}
)

// BaseNode can be used as base for a more concrete struct, which implements a concrete Collapse() method.
type BaseNode struct {
	id         NodeID
	neighbours NodeIDs
}

func (n *BaseNode) ID() NodeID {
	return n.id
}

func (n *BaseNode) Neighbours() NodeIDs {
	return n.neighbours
}

func (n *BaseNode) Collapse(rnd *rand.Rand, env NodeEnvironment) NodeState {
	return nil
}

// SuperpositionNode is a Node that, in order to determine state, collapses its superposition defined by probabilistic state producer functions.
type SuperpositionNode struct {
	BaseNode
	super NodeSuperposition
}

func (n *SuperpositionNode) Collapse(rnd *rand.Rand, env NodeEnvironment) NodeState {
	// Stop early when the Node's superposition is empty.
	num := len(n.super)
	if num == 0 {
		return nil
	}

	// We'll later compare the generated probabilities against this float, in this order.
	// Both are generated now, so that collapsing the superposition below won't interfere with these values.
	compare := rnd.Float64()
	order := rnd.Perm(num)

	// Call all functions in the superposition and collect their probabilities and states.
	sum := float64(0.0)
	probabilities := make([]NodeProbability, num)
	states := make([]NodeState, num)
	states_allow_nil := false
	for _, i := range order {
		ip, is := n.super[i](rnd, env)

		sum += ip
		probabilities[i] = ip
		states[i] = is

		if is == nil {
			states_allow_nil = true
		}
	}

	// Scale compare float according to the relative probability sum.
	compare *= math.Max(1, sum)

	// Collapse into the first state that had a high enough Nodeprobability to reach the compare float.
	for i, p := range probabilities {
		compare -= p
		if compare <= 0 {
			return states[i]
		}
	}

	// If no state was probable enough and nil state is not allowed, then return a random state.
	if !states_allow_nil {
		return states[rnd.Intn(len(states))]
	}
	return nil
}

// Applies a logical AND to the two index lists and returns the product.
func (nis NodeIDs) And(other NodeIDs) NodeIDs {
	xs := NodeIDs{}
	for _, a := range nis {
		for _, b := range other {
			if a == b {
				xs = append(xs, a)
			}
		}
	}
	return xs
}

// Applies a logical OR to the two index lists and returns the product.
func (nis NodeIDs) Or(other NodeIDs) NodeIDs {
	xs := NodeIDs{}
	xs = append(xs, nis...)
Outer:
	for _, b := range other {
		for _, x := range xs {
			if b == x {
				continue Outer
			}
		}
		xs = append(xs, b)
	}
	return xs
}

// Applies a logical XOR to the two index lists and returns the product.
func (nis NodeIDs) Xor(other NodeIDs) NodeIDs {
	xs := NodeIDs{}
Outer:
	for _, a := range nis {
		for _, b := range other {
			if a == b {
				continue Outer
			}
		}
		xs = append(xs, a)
	}
Outer2:
	for _, b := range other {
		for _, a := range nis {
			if a == b {
				continue Outer2
			}
		}
		xs = append(xs, b)
	}
	return xs
}
