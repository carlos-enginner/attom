package associations

import (
	"fmt"
	"src/post_relay/internal/utils"
	"src/post_relay/models/environment"
)

func LoadPainel(painel environment.Panels, cnes string, idCbo string) (*environment.Queue, error) {

	for _, painel := range painel.Items {

		if painel.Cnes != cnes {
			continue
		}

		if len(painel.Cbos) > 0 {

			cbo4Digit := utils.Substr(idCbo, 0, 4)

			if !utils.Contains(idCbo, painel.Cbos) && !utils.Contains(cbo4Digit, painel.Cbos) {
				continue
			}

		}

		return &painel.Queue, nil
	}

	return nil, fmt.Errorf("nenhum painel foi encontrado configurado para o CBO %s na unidade %s", idCbo, cnes)
}
