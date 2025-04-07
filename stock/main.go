package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

var (
	products   = map[int]Product{}
	productMux sync.Mutex
)

func registerProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	productMux.Lock()
	p.ID = len(products) + 1
	products[p.ID] = p
	productMux.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productMux.Lock()
	defer productMux.Unlock()

	json.NewEncoder(w).Encode(products)
}

type StockRequest struct {
	Products []int `json:"products"`
}

func validateStock(w http.ResponseWriter, r *http.Request) {
	var req StockRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	for _, id := range req.Products {
		if product, ok := products[id]; !ok || product.Stock <= 0 {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"ok": false}`))
			return
		}
	}

	w.Write([]byte(`{"ok": true}`))
}

func decreaseStock(w http.ResponseWriter, r *http.Request) {
	var req StockRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	for _, id := range req.Products {
		if product, ok := products[id]; !ok || product.Stock <= 0 {
			http.Error(w, "Error decreasing stock", http.StatusConflict)
			return
		}
	}

	for _, id := range req.Products {
		product := products[id]
		product.Stock--
		products[id] = product
	}

	w.Write([]byte(`{"ok": true}`))
}

func main() {
	http.HandleFunc("/products/validate", validateStock)
	http.HandleFunc("/products/decrease", decreaseStock)
	http.HandleFunc("/products", registerProduct)
	http.HandleFunc("/products/list", listProducts)

	log.Println("Stock service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
