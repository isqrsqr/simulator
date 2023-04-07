package mrn

// An Item is something we manage in a priority queue.
type PQ_Item struct {
	node_number int     // identity of the route node
	cost float64        // The cost of the path from the source to this node
	index int           // The index of the item in the heap.
}

// A PQ implements heap.Interface and holds PQ Items.
type PQ []*PQ_Item

func (pq *PQ) Len() int { return len(pq) }

func (pq *PQ) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].cost > pq[j].cost
}

func (pq *PQ) Swap(i, j int) {
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

// Update modifies the priority and value of an Item in the queue.
func (pq *PQ) Update(item *PQ_Item, node_number int, cost float64) {
	item.node_number = node_number
	item.cost = cost
	heap.Fix(pq, item.index)
}

