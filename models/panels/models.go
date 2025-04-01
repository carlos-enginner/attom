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

type TypeItem struct {
	Codigo    int64
	Descricao string
}

type TypesPanels struct {
	items []TypeItem
}

func (tp *TypesPanels) GetAll() []TypeItem {
	return tp.items
}

func (tp *TypesPanels) SetItems(items []TypeItem) {
	tp.items = items
}
