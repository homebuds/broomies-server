package service

import (
	"homebuds/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService interface {
	CreateAccount(account *model.Account) error
	GetAccountByEmail(email string) (model.Account, error)
	GetAccountByID(id uuid.UUID) (model.Account, error)
	GetAccountsWithLeastChores(householdID uuid.UUID) ([]model.Account, error)
	CreateAssignedChore(choreID uuid.UUID, accountId uuid.UUID, date time.Time) error
	CreateChore(chore *model.Chore) error
	GetChores(householdID uuid.UUID) ([]model.Chore, error)
	GetAssignedChores(householdID uuid.UUID) ([]model.AssignedChore, error)
	MarkAssignedChoreComplete(choreID uuid.UUID) error
	GetAssignedChore(assignedChoreID uuid.UUID) (model.AssignedChore, error)
	UpsertRoommateScore(accountID uuid.UUID, points uint) error
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

func (s *dbService) GetAccountByID(id uuid.UUID) (model.Account, error) {
	var account model.Account
	err := s.db.Where("id = ?", id).First(&account).Error
	return account, err
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

func (s *dbService) GetAccountsWithLeastChores(householdID uuid.UUID) ([]model.Account, error) {
	var accounts []model.Account
	err := s.db.Joins("INNER JOIN households ON households.id = accounts.household_id").
		Joins("LEFT JOIN assigned_chores ON accounts.id = assigned_chores.account_id").
		Where("household_id = ? AND (assigned_chores.id IS NULL OR assigned_chores.date > now() OR assigned_chores.completed = true)", householdID).
		Group("accounts.id").
		Order("count(assigned_chores.id) ASC").
		Find(&accounts).Error
	return accounts, err
}

func (s *dbService) CreateAssignedChore(choreID uuid.UUID, accountId uuid.UUID, date time.Time) error {
	assignedChore := model.AssignedChore{
		ChoreID:   choreID,
		AccountID: accountId,
		Date:      date,
		Completed: false,
	}
	return s.db.Create(&assignedChore).Error
}

func (s *dbService) GetChores(householdId uuid.UUID) ([]model.Chore, error) {
	var chores []model.Chore
	err := s.db.Where("household_id = ?", householdId).Find(&chores).Error
	return chores, err
}

func (s *dbService) GetAssignedChores(householdID uuid.UUID) ([]model.AssignedChore, error) {
	var assignedChores []model.AssignedChore
	err := s.db.Joins("Account", s.db.Where("household_id = ?", householdID)).
		Joins("Chore").
		Find(&assignedChores).Error
	return assignedChores, err
}

func (s *dbService) MarkAssignedChoreComplete(assignedChoreID uuid.UUID) error {
	return s.db.Model(&model.AssignedChore{}).Where("id = ?", assignedChoreID).Update("completed", true).Error
}

func (s *dbService) GetAssignedChore(assignedChoreID uuid.UUID) (model.AssignedChore, error) {
	var assignedChore model.AssignedChore
	err := s.db.Joins("Account").
		Joins("Chore").
		Where("id = ?", assignedChoreID).
		Find(&assignedChore).Error
	return assignedChore, err
}

func (s *dbService) UpsertRoommateScore(accountID uuid.UUID, points uint) error {
	var roommateScore model.RoommateScore
	err := s.db.Where("account_id = ?", accountID).First(&roommateScore).Error
	if err != nil {
		roommateScore = model.RoommateScore{
			AccountID: accountID,
			Points:    points,
		}
		return s.db.Create(&roommateScore).Error
	}
	return s.db.Model(&roommateScore).Update("points", roommateScore.Points+points).Error
}
