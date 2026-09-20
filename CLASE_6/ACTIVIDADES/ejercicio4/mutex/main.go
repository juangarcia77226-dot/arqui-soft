package main

import (
	"fmt"
	"sync"
	"time"
)

// VARIANTE: el formulario indicará qué bloque copiar y pegar aquí.
const (
	initialStock = 184
	totalSales   = 52
	saleDelay    = 13 * time.Millisecond
)

// TODO: Ejercicio 4 (a) — Contador seguro con sync.Mutex
//
// Resolver la misma race condition de ejercicio4/inseguro protegiendo el
// contador con sync.Mutex:
//  1. Declarar un sync.Mutex junto al contador de stock.
//  2. Antes de decrementar, mu.Lock(); después, mu.Unlock().
//  3. El time.Sleep que simula la compra puede quedar FUERA del lock (no es
//     parte de la sección crítica).
//  4. Verificar con go run -race main.go que ya no aparece el reporte, y
//     que el resultado es initialStock menos totalSales.
func main() {
	stock := initialStock
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(totalSales)
	for i := 0; i < totalSales; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(saleDelay)

			mu.Lock()
			stock--
			mu.Unlock()
		}()
	}
	wg.Wait()

	fmt.Println("Stock final:", stock)
}
