package database

import (
	"context"
	"errors"
	"rinhaV2/dto"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPool() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig("host=postgres port=5432 user=postgres password=example dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	config.MaxConns = int32(9)
	config.MinConns = int32(9)

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}
	return db
}

func CheckClient(db *pgxpool.Pool, id int64) bool {
	var found int64
	err := db.QueryRow(context.Background(), "SELECT user_id FROM cliente WHERE user_id = $1", id).Scan(&found)
	return err == nil
}

func Creditar(db *pgxpool.Pool, id int64, req dto.TransacoesRequest) (int64, int64, error) {
	var limite int64
	var saldo int64
	tx, err := db.Begin(context.Background())
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), "SELECT limite, saldo FROM cliente WHERE user_id = $1 FOR UPDATE", id).Scan(&limite, &saldo)
	if err != nil {
		return 0, 0, err
	}

	newSaldo := saldo + req.Valor
	batch := &pgx.Batch{}

	batch.Queue("UPDATE cliente SET saldo = saldo + $1 WHERE user_id = $2", req.Valor, id)
	batch.Queue("INSERT INTO transacoes (user_id, valor, tipo, descricao, realizada_em) values ($1,$2,$3,$4,$5)",
		id, req.Valor, "c", req.Descricao, time.Now().UTC())

	s := tx.SendBatch(context.Background(), batch)
	if err := s.Close(); err != nil {
		return 0, 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, 0, err
	}

	return newSaldo, limite, nil
}

func Debitar(db *pgxpool.Pool, id int64, req dto.TransacoesRequest) (int64, int64, error) {
	var limite int64
	var saldo int64
	tx, err := db.Begin(context.Background())
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), "SELECT limite, saldo FROM cliente WHERE user_id = $1 FOR UPDATE", id).Scan(&limite, &saldo)
	if err != nil {
		return 0, 0, err
	}

	newSaldo := saldo - req.Valor
	if (newSaldo + limite) < 0 {
		return 0, 0, errors.New("412")
	}

	batch := &pgx.Batch{}
	batch.Queue("UPDATE cliente SET saldo = saldo - $1 WHERE user_id = $2", req.Valor, id)
	batch.Queue("INSERT INTO transacoes (user_id, valor, tipo, descricao, realizada_em) values ($1,$2,$3,$4,$5)",
		id, req.Valor, "d", req.Descricao, time.Now().UTC())

	s := tx.SendBatch(context.Background(), batch)
	if err := s.Close(); err != nil {
		return 0, 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, 0, err
	}

	return newSaldo, limite, nil
}

func GetHistory(db *pgxpool.Pool, id int64) (dto.History, error) {
	var limite int64
	var saldo int64

	err := db.QueryRow(
		context.Background(),
		"SELECT limite, saldo FROM cliente WHERE user_id = $1",
		id,
	).Scan(&limite, &saldo)
	if err != nil {
		return dto.History{}, err
	}

	rows, err := db.Query(
		context.Background(),
		"SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE user_id = $1 ORDER BY realizada_em DESC  LIMIT 10",
		id,
	)
	if err != nil {
		return dto.History{}, err
	}

	ultimas := []dto.UltimasTransacoes{}
	for rows.Next() {
		var valor int64
		var tipo string
		var descricao string
		var realizada_em time.Time

		if err := rows.Scan(&valor, &tipo, &descricao, &realizada_em); err != nil {
			return dto.History{}, err
		}

		ultimas = append(ultimas, dto.UltimasTransacoes{
			Valor:       valor,
			Tipo:        tipo,
			Descricao:   descricao,
			RealizadaEm: realizada_em,
		})
	}

	clientSaldo := dto.Saldo{
		Total:       saldo,
		DataExtrato: time.Now().UTC(),
		Limite:      limite,
	}

	resp := dto.History{
		Saldo:             clientSaldo,
		UltimasTransacoes: ultimas,
	}

	return resp, nil
}
