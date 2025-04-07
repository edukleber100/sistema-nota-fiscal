package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Produto struct {
	ID    int     `json:"id"`
	Nome  string  `json:"nome"`
	Preco float64 `json:"preco"`
	Saldo int     `json:"saldo"`
}

var (
	produtos   = map[int]Produto{}
	produtoMux sync.Mutex
)

func cadastrarProduto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var p Produto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	produtoMux.Lock()
	p.ID = len(produtos) + 1
	produtos[p.ID] = p
	produtoMux.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func listarProdutos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	produtoMux.Lock()
	defer produtoMux.Unlock()

	json.NewEncoder(w).Encode(produtos)
}
