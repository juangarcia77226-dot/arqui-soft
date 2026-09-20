package main

import (
	"context"
	"fmt"
	"time"
)

// VARIANTE: el formulario indicará qué bloque copiar y pegar aquí.
const (
	productName    = "Producto 12"
	reviewsScore   = 42
	reviewsDelay   = 1860 * time.Millisecond
	timeoutLimit   = 370 * time.Millisecond
	defaultReviews = 0
)

func getReviews() int {
	time.Sleep(reviewsDelay)
	return reviewsScore
}

func getProductReviews() int {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutLimit)
	defer cancel()

	resultCh := make(chan int, 1)
	go func() {
		resultCh <- getReviews()
	}()

	select {
	case reviews := <-resultCh:
		return reviews
	case <-ctx.Done():
		return defaultReviews
	}
}

func main() {
	startedAt := time.Now()
	reviews := getProductReviews()
	fmt.Printf("%s reviews: %d\n", productName, reviews)
	fmt.Println("Time:", time.Since(startedAt))
}
