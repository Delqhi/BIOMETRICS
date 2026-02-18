package database

import (
	"biometrics/internal/config"
	"biometrics/pkg/utils"

	"gorm.io/gorm"
)

type Postgres struct {
	DB  *gorm.DB
	log *utils.Logger
}

func NewPostgres(cfg config.DatabaseConfig) (*Postgres, error) {
	log := utils.NewLogger("info", "development")

	return &Postgres{
		DB:  nil,
		log: log,
	}, nil
}

func (p *Postgres) Close() error {
	return nil
}

func (p *Postgres) Transaction(fn func(*gorm.DB) error) error {
	return nil
}

func (p *Postgres) Raw(sql string, values ...interface{}) *gorm.DB {
	return nil
}

func (p *Postgres) Exec(sql string, values ...interface{}) *gorm.DB {
	return nil
}
