package evtq

import (
	"container/heap"
	"errors"
	_ "fmt"
	"github.com/isqrsqr/simulator/i2s2/vrtime"
)

// EventQueue represents the queue
type EventQueue struct {
	evt_id   int
	itemHeap *itemHeap
	lookup map[int]*item
}

// New initializes an empty priority queue.
func New() *EventQueue {
	return &EventQueue{
		evt_id:   0,
		itemHeap: &itemHeap{},
		lookup:   make(map[int]*item),
	}
}

// Len returns the number of elements in the queue.
func (p *EventQueue) Len() int {
	return p.itemHeap.Len()
}

func (p *EventQueue) MinTime() vrtime.Time {
	return (*p.itemHeap)[0].Priority
}

// Insert inserts a new element into the queue. No action is performed on duplicate elements.
func (p *EventQueue) Insert(v any, priority vrtime.Time) int {
	p.evt_id += 1
	newItem := &item{
		item_id:  p.evt_id,
		Value:    v,
		Priority: priority,
	}
	heap.Push(p.itemHeap, newItem)
	p.lookup[p.evt_id] = newItem
	return p.evt_id
}

// Pop removes the element with the highest priority from the queue and returns it.
// In case of an empty queue, an error is returned.
func (p *EventQueue) Pop() (any, vrtime.Time, int, error) {

	if len(*p.itemHeap) == 0 {
		return nil, vrtime.ZeroTime(), 0, errors.New("popping empty queue")
	}

	item := heap.Pop(p.itemHeap).(*item)
	delete(p.lookup, item.item_id)
	return item.Value, item.Priority, item.item_id, nil
}

// UpdatePriority changes the priority of a given item.
// If the specified item is not present in the queue, no action is performed.
func (p *EventQueue) UpdatePriority(evt_id int, newPriority vrtime.Time) {
	item, present := p.lookup[evt_id]
	if !present {
		return
	}
	item.Priority = newPriority
	heap.Fix(p.itemHeap, item.index)
}

// remove an element
func (p *EventQueue) Remove(evt_id int) bool {
	item, present := p.lookup[evt_id]
	if !present {
		return false
	}

	// we're going to push the element to be rid of to the top
	item.Priority = vrtime.ZeroTime()
	heap.Fix(p.itemHeap, item.index)

	// now pop it off
	p.Pop()
    return true
}

type itemHeap []*item

type item struct {
	item_id  int
	Value    any 
	Priority vrtime.Time
	index    int
}

func (ih *itemHeap) Len() int {
	return len(*ih)
}

func (ih *itemHeap) Less(i, j int) bool {
	return (*ih)[i].Priority.LT((*ih)[j].Priority)
	// invert for getting minimum at the top
	// return (*ih)[i].priority.GT(  (*ih)[j].priority )
}

func (ih *itemHeap) Swap(i, j int) {
	(*ih)[i], (*ih)[j] = (*ih)[j], (*ih)[i]
	(*ih)[i].index = i
	(*ih)[j].index = j
}

func (ih *itemHeap) Push(x any) {
	it := x.(*item)
	it.index = len(*ih)
	*ih = append(*ih, it)
}

func (ih *itemHeap) Pop() any {
	old := *ih
	item := old[len(old)-1]
	*ih = old[0 : len(old)-1]
	return item
}
