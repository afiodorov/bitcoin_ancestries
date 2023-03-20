package graph

import (
	"fmt"
)

type strings []string

func (s strings) Remove(el string) []string {
	res := make([]string, 0, len(s)-1)

	found := false

	for _, e := range s {
		if e == el && !found {
			found = true
			continue
		}

		res = append(res, e)
	}

	return res
}

type Graph struct {
	Outgoing map[string][]string
	Incoming map[string][]string
}

func MakeGraph() Graph {
	g := Graph{}
	g.Incoming = make(map[string][]string)
	g.Outgoing = make(map[string][]string)

	return g
}

func (g Graph) AddEdge(from, to string) {
	if from == "" {
		if _, ok := g.Incoming[to]; ok {
			// nothing to do
		} else {
			g.Incoming[to] = []string{}
		}

		return
	}

	if val, ok := g.Outgoing[from]; ok {
		g.Outgoing[from] = append(val, to)
	} else {
		g.Outgoing[from] = []string{to}
	}

	if val, ok := g.Incoming[to]; ok {
		g.Incoming[to] = append(val, from)
	} else {
		g.Incoming[to] = []string{from}
	}

	if _, ok := g.Incoming[from]; !ok {
		g.Incoming[from] = []string{}
	}
}

func (g Graph) IncomingEdges(to string) []string {
	if val, ok := g.Incoming[to]; ok {
		return val
	}

	return []string{}
}

func (g Graph) OutgoingEdges(from string) []string {
	if val, ok := g.Outgoing[from]; ok {
		return val
	}

	return []string{}
}

func (g Graph) RemoveEdge(from, to string) {
	if val, ok := g.Outgoing[from]; ok {
		g.Outgoing[from] = strings(val).Remove(to)
	} else {
		// nothing to do
	}

	if val, ok := g.Incoming[to]; ok {
		g.Incoming[to] = strings(val).Remove(from)
	} else {
		// nothing to do
	}
}

func (g Graph) HasEdges() bool {
	for _, edges := range g.Incoming {
		if len(edges) > 0 {
			return true
		}
	}

	return false
}

func (g Graph) TSort() ([]string, error) {
	res := make([]string, 0, len(g.Incoming))
	orphanNodes := make([]string, 0)

	for node, parents := range g.Incoming {
		if len(parents) == 0 {
			orphanNodes = append(orphanNodes, node)
		}
	}

	for {
		if len(orphanNodes) == 0 {
			break
		}

		node := orphanNodes[0]
		orphanNodes = orphanNodes[1:]

		res = append(res, node)

		for _, child := range g.OutgoingEdges(node) {
			g.RemoveEdge(node, child)

			if len(g.IncomingEdges(child)) == 0 {
				orphanNodes = append(orphanNodes, child)
			}
		}
	}

	if g.HasEdges() {
		return nil, fmt.Errorf("%v shouldn't have edges", g)
	}

	return res, nil
}
