package conf

import (
	"encoding/json"
	"github.com/google/wire"
	"os"
)

type Config struct {
	Database db `json:"database"`
}

type db struct {
	Dsn string `json:"dsn"`
}

var Provider = wire.NewSet(New)

func New() (*Config, error) {
	fp, err := os.Open("../config/config.json")
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var cfg *Config
	if err = json.NewDecoder(fp).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
