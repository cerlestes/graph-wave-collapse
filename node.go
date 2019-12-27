package gwc

import (
	"math"
	"math/rand"
)

// Builds a Node from the provided state function and neighbours.
func NewNode(id NodeID, fn NodeStateFn, neighbours ...NodeID) Node {
	return &BaseNode{id, neighbours, fn}
}

// Builds a Node from the provided superposition and neighbours.
func NewSuperpositionNode(id NodeID, super NodeSuperposition, neighbours ...NodeID) Node {
	return NewNode(id, SuperpositionStateFn(super), neighbours...)
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

	NodeProbability = float64
	NodeState       = interface{}
	NodeStateFn     = func(*rand.Rand, NodeEnvironment) NodeState
)

// BaseNode can be used as base for a more concrete struct, which implements a concrete Collapse() method.
type BaseNode struct {
	id         NodeID
	neighbours NodeIDs
	fn         NodeStateFn
}

func (n *BaseNode) ID() NodeID {
	return n.id
}

func (n *BaseNode) Neighbours() NodeIDs {
	return n.neighbours
}

func (n *BaseNode) Collapse(rnd *rand.Rand, env NodeEnvironment) NodeState {
	if n.fn != nil {
		return n.fn(rnd, env)
	}
	return nil
}

// Applies a logical AND to the two index lists and returns the product.
func (ids NodeIDs) And(other NodeIDs) NodeIDs {
	xs := NodeIDs{}
	for _, a := range ids {
		for _, b := range other {
			if a == b {
				xs = append(xs, a)
			}
		}
	}
	return xs
}

// Applies a logical OR to the two index lists and returns the product.
func (ids NodeIDs) Or(other NodeIDs) NodeIDs {
	xs := NodeIDs{}
	xs = append(xs, ids...)
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
func (ids NodeIDs) Xor(other NodeIDs) NodeIDs {
	xs := NodeIDs{}
Outer:
	for _, a := range ids {
		for _, b := range other {
			if a == b {
				continue Outer
			}
		}
		xs = append(xs, a)
	}
Outer2:
	for _, b := range other {
		for _, a := range ids {
			if a == b {
				continue Outer2
			}
		}
		xs = append(xs, b)
	}
	return xs
}

type (
	NodeSuperpositionFn = func(*rand.Rand, NodeEnvironment) (NodeProbability, NodeState)
	NodeSuperposition   = []NodeSuperpositionFn
)

func SuperpositionStateFn(super NodeSuperposition) NodeStateFn {
	return func(rnd *rand.Rand, env NodeEnvironment) NodeState {
		// Stop early when the Node's superposition is empty.
		num := len(super)
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
		for _, i := range order {
			ip, is := super[i](rnd, env)

			sum += ip
			probabilities[i] = ip
			states[i] = is
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

		// If no state was probable enough but there are states available, return a random one.
		if len(states) > 0 {
			return states[rnd.Intn(len(states))]
		}

		// Fallback to nil when there were no states.
		return nil
	}
}
