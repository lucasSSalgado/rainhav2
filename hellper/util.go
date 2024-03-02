package hellper

import (
	"errors"
	"rinhaV2/dto"
)

func ValidarEntradaTransacoes(req dto.TransacoesRequest) error {
	if req.Tipo != "c" && req.Tipo != "d" {
		return errors.New("422")
	}

	if len(req.Descricao) > 10 || len(req.Descricao) <= 0 {
		return errors.New("422")
	}

	return nil
}
