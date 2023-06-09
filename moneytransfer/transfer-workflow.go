package moneytransfer

import (
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

/* Transfer Workflow */
func Transfer(ctx workflow.Context, input PaymentDetails, delay int) (string, error) {

	log.Printf("%sTransfer-Workflow:%s Started", ColorGreen, ColorReset)

	// Define workflow logger (avoid repeating messages)
	logger := workflow.GetLogger(ctx)
	logger.Info("Workflow Logger: Transfer Workflow: Invoked -")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        0, // unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,

		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	// Set search attribute status to PROCESSING
	_ = UpcertSearchAttribute(ctx, "CustomStringField", "PROCESSING")

	/* Withdraw money Activity - (blocks until completion with .Get function) */
	var withdrawOutput string

	withdrawErr := workflow.ExecuteActivity(ctx, Withdraw, input).Get(ctx, &withdrawOutput)

	if withdrawErr != nil {
		// Set search attribute status to FAILED
		_ = UpcertSearchAttribute(ctx, "CustomStringField", "FAILED")
		log.Printf("%sTransfer-Workflow:%s Complete. %s(Withdraw Failed)%s", ColorGreen, ColorReset, ColorRed, ColorReset)
		return "", fmt.Errorf("Withdraw: failed to withdraw funds from: %v, %w", input.SourceAccount, withdrawErr)
	}

	// For demo - sleep between activities
	logger.Debug("Workflow Logger: Transfer Workflow: Sleeping between activity calls -")
	log.Printf("%sTransfer-Workflow:%s workflow.Sleep duration %d seconds%s", ColorGreen, ColorBlue, delay, ColorReset)
	workflow.Sleep(ctx, time.Duration(delay)*time.Second)

	/* Deposit money Activity - (blocks until completion with .Get function) */
	var depositOutput string

	depositErr := workflow.ExecuteActivity(ctx, Deposit, input).Get(ctx, &depositOutput)

	if depositErr != nil {
		// The deposit failed; put money back in original account.

		// Set search attribute status to FAILED
		_ = UpcertSearchAttribute(ctx, "CustomStringField", "FAILED")

		/* Refund money Activity */
		var result string
		refundErr := workflow.ExecuteActivity(ctx, Refund, input).Get(ctx, &result)

		if refundErr != nil {
			log.Printf("%sTransfer-Workflow:%s Complete. %s(Deposit & Refund Failed)%s", ColorGreen, ColorReset, ColorRed, ColorReset)
			return "", fmt.Errorf("Refund: failed to Deposit funds to: %v, %w. Money could NOT be returned to %v: %w",
				input.TargetAccount, depositErr, input.SourceAccount, refundErr)
		}

		log.Printf("%sTransfer-Workflow:%s Complete. %s(Deposit Failed)%s", ColorGreen, ColorReset, ColorRed, ColorReset)
		return "", fmt.Errorf("Deposit: failed to deposit funds to: %v, Funds returned to: %v, %w",
			input.TargetAccount, input.SourceAccount, depositErr)
	}

	// Tranfer complete.
	result := fmt.Sprintf("Transfer complete (transaction IDs: Withdraw: %s, Deposit: %s)", withdrawOutput, depositOutput)

	// Set search attribute status to COMPLETED
	_ = UpcertSearchAttribute(ctx, "CustomStringField", "COMPLETED")

	logger.Info("Workflow Logger: Transfer Workflow: Complete -")
	log.Printf("%sTransfer-Workflow:%s Complete.", ColorGreen, ColorReset)

	return result, nil
}
