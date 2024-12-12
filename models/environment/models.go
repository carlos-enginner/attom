package environment

import "time"

// Queue define a estrutura para a configuração da fila (queue)
type Queue struct {
	PanelUuid  string `toml:"panelUuid"`
	SectorUuid string `toml:"sectorUuid"`
}

// PanelItem define a estrutura de um item de painel
type PanelItem struct {
	Cnes        string   `toml:"cnes"`
	Description string   `toml:"description"`
	Type        []string `toml:"type"`
	Queue       Queue    `toml:"queue"`
}

// Panels define um conjunto de itens de painel
type Panels struct {
	Items []PanelItem `toml:"items"`
}

// API define as configurações de API
type API struct {
	Endpoint string `toml:"endpoint"`
	Token    string `toml:"token"`
	IBGE     string `toml:"ibge"`
}

// Database define as configurações de conexão com o banco de dados
type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbName"`
}

// Application define a configuração da aplicação
type Application struct {
	TimeoutConnection time.Duration `toml:"TimeoutConnectionInSeconds"`
	HttpDebug         bool          `toml:"HttpDebug"`
}

// Config é a estrutura principal que contém todas as outras configurações do ambiente
type Config struct {
	Application Application `toml:"application"`
	API         API         `toml:"api"`
	Database    Database    `toml:"database"`
	Panels      Panels      `toml:"panels"`
}
