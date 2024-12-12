package associations

import (
	"fmt"
	"src/post_relay/internal/logger"
	"src/post_relay/internal/utils"
	"src/post_relay/models/environment"
)

func LoadPainel(painel environment.Panels, cnes string, localChamada string) (*environment.Queue, error) {

	logger.GetLogger().Info("Associations.LoadPainel.start")

	for _, painel := range painel.Items {

		if painel.Cnes != cnes {
			continue
		}

		if !utils.Contains(localChamada, painel.Type) {
			continue
		}

		return &painel.Queue, nil
	}

	logger.GetLogger().Infof("Associations.LoadPainel.data - [cnes: %s local_chamada: %s])", cnes, localChamada)

	return nil, fmt.Errorf("nenhum painel de %s foi encontrado configurado na unidade %s. Por favor, verificar os campos: cnes e type", localChamada, cnes)
}
