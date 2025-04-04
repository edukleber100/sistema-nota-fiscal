package main

import (
	"encoding/json"
	"log"
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

type RequisicaoEstoque struct {
	Produtos []int `json:"produtos"`
}

func validarEstoque(w http.ResponseWriter, r *http.Request) {
	var req RequisicaoEstoque

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	for _, id := range req.Produtos {
		if produto, ok := produtos[id]; !ok || produto.Saldo <= 0 {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"ok": false}`))
			return
		}
	}

	w.Write([]byte(`{"ok": true}`))
}

func baixarEstoque(w http.ResponseWriter, r *http.Request) {
	var req RequisicaoEstoque

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	for _, id := range req.Produtos {
		if produto, ok := produtos[id]; !ok || produto.Saldo <= 0 {
			http.Error(w, "Erro ao baixar estoque", http.StatusConflict)
			return
		}
	}

	for _, id := range req.Produtos {
		produto := produtos[id]
		produto.Saldo--
		produtos[id] = produto
	}

	w.Write([]byte(`{"ok": true}`))
}

func main() {
	http.HandleFunc("/produtos/validar", validarEstoque)
	http.HandleFunc("/produtos/baixar", baixarEstoque)
	http.HandleFunc("/produtos", cadastrarProduto)
	http.HandleFunc("/produtos/listar", listarProdutos)

	log.Println("O serviço de estoque rodando na porta 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
