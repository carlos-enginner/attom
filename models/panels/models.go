package panels

type APIPayload struct {
	Context            interface{} `json:"event_dispatched"`
	NomePaciente       string      `json:"nomePaciente"`
	IdPainel           string      `json:"idPainel"`
	IdLocalAtendimento string      `json:"idLocalAtendimento"`
}

func (payload APIPayload) IsValid() bool {
	return payload.NomePaciente != "" && payload.IdPainel != "" && payload.IdLocalAtendimento != ""
}
