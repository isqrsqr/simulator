package mrn

func edgeCost(p1 int, p2 int) float64 {
    if p1 >= 0 && p2>= 0 {
        return float64(1.0)
    }
    return math.MaxFloat64/4.0
}

type nodeHeap struct {
    nodes []int
    edges map[int][]int
    cost map[int]float64
    thru map[int]int
}
func (h *nodeHeap) initialize() {
    h.nodes = make([]int,0)
    h.edges = make(map[int][]int)
    h.cost  = make(map[int]float64)
    h.thru  = make(map[int]int)
}

func (h *nodeHeap) Len() int           { return len(h.nodes) }
func (h *nodeHeap) Less(i, j int) bool { return h.cost[j] < h.cost[i] }  // negative priority
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

func RoutesFromHere(here int, Edges map[int][]int, destinations map[int]bool) map[int]int {

    h := new(nodeHeap)
    h.initialize()
    h.edges = Edges

    settled := make(map[int]bool)

    // set up the analysis
    for node, _ := range(Edges) {
        h.cost[node] = math.MaxFloat64/2
    }

    h.cost[here] = 0.0
    heap.Push(h, here)

    local_destinations := destinations
    number_destinations := len(destinations)

    for len(h.nodes) > 0 and number_destinations > 0 {
        // get the current node
        current := h.Pop().(int)

        // might be we already settled, which means we can skip
        _, ok := settled[current]
        if ok {
            continue
        }
        settled[current] = true

        _, ok := local_destinations[current]
        if ok {
            num_destinations -= 1
            delete(local_destinations, current)
        }


        // put the unsettled
        for _, peer := range(Edges[current]) {
            _, ok := settled[peer]

            if ok {
                continue
            }

            trial_cost := h.cost[current] + EdgeCost(peer, current)
            if trial_cost < h.cost[peer] {
                h.cost[peer] = trial_cost
                h.thru[peer] = current
                heap.Push(h, peer)
            }
         }
    }

    // for each destination we can work back the link to get there
    // h.thru[destination] is id of node leading to destination,
    // h.thru[h.thru[destination]] is the next one, and so on
    //
    //  rtn[ path_step ][destination] gives the next hop from path_step heading towards destination
    //  need to be able to represent map[int]map[int]int


    // h.thru[x] gives the id of the node closer to the root on a traversal towards the root
    // rtn indexed by destination, returns id of next step to follow
    rtn := make(map[int]int)
    return rtn
}

