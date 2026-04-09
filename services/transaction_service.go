package services

import (
	"errors"
	"time"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

type TransactionService struct {
	transactionRepository *repository.TransactionRepository
	accountRepository     *repository.AccountRepository
}

func NewTransactionService(
	transactionRepository *repository.TransactionRepository,
	accountRepository *repository.AccountRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
	}
}

func (transactionService *TransactionService) Create(
	accountID uint,
	operationType models.OperationType,
	amountInPaisa int64,
) (*models.Transaction, error) {

	if !operationType.IsValid() {
		return nil, errors.New("invalid operation_type_id")
	}

	account, err := transactionService.accountRepository.GetById(accountID)
	if err != nil {
		return nil, err
	}
	amount := absInt64(amountInPaisa)
	if !operationType.IsCredit() {
		amount = -amount
	}

	resultingBalance := account.BalanceInPaisa + amount
	if resultingBalance < 0 {
		return nil, errors.New("insufficient balance")
	}

	tx := &models.Transaction{
		AccountId:       accountID,
		OperationTypeId: operationType,
		AmountInPaisa:   amount,
		EventDate:       time.Now().UTC(),
	}

	if err := transactionService.transactionRepository.Create(tx); err != nil {
		return nil, err
	}
	if err := transactionService.accountRepository.UpdateBalanceByID(accountID, resultingBalance); err != nil {
		return nil, err
	}
	return tx, nil
}
