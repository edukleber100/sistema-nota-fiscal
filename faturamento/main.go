package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func imprimirNotaFiscal(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/notas/imprimir/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for i, nota := range notasFiscais {
		if nota.ID == id {
			if nota.Status != "aberto" {
				http.Error(w, "Nota já foi processada", http.StatusBadRequest)
				return
			}

			if !validarEstoque(nota.Produtos) {
				http.Error(w, "Estoque insuficiente", http.StatusConflict)
				return
			}

			if !baixarEstoque(nota.Produtos) {
				http.Error(w, "Erro ao baixar estoque", http.StatusInternalServerError)
				return
			}

			notasFiscais[i].Status = "fechada"

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(notasFiscais[i])
			return
		}
	}

	http.Error(w, "Nota não encontrada", http.StatusNotFound)
}

type RequisicaoEstoque struct {
	Produtos []int `json:"produtos"`
}

func validarEstoque(produtos []int) bool {
	req := RequisicaoEstoque{Produtos: produtos}
	body, _ := json.Marshal(req)

	resp, err := http.Post("http://localhost:8081/produtos/validar", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func baixarEstoque(produtos []int) bool {
	req := RequisicaoEstoque{Produtos: produtos}
	body, _ := json.Marshal(req)

	resp, err := http.Post("http://localhost:8081/produtos/baixar", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func main() {
	http.HandleFunc("/notas", criarNotaFiscal)
	http.HandleFunc("/notas/listar", listarNotasFiscais)
	http.HandleFunc("/notas/imprimir/", imprimirNotaFiscal)

	fmt.Println("Serviço de Faturamento rodando na porta 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
