package dto

import "time"

type TransacoesRequest struct {
	Valor     int64  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacoesResponse struct {
	Limite int64 `json:"limite"`
	Saldo  int64 `json:"saldo"`
}

type History struct {
	Saldo             Saldo               `json:"saldo"`
	UltimasTransacoes []UltimasTransacoes `json:"ultimas_transacoes"`
}

type UltimasTransacoes struct {
	Valor       int64     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type Saldo struct {
	Total       int64     `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int64     `json:"limite"`
}
