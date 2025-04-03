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

// models: tipos de painel

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

// models: unidades

type Unidade struct {
	NuCnes        string
	NomeUnidade   string
	NomeMunicipio string
}

type Unidades struct {
	items []Unidade
}

func (u *Unidades) GetAll() []Unidade {
	return u.items
}

func (u *Unidades) SetItems(items []Unidade) {
	u.items = items
}

// models: painels ativos

type LocalAtendimento struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type PainelActive struct {
	NuCnes           string             `json:"cnes"`
	Descricao        string             `json:"descricao"`
	IDPainel         string             `json:"idPainel"`
	NomePainel       string             `json:"nomePainel"`
	DuracaoChamada   int                `json:"duracaoChamada"`
	LocalAtendimento []LocalAtendimento `json:"localAtendimento"`
}

type ActivesResponse struct {
	Error bool           `json:"error"`
	Msg   string         `json:"msg"`
	Obj   []PainelActive `json:"obj"`
}

func (pa *PainelActive) RequestGetAll(cnes string) bool {
	return cnes == "0000000"
}

func (ar *ActivesResponse) SetItems(items []PainelActive) {
	ar.Obj = items
}

func (ar *ActivesResponse) GetPanelsActives() []PainelActive {
	return ar.Obj
}
