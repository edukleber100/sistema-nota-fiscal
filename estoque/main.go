package main

import (
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
