package gwc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newDefaultTestNodesEnvironment(super ...NodeSuperpositionFn) *NodeEnvironment {
	// The graph has the following form:
	//     6
	//     |
	//     4 - 5
	//     |   |
	// 1 - 2 - 3
	//     |
	//     0

	nodes := newDefaultTestNodes(super...)
	ne := NewNodeEnvironment(nodes)

	return ne
}

func Test_Methods(t *testing.T) {
	ne := newDefaultTestNodesEnvironment()
	ne.CollapsedMap = NodeCollapsedMap{
		"0": 6,
		"1": 5,
		"2": 4,
		"3": 3,
		"4": 2,
		"5": 1,
		"6": 0,
	}

	assert.Equal(t, "0", ne.GetID(0))
	assert.Equal(t, "3", ne.GetID(3))
	assert.Equal(t, "", ne.GetID(7))

	assert.Equal(t, 0, ne.GetIndex("0"))
	assert.Equal(t, 3, ne.GetIndex("3"))
	assert.Equal(t, -1, ne.GetIndex("7"))

	assert.EqualValues(t, NodeStates{nil, nil, nil, nil, nil, nil, nil}, ne.States())

	assert.EqualValues(t, NodeIDs{"6", "5", "4", "3", "2", "1", "0"}, ne.Collapsed())

	nodes := Nodes{
		ne.Nodes[6],
		ne.Nodes[5],
		ne.Nodes[4],
		ne.Nodes[3],
		ne.Nodes[2],
		ne.Nodes[1],
		ne.Nodes[0],
	}
	assert.EqualValues(t, nodes, ne.CollapsedNodes())

}

func Test_IsNeighbour(t *testing.T) {
	ne := newDefaultTestNodesEnvironment()
	ne.Current = "5"

	is_list := []bool{
		ne.IsNeighbour("3"),
		ne.IsNeighbour("4"),
	}
	isnt_list := []bool{
		ne.IsNeighbour("0"),
		ne.IsNeighbour("1"),
		ne.IsNeighbour("2"),
		ne.IsNeighbour("6"),
	}

	for _, is := range is_list {
		assert.True(t, is)
	}
	for _, isnt := range isnt_list {
		assert.False(t, isnt)
	}
}

func Test_IsWithinRange(t *testing.T) {
	ne := newDefaultTestNodesEnvironment()
	ne.Current = "6"

	is_list := []bool{
		ne.IsWithinRange("4", 1),
		ne.IsWithinRange("4", 2),
		ne.IsWithinRange("5", 2),
		ne.IsWithinRange("3", 3),
		ne.IsWithinRange("2", 2),
		ne.IsWithinRange("1", 3),
		ne.IsWithinRange("0", 3),
	}
	isnt_list := []bool{
		ne.IsWithinRange("4", 0),
		ne.IsWithinRange("5", 1),
		ne.IsWithinRange("3", 2),
		ne.IsWithinRange("2", 1),
		ne.IsWithinRange("1", 2),
		ne.IsWithinRange("0", 2),
	}

	for _, is := range is_list {
		assert.True(t, is)
	}
	for _, isnt := range isnt_list {
		assert.False(t, isnt)
	}
}

func Test_NodesWithinRange(t *testing.T) {
	ne := newDefaultTestNodesEnvironment()
	ne.Current = "6"

	// Check from current NodeID.
	range1 := ne.NodesWithinRangeIncl(1)
	range2 := ne.NodesWithinRangeIncl(2)
	range3 := ne.NodesWithinRangeExcl(3)

	assert.EqualValues(t, NodeIDs{"6", "4"}, range1)
	assert.EqualValues(t, NodeIDs{"6", "4", "2", "5"}, range2)
	assert.EqualValues(t, NodeIDs{"4", "2", "5", "0", "1", "3"}, range3)

	// Check from NodeID 3.
	range1 = ne.NodesWithinRangeOfIncl("3", 1)
	range2 = ne.NodesWithinRangeOfIncl("3", 2)
	range3 = ne.NodesWithinRangeOfExcl("3", 3)

	assert.EqualValues(t, NodeIDs{"3", "2", "5"}, range1)
	assert.EqualValues(t, NodeIDs{"3", "2", "5", "0", "1", "4"}, range2)
	assert.EqualValues(t, NodeIDs{"2", "5", "0", "1", "4", "6"}, range3)

	// Check from NodeID 2.
	range1 = ne.NodesWithinRangeOfIncl("2", 1)
	range2 = ne.NodesWithinRangeOfIncl("2", 2)
	range3 = ne.NodesWithinRangeOfExcl("2", 3)

	assert.EqualValues(t, NodeIDs{"2", "0", "1", "3", "4"}, range1)
	assert.EqualValues(t, NodeIDs{"2", "0", "1", "3", "4", "5", "6"}, range2)
	assert.EqualValues(t, NodeIDs{"0", "1", "3", "4", "5", "6"}, range3)
}

func Test_FilterNodes(t *testing.T) {
	ne := newDefaultTestNodesEnvironment()

	filtered_indexes_slice := ne.FilterNodes([]NodeID{"1", "2", "3"})
	filtered_indexes := ne.FilterNodes(NodeIDs{"4", "5", "6"})
	filtered := ne.FilterNodes(func(i NodeID, _ NodeState) bool {
		return (i == "0" || i == "1" || i == "2")
	})

	assert.EqualValues(t, NodeIDs{"1", "2", "3"}, filtered_indexes_slice)
	assert.EqualValues(t, NodeIDs{"4", "5", "6"}, filtered_indexes)
	assert.EqualValues(t, NodeIDs{"0", "1", "2"}, filtered)

	assert.Panics(t, func() {
		ne.FilterNodes("invalid")
	})

	as := ne.FilterNodesAnd(filtered_indexes_slice, filtered)
	bs := ne.FilterNodesOr(filtered_indexes, filtered)
	cs := ne.FilterNodesAnd(filtered, filtered)
	ds := ne.FilterNodesOr(filtered, filtered)

	assert.EqualValues(t, NodeIDs{"1", "2"}, as)
	assert.EqualValues(t, NodeIDs{"4", "5", "6", "0", "1", "2"}, bs)
	assert.EqualValues(t, filtered, cs)
	assert.EqualValues(t, filtered, ds)
}

func Test_NewNodeEnvironment(t *testing.T) {
	ne_empty := NewNodeEnvironment(Nodes{})
	ne_filled := NewNodeEnvironment(Nodes{NewSuperpositionNode("id", nil)})

	assert.NotNil(t, ne_empty)
	assert.NotNil(t, ne_filled)
}
