package service

import (
	"htn-backend/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService interface {
	CreateAccount(account *model.Account) error
	GetAccountByEmail(email string) (model.Account, error)
	CreateChore(chore *model.Chore) error
}

type dbService struct {
	db *gorm.DB
}

func NewDBService(connUrl string) (DBService, error) {
	db, err := gorm.Open(postgres.Open(connUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(model.Account{}, model.Chore{}, model.AssignedChore{}, model.Household{})
	return &dbService{db}, nil
}

func (s *dbService) CreateAccount(account *model.Account) error {
	return s.db.Create(account).Error
}

func (s *dbService) GetAccountByEmail(email string) (model.Account, error) {
	var account model.Account
	err := s.db.Where("email = ?", email).First(&account).Error
	return account, err
}

func (s *dbService) CreateChore(chore *model.Chore) error {
	return s.db.Create(chore).Error
}

func (s *dbService) CreateAssignedChore(assignedChore *model.AssignedChore) error {
	return s.db.Create(assignedChore).Error
}
