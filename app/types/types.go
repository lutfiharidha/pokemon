package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Pokedex struct {
	Name    string    `json:"name_pokedex"`
	Pokemon []Pokemon `json:"pokemon,omitempty"`
	Stock   int       `json:"stock"`
}

type Pokemon struct {
	Name  string  `json:"name_pokemon,omitempty"`
	Hp    int     `json:"hp,omitempty"`
	Moves []Moves `json:"moves,omitempty"`
	Point int64   `json:"point,omitempty"`
}

type Moves struct {
	Name  string
	Power int
}

type Logging struct {
	DataLog []string
	Winner  string
}

type Log struct {
	ID        string         `gorm:"uniqueIndex" json:"id"`
	DataLog   datatypes.JSON `json:"data_log" gorm:"not null"`
	Winner    string         `json:"winner"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

func (friend *Log) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.New()
	tx.Statement.SetColumn("ID", uuid)
	return
}
