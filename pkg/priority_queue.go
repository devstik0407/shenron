package pkg

import "container/heap"

type Item struct {
	value    string
	priority int
	index    int
}

func (it *Item) Value() string {
	return it.value
}

func (it *Item) Priority() int {
	return it.priority
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest priority, so we use less than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) PushItem(value string, priority int) {
	n := len(*pq)
	heap.Push(pq, &Item{
		value:    value,
		priority: priority,
		index:    n,
	})
}

func (pq PriorityQueue) TopItem() *Item {
	n := len(pq)
	if n > 0 {
		return pq[0]
	}
	return nil
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	if n == 0 {
		return nil
	}

	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) PopItem() *Item {
	return heap.Pop(pq).(*Item)
}

// Update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
