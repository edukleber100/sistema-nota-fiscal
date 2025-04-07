package main

import (
	"encoding/json"
	"net/http"
)

type NotaFiscal struct {
	ID       int    `json:"id"`
	Numero   string `json:"numero"`
	Status   string `json:"status"`
	Produtos []int  `json:"produtos"`
}

var notasFiscais []NotaFiscal
var nextID = 1

func criarNotaFiscal(w http.ResponseWriter, r *http.Request) {
	var nota NotaFiscal
	if err := json.NewDecoder(r.Body).Decode(&nota); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	nota.ID = nextID
	nota.Status = "aberto"
	notasFiscais = append(notasFiscais, nota)
	nextID++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nota)
}

func listarNotasFiscais(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notasFiscais)
}
