package main

import "container/heap"
import "fmt"

type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work: // received a Request...
			b.dispatch(req) // ...so send it to a Worker
		case w := <-b.done: // a worker has finished ...
			// fmt.Printf("%d has %d tasks\n", w.index, w.pending)
			fmt.Println(b.pool)
			b.completed(w) // ...so update its info
		}
	}
}

// Send Request to worker
func (b *Balancer) dispatch(req Request) {
	// Grab the least loaded worker...
	w := heap.Pop(&b.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in its work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}

// Job is complete; update heap
func (b *Balancer) completed(w *Worker) {
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	heap.Remove(&b.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}
