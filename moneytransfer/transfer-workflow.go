package moneytransfer

import (
	"fmt"
	"time"
	"webapp/utils"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

/* Transfer Workflow */
func TransferWorkflow(ctx workflow.Context, input PaymentDetails, delay int) (string, error) {

	// Define workflow logger (avoid repeating messages on replay etc)
	logger := workflow.GetLogger(ctx)
	logger.Info(ColorGreen, "Transfer-Workflow:", ColorReset, "Started")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	activityretrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        0, // unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	activityoptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,         // Timeout options specify when to automatically timeout Activity functions.
		RetryPolicy:         activityretrypolicy, // Temporal retries failed Activities by default.
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, activityoptions)

	// Set search attribute status to PROCESSING
	_ = utils.UpcertSearchAttribute(ctx, "CustomStringField", "PROCESSING")

	/* Withdraw money Activity - (blocks until completion with .Get function) */
	var withdrawOutput string

	withdrawErr := workflow.ExecuteActivity(ctx, Withdraw, input).Get(ctx, &withdrawOutput)

	if withdrawErr != nil {
		// Set search attribute status to FAILED
		_ = utils.UpcertSearchAttribute(ctx, "CustomStringField", "FAILED")
		logger.Info(ColorGreen, "Transfer-Workflow:", ColorReset, "Complete.", ColorRed, "(Withdraw Failed)", ColorReset)
		return "", fmt.Errorf("Withdraw: failed to withdraw funds from: %v, %w", input.SourceAccount, withdrawErr)
	}

	// For demo - sleep between activities
	logger.Debug("Transfer-Workflow: (DEBUG) Sleeping between activity calls -")
	logger.Info(ColorGreen, "Transfer-Workflow:", ColorBlue, "workflow.Sleep duration", delay, "seconds", ColorReset)
	workflow.Sleep(ctx, time.Duration(delay)*time.Second)

	/* Deposit money Activity - (blocks until completion with .Get function) */
	var depositOutput string

	depositErr := workflow.ExecuteActivity(ctx, Deposit, input).Get(ctx, &depositOutput)

	if depositErr != nil {
		// The deposit failed; put money back in original account.

		// Set search attribute status to FAILED
		_ = utils.UpcertSearchAttribute(ctx, "CustomStringField", "FAILED")

		/* Refund money Activity */
		var result string
		refundErr := workflow.ExecuteActivity(ctx, Refund, input).Get(ctx, &result)

		if refundErr != nil {
			logger.Info(ColorGreen, "Transfer-Workflow:", ColorReset, "Complete.", ColorRed, "(Deposit & Refund Failed)", ColorReset)
			return "", fmt.Errorf("Refund: failed to Deposit funds to: %v, %w. Money could NOT be returned to %v: %w",
				input.TargetAccount, depositErr, input.SourceAccount, refundErr)
		}

		logger.Info(ColorGreen, "Transfer-Workflow:", ColorReset, "Complete.", ColorRed, "(Deposit Failed)", ColorReset)
		return "", fmt.Errorf("Deposit: failed to deposit funds to: %v, Funds returned to: %v, %w",
			input.TargetAccount, input.SourceAccount, depositErr)
	}

	// Tranfer complete.
	result := fmt.Sprintf("Transfer complete (transaction IDs: Withdraw: %s, Deposit: %s)", withdrawOutput, depositOutput)

	// Set search attribute status to COMPLETED
	_ = utils.UpcertSearchAttribute(ctx, "CustomStringField", "COMPLETED")

	logger.Info(ColorGreen, "Transfer-Workflow:", ColorReset, "Complete.")

	return result, nil
}
