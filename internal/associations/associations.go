package associations

import (
	"fmt"
	"regexp"
	"src/post_relay/internal/logger"
	"src/post_relay/internal/utils"
	"src/post_relay/models/environment"
	"strings"
)

func LoadPainel(painel environment.Panels, cnes string, idCbo string, localChamada string) (*environment.Queue, error) {

	logger.GetLogger().Info("Associations.LoadPainel.start")

	for _, painel := range painel.Items {

		if painel.Cnes != cnes {
			continue
		}

		if strings.EqualFold(strings.TrimSpace(painel.Type), "triagem") {
			regexp, err := regexp.Compile(`(?i)\bESCUTA\b`)
			if err != nil {
				return nil, fmt.Errorf("erro na compilação da regex: %s", err)
			}

			if regexp.MatchString(localChamada) {
				return &painel.Queue, nil
			}
		}

		if strings.EqualFold(strings.TrimSpace(painel.Type), "atendimento") && localChamada == "ATENDIMENTO" && len(painel.Cbos) > 0 {

			cbo4Digit := utils.Substr(idCbo, 0, 4)

			if !utils.Contains(idCbo, painel.Cbos) && !utils.Contains(cbo4Digit, painel.Cbos) {
				continue
			}

			return &painel.Queue, nil
		}

	}

	logger.GetLogger().Infof("Associations.LoadPainel.data - [cnes: %s id_cbo: %s local_chamada: %s])", cnes, idCbo, localChamada)

	return nil, fmt.Errorf("nenhum painel de %s foi encontrado configurado na unidade %s. Por favor, verificar os campos: cnes, type e cbo", localChamada, cnes)
}
