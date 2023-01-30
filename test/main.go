package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	var val1, val2 map[string]int

	go func() {
		defer wg.Done()
		resp, err := http.Get("https://api.com/endpoint1")
		if err != nil {
			fmt.Println("Erro ao pegar endpoint 1:", err)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&val1); err != nil {
			fmt.Println("Erro ao Unmarshal endpoint 1:", err)
			return
		}
		fmt.Println("Endpoint 1:", val1)
	}()

	go func() {
		defer wg.Done()
		resp, err := http.Get("https://api.com/endpoint2")
		if err != nil {
			fmt.Println("Erro ao pegar endpoint 2:", err)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&val2); err != nil {
			fmt.Println("Erro ao Unmarshal endpoint 2:", err)
			return
		}
		fmt.Println("Endpoint 2:", val2)
	}()

	wg.Wait()
	fmt.Println("Todos os endpoints foram pegos e Unmarshaled")
}
