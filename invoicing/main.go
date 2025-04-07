package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Invoice struct {
	ID       int    `json:"id"`
	Number   string `json:"number"`
	Status   string `json:"status"`
	Products []int  `json:"products"`
}

var invoices []Invoice
var nextID = 1

func createInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	invoice.ID = nextID
	invoice.Status = "open"
	invoices = append(invoices, invoice)
	nextID++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

func listInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func printInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/invoices/print/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, invoice := range invoices {
		if invoice.ID == id {
			if invoice.Status != "open" {
				http.Error(w, "Invoice already processed", http.StatusBadRequest)
				return
			}

			if !validateStock(invoice.Products) {
				http.Error(w, "Insufficient stock", http.StatusConflict)
				return
			}

			if !decreaseStock(invoice.Products) {
				http.Error(w, "Failed to decrease stock", http.StatusInternalServerError)
				return
			}

			invoices[i].Status = "closed"

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invoices[i])
			return
		}
	}

	http.Error(w, "Invoice not found", http.StatusNotFound)
}

type StockRequest struct {
	Products []int `json:"products"`
}

func validateStock(products []int) bool {
	req := StockRequest{Products: products}
	body, _ := json.Marshal(req)

	resp, err := http.Post("http://localhost:8081/products/validate", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func decreaseStock(products []int) bool {
	req := StockRequest{Products: products}
	body, _ := json.Marshal(req)

	resp, err := http.Post("http://localhost:8081/products/decrease", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func main() {
	http.HandleFunc("/invoices", createInvoice)
	http.HandleFunc("/invoices/list", listInvoices)
	http.HandleFunc("/invoices/print/", printInvoice)

	fmt.Println("Invoice Service running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
