package moneytransfer

// This code simulates a client for a hypothetical banking service.
// It supports both withdrawals and deposits, and generates a
// pseudorandom transaction ID for each request.
import (
	"errors"
	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

type account struct {
	AccountNumber string
	Balance       int64
}

type bank struct {
	Accounts []account
}

/* findAccount */
//func (b bank) findAccount(accountNumber string) (account, error) {
//
//  for _, v := range b.Accounts {
//    if v.AccountNumber == accountNumber {
//      return v, nil
//    }
//  }
//  return account{}, errors.New("account not found")
//}

/* findDbAccount */
func (b bank) findDbAccount(accountNumber string) (account, error) {

	dbaccs := ReadDbAccounts()

	for _, dbacc := range dbaccs {
		if dbacc.AccountName == accountNumber {
			return account{AccountNumber: dbacc.AccountName,
				Balance: int64(dbacc.AccountBalance)}, nil
		}
	}
	return account{}, errors.New("account not found in database")
}

/* InsufficientFundsError - raised when the account doesn't have enough money. */
type InsufficientFundsError struct{}

func (m *InsufficientFundsError) Error() string {
	return "Insufficient Funds"
}

/* InvalidAccountError - raised when the account name is invalid */
type InvalidAccountError struct{}

func (m *InvalidAccountError) Error() string {
	return "Account name supplied is invalid"
}

/* BankIntermittentError - raised but retryable */
type BankIntermittentError struct{}

func (m *BankIntermittentError) Error() string {
	return "Banking Service currently unavailable"
}

/* mockBank accounts
 *
 * Moved accounts to mysql database
 */
var mockBank = &bank{
	Accounts: []account{
		{AccountNumber: "85-150", Balance: 2000},
		{AccountNumber: "43-812", Balance: 0},
	},
}

/* BankingService mocks interaction with a bank API. It supports withdrawals and deposits */
type BankingService struct {
	// the hostname is to make it more realistic. This code does not
	// actually make any network calls.
	Hostname string
}

/* Withdraw - simulates a Withdrawal from a bank.
 * Accepts the account name (string), amount (int), and a reference ID (string)
 * for idempotent transaction tracking.
 * Returns a transaction id when successful
 * Returns various errors based on amount and account name.
 */
func (client BankingService) Withdraw(accountNumber string, amount int, referenceID string) (string, error) {

	// Check Bank Status
	if !checkBankService() {
		return "", &BankIntermittentError{}
	}

	// Check Account
	acct, err := mockBank.findDbAccount(accountNumber)

	if err != nil {
		return "", &InvalidAccountError{}
	}
	if amount > int(acct.Balance) {
		return "", &InsufficientFundsError{}
	}

	return generateTransactionID("W", 10), nil
}

/* Deposit - simulates a Withdrawal from a bank.
 * Acceptsthe account name (string), amount (int), and a reference ID (string)
 * for idempotent transaction tracking.
 * Returns a transaction id when successful
 * Returns InvalidAccountError if the account is invalid
 */
func (client BankingService) Deposit(accountNumber string, amount int, referenceID string) (string, error) {

	// Check Bank Status
	if !checkBankService() {
		return "", &BankIntermittentError{}
	}

	// Check Account
	_, err := mockBank.findDbAccount(accountNumber)
	if err != nil {
		return "", &InvalidAccountError{}
	}

	return generateTransactionID("D", 10), nil
}

/* generateTransactionID */
func generateTransactionID(prefix string, length int) string {
	randChars := make([]byte, length)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return prefix + string(randChars)
}
