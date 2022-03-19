package queue

type TaskHeap []Task

func (h TaskHeap) Len() int { return len(h) }

func (h TaskHeap) Less(i, j int) bool { return h[i].RunAt.Before(h[j].RunAt) }

func (h TaskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *TaskHeap) Push(x interface{}) {
	*h = append(*h, x.(Task))
}

// func (h *TaskHeap) Pop() interface{} {
// 	val := (*h)[h.Len()-1]
// 	*h = (*h)[0 : h.Len()-1]
// 	return val
// }

func (h *TaskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
