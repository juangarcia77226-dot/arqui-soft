package main

import (
	"fmt"
	"sync"
	"time"
)

// VARIANTE: el formulario indicará qué bloque copiar y pegar aquí.
const (
	firstStoreName  = "Tienda 12 A"
	firstStorePrice = 1876
	firstStoreDelay = 854 * time.Millisecond

	secondStoreName  = "Tienda 12 B"
	secondStorePrice = 1632
	secondStoreDelay = 406 * time.Millisecond
)

type Price struct {
	Store string
	Value int
}

func searchFirstStore() Price {
	time.Sleep(firstStoreDelay)
	return Price{Store: firstStoreName, Value: firstStorePrice}
}

func searchSecondStore() Price {
	time.Sleep(secondStoreDelay)
	return Price{Store: secondStoreName, Value: secondStorePrice}
}

func comparePrices() []Price {
	results := make([]Price, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		results[0] = searchFirstStore()
	}()

	go func() {
		defer wg.Done()
		results[1] = searchSecondStore()
	}()

	wg.Wait()
	return results
}

func main() {
	startedAt := time.Now()
	for _, price := range comparePrices() {
		fmt.Printf("%s: $%d\n", price.Store, price.Value)
	}
	fmt.Println("Time:", time.Since(startedAt))
}
