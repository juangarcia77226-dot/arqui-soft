package main

import (
	"fmt"
	"sync"
	"time"
)

// VARIANTE: el formulario indicará qué bloque copiar y pegar aquí.
const (
	workerCount  = 2
	processDelay = 282 * time.Millisecond
)

var productIDs = []int{220, 221, 222, 223, 224}

func updateProduct(workerID int, productID int) {
	fmt.Printf("Worker %d updated product %d\n", workerID, productID)
	time.Sleep(processDelay)
}

func processWithWorkers() {
	jobs := make(chan int)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for w := 1; w <= workerCount; w++ {
		go func(workerID int) {
			defer wg.Done()
			for productID := range jobs {
				updateProduct(workerID, productID)
			}
		}(w)
	}

	for _, productID := range productIDs {
		jobs <- productID
	}
	close(jobs)

	wg.Wait()
}

func main() {
	startedAt := time.Now()
	processWithWorkers()
	fmt.Println("Time:", time.Since(startedAt))
}
