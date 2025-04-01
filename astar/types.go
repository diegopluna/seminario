package astar

import "container/heap"

type Node any

type Graph interface {
	Neighbors(node Node) []Node // Retorna os vizinhos
	Cost(from, to Node) float64 // Retorna custo de ir de um nó para o outro
}

// PqItem é um item na nossa fila de prioridade
type PqItem struct {
	node     Node    // O nó em si
	priority float64 // Prioridade (quanto menor, mais importante)
	index    int     // Posição na fila
}

// PriorityQueue é nossa fila de prioridade
type PriorityQueue []*PqItem

// Len diz quantos itens tem na fila
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compara dois itens - quem tem prioridade menor vem primeiro
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

// Swap troca dois itens de lugar na fila
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adiciona um novo item na fila
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*PqItem)
	item.index = n
	*pq = append(*pq, item)
}

// Pop remove o item com maior prioridade da fila
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// update atualiza a prioridade de um item na fila
func (pq *PriorityQueue) Update(item *PqItem, priority float64) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

// HeuristicFunc é uma função que estima o custo de um nó até o objetivo
type HeuristicFunc func(a, b Node) float64
