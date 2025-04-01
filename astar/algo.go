package astar

import (
	"container/heap"
	"errors"
	"math"
)

func AStar(start, goal Node, graph Graph, h HeuristicFunc) ([]Node, float64, error) {
	// gScore guarda o custo real do início até cada nó
	gScore := make(map[Node]float64)
	gScore[start] = 0 // Custo para ir do início ao início é zero

	// fScore guarda o custo estimado total (gScore + heurística)
	fScore := make(map[Node]float64)
	fScore[start] = h(start, goal) // Estimativa inicial

	// cameFrom guarda de onde veio cada nó (para reconstruir o caminho)
	cameFrom := make(map[Node]Node)

	// openSetMap é um mapa que acompanha nossa fila de prioridade
	openSetMap := make(map[Node]*PqItem)
	pq := make(PriorityQueue, 0) // Nossa fila de prioridade

	// Começamos pelo nó inicial
	startItem := &PqItem{
		node:     start,
		priority: fScore[start],
	}
	heap.Push(&pq, startItem)     // Adiciona na fila
	openSetMap[start] = startItem // Marca como na fila

	// Enquanto houver nós para explorar
	for pq.Len() > 0 {
		// Pega o nó com menor fScore da fila
		currentItem := heap.Pop(&pq).(*PqItem)
		current := currentItem.node
		delete(openSetMap, current) // Remove do mapa de nós abertos

		// Se chegamos no objetivo, reconstruímos o caminho
		if current == goal {
			path := reconstructPath(cameFrom, current)
			return path, gScore[current], nil
		}

		// Para cada vizinho do nó atual
		for _, neighbor := range graph.Neighbors(current) {
			// Calcula o custo para ir até esse vizinho
			cost := graph.Cost(current, neighbor)
			tentativeGScore := gScore[current] + cost

			// Pega o gScore atual do vizinho (ou infinito se não conhecido)
			neighborGScore, ok := gScore[neighbor]
			if !ok {
				neighborGScore = math.Inf(1)
			}

			// Se encontramos um caminho melhor para esse vizinho
			if tentativeGScore < neighborGScore {
				// Atualizamos de onde ele veio
				cameFrom[neighbor] = current
				// Atualizamos o custo real
				gScore[neighbor] = tentativeGScore
				// Atualizamos a estimativa total
				fScore[neighbor] = tentativeGScore + h(neighbor, goal)

				// Se o vizinho já está na fila, atualizamos sua prioridade
				neighborItem, existsInOpenSet := openSetMap[neighbor]
				if existsInOpenSet {
					pq.Update(neighborItem, fScore[neighbor])
				} else {
					// Se não está, adicionamos na fila
					newItem := &PqItem{
						node:     neighbor,
						priority: fScore[neighbor],
					}
					heap.Push(&pq, newItem)
					openSetMap[neighbor] = newItem
				}
			}
		}
	}

	// Se a fila acabou e não chegamos no objetivo, não há caminho
	return nil, 0, errors.New("no path found")
}

// reconstructPath reconstrói o caminho do início até o objetivo
func reconstructPath(cameFrom map[Node]Node, current Node) []Node {
	totalPath := []Node{current} // Começa pelo final
	for {
		prev, ok := cameFrom[current] // De onde veio esse nó?
		if !ok {
			break // Quando não tem mais, chegamos no início
		}
		current = prev // Volta um passo
		// Adiciona no começo do caminho
		totalPath = append([]Node{current}, totalPath...)
	}
	return totalPath
}
