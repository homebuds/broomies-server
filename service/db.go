package service

import (
	"time"

	"github.com/google/uuid"
	model "github.com/homebuds/broomies-server/internal/dbmodel"
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
	GetAccountsWitMostPoints(householdID uuid.UUID, topN int) ([]model.RoommateScore, error)
	UpsertRoommateScore(accountID uuid.UUID, householdID uuid.UUID, points int) error
	CreateFinancialTransaction(amount float64, name string, accountID uuid.UUID, householdID uuid.UUID) (model.FinancialTransaction, error)
	GetFinancialTransactions(householdID uuid.UUID) ([]model.FinancialTransaction, error)
	GetUnsettledFinancialTransactions(householdID uuid.UUID) ([]model.FinancialTransaction, error)
	GetPointsForAccount(accountID uuid.UUID) (int, error)
	GetAccounts(householdId uuid.UUID) ([]model.Account, error)
	CreateNotification(accountID uuid.UUID, actorAccountID uuid.UUID, action string, choreID *uuid.UUID, financialTransactionID *uuid.UUID) (uuid.UUID, error)
	CreateUserNotification(accountID uuid.UUID, notificationID uuid.UUID) error
	GetUserNotifications(accountID uuid.UUID) ([]model.UserNotification, error)
	GetChore(choreId uuid.UUID) (model.Chore, error)
	SetNotificationSeen(notificationID uuid.UUID, accountID uuid.UUID) error
}

type dbService struct {
	db *gorm.DB
}

func NewDBService(connUrl string) (DBService, error) {
	db, err := gorm.Open(postgres.Open(connUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Migrator().AddColumn(&model.UserNotification{}, "seen")
	db.Migrator().DropColumn(&model.UserNotification{}, "reviewed")
	db.Migrator().AddColumn(&model.Notification{}, "created_at")
	db.Migrator().CreateTable(&model.Notification{}, &model.UserNotification{}, &model.Account{}, &model.Chore{}, &model.ChoreCompletionReview{})
	return &dbService{db}, nil
}

func (s *dbService) GetAccounts(householdId uuid.UUID) ([]model.Account, error) {
	var accounts []model.Account
	err := s.db.Where("household_id = ?", householdId).Find(&accounts).Error
	return accounts, err
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
		Where("assigned_chores.id = ?", assignedChoreID).
		Find(&assignedChore).Error
	return assignedChore, err
}

func (s *dbService) UpsertRoommateScore(accountID uuid.UUID, householdID uuid.UUID, points int) error {
	var roommateScore model.RoommateScore
	err := s.db.Where("account_id = ?", accountID).First(&roommateScore).Error
	if err == gorm.ErrRecordNotFound {
		roommateScore = model.RoommateScore{
			AccountID:   accountID,
			Points:      uint(500 + points),
			HouseholdID: householdID,
		}
		return s.db.Create(&roommateScore).Error
	} else if err != nil {
		return err
	}

	if points < 0 && roommateScore.Points < uint(-points) {
		points = 0
	} else {
		points = int(roommateScore.Points) + points
	}
	return s.db.Model(&roommateScore).Update("points", points).Error
}

func (s *dbService) GetAccountsWitMostPoints(householdID uuid.UUID, topN int) ([]model.RoommateScore, error) {
	var scores []model.RoommateScore
	err := s.db.Joins("Account").
		Joins("Household", s.db.Where("roommate_scores.household_id = ?", householdID)).
		Order("roommate_scores.points").
		Find(&scores).Limit(topN).Error
	return scores, err
}

func (s *dbService) CreateFinancialTransaction(
	amount float64, name string, accountID uuid.UUID, householdID uuid.UUID,
) (model.FinancialTransaction, error) {
	tx := model.FinancialTransaction{
		Amount:      amount,
		Name:        name,
		AccountID:   accountID,
		HouseholdID: householdID,
	}
	err := s.db.Create(&tx).Error
	return tx, err
}

func (s *dbService) GetFinancialTransactions(householdID uuid.UUID) ([]model.FinancialTransaction, error) {
	var txs []model.FinancialTransaction
	err := s.db.Joins("Account").
		Joins("Household", s.db.Where("financial_transactions.household_id = ?", householdID)).
		Find(&txs).Error
	return txs, err
}

func (s *dbService) GetUnsettledFinancialTransactions(householdID uuid.UUID) ([]model.FinancialTransaction, error) {
	var txs []model.FinancialTransaction
	err := s.db.
		Where("settled_at IS NULL AND household_id = ?", householdID).
		Find(&txs).Error
	return txs, err
}

func (s *dbService) GetPointsForAccount(accountID uuid.UUID) (int, error) {
	var roommateScore model.RoommateScore
	err := s.db.Where("account_id = ?", accountID).First(&roommateScore).Error
	return int(roommateScore.Points), err
}

func (s *dbService) CreateNotification(accountID uuid.UUID, actorID uuid.UUID, action string, choreID *uuid.UUID, financialTransactionID *uuid.UUID) (uuid.UUID, error) {
	notification := model.Notification{
		ActorAccountID:         actorID,
		Action:                 action,
		ActorChoreID:           choreID,
		FinancialTransactionID: financialTransactionID,
	}
	err := s.db.Create(&notification).Error
	return notification.ID, err
}

func (s *dbService) CreateUserNotification(accountID uuid.UUID, notificationID uuid.UUID) error {
	userNotification := model.UserNotification{
		AccountID:      accountID,
		NotificationID: notificationID,
		Seen:           false,
	}
	return s.db.Create(&userNotification).Error
}

func (s *dbService) GetUserNotifications(accountID uuid.UUID) ([]model.UserNotification, error) {
	var userNotifications []model.UserNotification
	err := s.db.Joins("Notification").
		Joins("Notification.ActorAccount").
		Joins("Notification.ActorChore").
		Where("user_notifications.account_id = ?", accountID).
		Find(&userNotifications).Error
	return userNotifications, err
}

func (s *dbService) GetChore(choreId uuid.UUID) (model.Chore, error) {
	var chore model.Chore
	err := s.db.Where("id = ?", choreId).First(&chore).Error
	return chore, err
}

func (s *dbService) SetNotificationSeen(notificationID uuid.UUID, accountID uuid.UUID) error {
	return s.db.Model(&model.UserNotification{}).Where("account_id = ? AND notification_id = ?", accountID, notificationID).Update("seen", true).Error
}
