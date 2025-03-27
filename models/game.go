package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Table struct {
	gorm.Model
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	Members   []User         `json:"members"`
	Started   bool           `gorm:"not null" json:"started"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (table *Table) BeforeCreate(tx *gorm.DB) error {
	table.ID = uuid.New()
	return nil
}

func (table *Table) AddPlayer(player *User) error {
	for _, member := range table.Members {
		if member.ID.String() == player.ID.String() {
			return errors.New("player already in table")
		}
	}
	fmt.Printf("adding player: %+v\n", player)
	table.Members = append(table.Members, *player)
	return nil
}

type User struct {
	gorm.Model
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	TableID   uuid.UUID      `json:"tableId"`
	Table     Table          `json:"-"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}
