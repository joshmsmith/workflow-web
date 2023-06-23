package moneytransfer

import (
	"os"
	"time"
)

var MoneyTransferTaskQueueName = os.Getenv("TRANSFER_MONEY_TASK_QUEUE")
var StandingOrdersTaskQueueName = os.Getenv("STANDING_ORDERS_TASK_QUEUE")
var DelayTimerBetweenWithdrawDeposit = os.Getenv("DELAY_TIMER_BETWEEN_WITHDRAW_DEPOSIT")

type PaymentDetails struct {
	SourceAccount string
	TargetAccount string
	ReferenceID   string
	Amount        int
}

type PaymentSchedule struct {
	PeriodDuration time.Duration // seconds
	Active         bool
}

type StandingOrder struct {
	Schedule PaymentSchedule
	Details  PaymentDetails
}

type WorkflowInfo struct {
	Id         int
	WorkflowID string
	RunID      string
	TaskQueue  string
	Info       string
	Status     string
}

var ColorReset = "\033[0m"
var ColorRed = "\033[31m"
var ColorGreen = "\033[32m"
var ColorYellow = "\033[33m"
var ColorBlue = "\033[94m"
var ColorMagenta = "\033[35m"
var ColorCyan = "\033[36m"
var ColorWhite = "\033[37m"
