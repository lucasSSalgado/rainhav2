package controller

import (
	"net/http"
	"rinhaV2/database"
	"rinhaV2/dto"
	"rinhaV2/hellper"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
)

func InitRoutes(db *pgxpool.Pool) {
	mux := http.NewServeMux()
	c := cache.New(cache.NoExpiration, cache.NoExpiration)

	mux.HandleFunc("POST /clientes/{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			w.WriteHeader(422)
			return
		}

		if exists := database.CheckClient(db, id, c); !exists {
			w.WriteHeader(404)
			return
		}

		var req dto.TransacoesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(422)
			return
		}

		if err := hellper.ValidarEntradaTransacoes(req); err != nil {
			w.WriteHeader(422)
			return
		}

		if req.Tipo == "c" {
			limite, saldo, err := database.Creditar(db, id, req)
			if err != nil {
				w.WriteHeader(422)
				return
			}

			w.WriteHeader(200)
			json.NewEncoder(w).Encode(dto.TransacoesResponse{
				Limite: limite,
				Saldo:  saldo,
			})
			return
		}
		if req.Tipo == "d" {
			limite, saldo, err := database.Debitar(db, id, req)
			if err != nil {
				w.WriteHeader(422)
				return
			}

			w.WriteHeader(200)
			json.NewEncoder(w).Encode(dto.TransacoesResponse{
				Limite: limite,
				Saldo:  saldo,
			})
			return
		}
	})

	mux.HandleFunc("GET /clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			w.WriteHeader(422)
			return
		}

		if exists := database.CheckClient(db, id, c); !exists {
			w.WriteHeader(404)
			return
		}

		history, err := database.GetHistory(db, id)
		if err != nil {
			w.WriteHeader(422)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&history)
	})

	http.ListenAndServe(":8080", mux)
}
