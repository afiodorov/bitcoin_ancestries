package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTSort1(t *testing.T) {
	g := MakeGraph()
	g.AddEdge("a", "b")

	r, err := g.TSort()
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b"}, r)
}

func TestTSort2(t *testing.T) {
	g := MakeGraph()
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")

	r, err := g.TSort()
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, r)
}

func TestTSort3(t *testing.T) {
	g := MakeGraph()
	g.AddEdge("5", "11")
	g.AddEdge("7", "11")
	g.AddEdge("7", "8")
	g.AddEdge("3", "8")
	g.AddEdge("3", "10")
	g.AddEdge("11", "2")
	g.AddEdge("11", "9")
	g.AddEdge("11", "10")
	g.AddEdge("8", "9")

	_, err := g.TSort()
	assert.Nil(t, err)
}

func TestTSort4(t *testing.T) {
	g := MakeGraph()
	g.AddEdge("9e7bc9715975579b36d0a81383dc76fcb41c9db45c77557083da5f4f76556b16", "e29a9dc0b293f0d84775adf316750255f267bc94a852e27a31db1f01c7a10e8d")
	g.AddEdge("4c0135d31a903852d488d908637478dedbea80167e532d23cfc188d2e5b5e3db", "e29a9dc0b293f0d84775adf316750255f267bc94a852e27a31db1f01c7a10e8d")
	g.AddEdge("", "d0a7e011c51cf17cf70b9f74aca9519fc9784e3b043c2927c6f25143c2f756ca")
	g.AddEdge("9e7bc9715975579b36d0a81383dc76fcb41c9db45c77557083da5f4f76556b16", "4c0135d31a903852d488d908637478dedbea80167e532d23cfc188d2e5b5e3db")
	r, err := g.TSort()
	assert.Nil(t, err)
	assert.Equal(t, []string{}, r)
}
