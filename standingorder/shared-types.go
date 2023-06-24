package standingorder

import (
	"os"
	"time"

  mt "webapp/moneytransfer"
)

var StandingOrdersTaskQueueName = os.Getenv("STANDING_ORDERS_TASK_QUEUE")

type PaymentSchedule struct {
	PeriodDuration time.Duration // seconds
	Active         bool
}

type StandingOrder struct {
	Schedule PaymentSchedule
	Details  mt.PaymentDetails
}

var ColorReset = "\033[0m"
var ColorRed = "\033[31m"
var ColorGreen = "\033[32m"
var ColorYellow = "\033[33m"
var ColorBlue = "\033[94m"
var ColorMagenta = "\033[35m"
var ColorCyan = "\033[36m"
var ColorWhite = "\033[37m"
