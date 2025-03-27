package storage

import (
	"fmt"

	"github.com/7onn/croupier/models"
	"github.com/google/uuid"
)

func GetTableByID(id string) (*models.Table, error) {
	var table models.Table
	err := db.Preload("Members").Where("id = ?", id).First(&table).Error
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func GetFormingTable() (*models.Table, error) {
	var table *models.Table
	err := db.Where("started = ?", false).First(&table).Error
	if err == nil {
		return table, nil
	}

	table, err = CreateTable()
	if err != nil {
		return nil, err
	}
	return table, nil

}

func CreateTable() (*models.Table, error) {
	table := &models.Table{
		Members: []models.User{},
		Started: false,
	}
	err := db.Create(table).Error
	if err != nil {
		return nil, err
	}
	return table, nil
}

func UpdateTable(table *models.Table) error {
	err := db.Save(table).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(id string) *models.User {
	var user models.User
	db.Where("id = ?", id).First(&user)
	return &user
}

func CreateUser(name string, tableId string) (*models.User, error) {
	var user *models.User
	db.Where("name LIKE ? AND table_id = ?", name, tableId).First(&user)
	if user.ID != uuid.Nil {
		fmt.Printf("user found, not creating: %+v\n", user)
		return user, nil
	}
	user.Name = name
	user.TableID = uuid.MustParse(tableId)
	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
