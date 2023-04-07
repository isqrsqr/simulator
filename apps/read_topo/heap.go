package main

import (
    "fmt"
    "container/heap"
)

type nodeHeap struct {
    nodes []int
    edges map[int][]int
    cost map[int]float64
}
func (h *nodeHeap) initialize() {
    h.nodes = make([]int,0)
    h.edges = make(map[int][]int)
    h.cost  = make(map[int]float64)
}

func (h *nodeHeap) Len() int           { return len(h.nodes) }
func (h *nodeHeap) Less(i, j int) bool { return h.cost[j] < h.cost[i] }
func (h *nodeHeap) Swap(i, j int)      { h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i] }

func (h *nodeHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	(*h).nodes = append((*h).nodes, x.(int))
}

func (h *nodeHeap) Pop() any {
	old := (*h).nodes
	n := len(old)
	x := old[n-1]
	(*h).nodes = old[0 : n-1]
	return x
}

func edgeCost(int i, int j) float64 {
    if i%2 + j%2 == 0 {
        return 0.5
    }

    if i%2 + j%2 == 2 {
        return 1.5
    }
    return 1.0
} 

func main() {
    h  := new(nodeHeap)
    h.initialize()
    h.nodes = append(h.nodes, 0)
    h.nodes = append(h.nodes, 1)
    h.nodes = append(h.nodes, 2)
    h.nodes = append(h.nodes, 3)
    h.nodes = append(h.nodes, 4)

    edges := make([]int,2)
    edges = append(edges,1,2)

    h.edges[0] = edges

    edges = nil
    edges = append(edges,0,2)
    h.edges[1] = edges

    edges = nil
    edges = append(edges,2,4)
    h.edges[3] = edges

    edges = nil
    edges = append(edges,3)
    h.edges[4] = edges

    heap.Init(h)
    heap.Push(h,0)

    cost := make(map[int]float64)
    cost[0] = 0.0

    

}
 
