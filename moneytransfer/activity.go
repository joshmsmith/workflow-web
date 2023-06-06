package moneytransfer

import (
  "context"
  "fmt"
  "log"
)

/* Withdraw Activity */
func Withdraw (ctx context.Context, data PaymentDetails) (string, error) {

  log.Printf("%sTransfer-Workflow-Withdraw-Activity:%s Withdrawing $%d from account %s.%s\n",
    ColorGreen, ColorBlue, data.Amount, data.SourceAccount, ColorReset)

  referenceID := fmt.Sprintf("%s-withdrawal", data.ReferenceID)
  bank := BankingService{"bank-api.example.com"}
  confirmation, err := bank.Withdraw(data.SourceAccount, data.Amount, referenceID)

  return confirmation, err
}

/* Deposit Activity */
func Deposit (ctx context.Context, data PaymentDetails) (string, error) {

  log.Printf("%sTransfer-Workflow-Deposit-Activity:%s Depositing $%d into account %s.%s\n",
    ColorGreen, ColorBlue, data.Amount, data.TargetAccount, ColorReset)

  referenceID := fmt.Sprintf("%s-deposit", data.ReferenceID)
  bank := BankingService{"bank-api.example.com"}
  confirmation, err := bank.Deposit(data.TargetAccount, data.Amount, referenceID)

  return confirmation, err
}

/* Refund Activity */
func Refund (ctx context.Context, data PaymentDetails) (string, error) {

  log.Printf("%sTransfer-Workflow-Refund-Activity:%s Refunding $%v back into account %v.%s\n\n",
    ColorGreen, ColorBlue, data.Amount, data.SourceAccount, ColorReset)

  referenceID := fmt.Sprintf("%s-refund", data.ReferenceID)
  bank := BankingService{"bank-api.example.com"}
  confirmation, err := bank.Deposit(data.SourceAccount, data.Amount, referenceID)

  return confirmation, err
}

