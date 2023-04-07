package mrn

import (
    "container/heap"
    "math"
)

func edgeCost(p1 int, p2 int) float64 {
    if p1 >= 0 && p2>= 0 {
        return float64(1.0)
    }
    return math.MaxFloat64/4.0
}

// An Item is something we manage in a priority queue.
type PQ_Item struct {
	node_number int     // identity of the route node
	cost float64        // The cost of the path from the source to this node
	index int           // The index of the item in the heap.
}

// A PQ implements heap.Interface and holds PQ Items.
type PQ []*PQ_Item

func (pq PQ) Len() int { return len(pq) }

func (pq PQ) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].cost > pq[j].cost
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PQ) Push(x any) {
	n := len(*pq)
	item := x.(*PQ_Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PQ) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func RoutesFromHere(here int, Edges map[int][]int, destinations map[int]bool) map[int][]int {

    var h PQ
    settled := make(map[int]bool)

    thru := make(map[int]int)
    cost_from_src := make(map[int]float64)

    cost_from_src[here] = 0.0

    first_node := PQ_Item{ here, 0.0, 0 }
    heap.Push(&h, first_node)

    local_destinations := destinations
    number_destinations := len(destinations)

    for len(h) > 0 && number_destinations > 0 {
        // get the current node
        current_item := heap.Pop(&h).(PQ_Item)
        current_node := current_item.node_number
        current_cost := current_item.cost

        // might be we already settled, which means we can skip
        _, ok := settled[current_node]
        if ok {
            continue
        }

        settled[current_node] = true
        cost_from_src[current_node] = current_cost

        _, ok = local_destinations[current_node]

        // we need to find one fewer destination
        if ok {
            number_destinations -= 1
            delete(local_destinations, current_node)
        }

        // expand out to as yet unsettled peers
        for _, peer := range(Edges[current_node]) {
            _, ok = settled[peer]

            if ok {
                continue
            }

            trial_cost := current_cost + edgeCost(peer, current_node)

            peer_cost, ok := cost_from_src[peer] 
            if !ok {
                peer_cost = math.MaxFloat64
            }

            if trial_cost < peer_cost {
                cost_from_src[peer] = trial_cost
                thru[peer] = current_node
                peer_item := PQ_Item{ peer, trial_cost, 0}
                heap.Push(&h, peer_item)
            }
         }
    }

    // for each destination we can work back the link to get there
    // h.thru[destination] is id of node leading to destination,
    // h.thru[h.thru[destination]] is the next one, and so on

    // route_to[j] gives the sequence of nodes to visit
    // from the source to j
    //
    route_to := make(map[int][]int)

    for dest, _ := range(destinations) {

        reversed := make([]int, 0)
        tgt := dest

        for tgt != here {
            reversed = append(reversed,tgt)
            tgt = thru[tgt]
        }
        route_to[dest] = make([]int,len(reversed))
        last_reversed := len(reversed)-1
        for idx := last_reversed; idx>=0; idx-- { 
            fwd_idx := last_reversed-idx
            route_to[dest][fwd_idx] = reversed[idx]
        }

    }
    return route_to 
}

