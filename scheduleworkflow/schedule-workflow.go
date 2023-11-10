package scheduleworkflow

import (
	"time"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"

  u "webapp/utils"
)

func ScheduleWorkflow(ctx workflow.Context, sd ScheduleDetails) error {

	logger := workflow.GetLogger(ctx)
	logger.Info(u.ColorGreen, "ScheduleWorkflow:", u.ColorReset, "Started - StartTime:", workflow.Now(ctx), "-", workflow.GetInfo(ctx).WorkflowExecution.ID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx2 := workflow.WithActivityOptions(ctx, ao)
	workflowInfo := workflow.GetInfo(ctx2)

	// Workflow Executions started by a Schedule have the following additional properties appended to
	// their search attributes
	// Doc ref: https://docs.temporal.io/workflows#action
	scheduledByIDPayload := workflowInfo.SearchAttributes.IndexedFields["TemporalScheduledById"]
	var scheduledByID string
	err := converter.GetDefaultDataConverter().FromPayload(scheduledByIDPayload, &scheduledByID)
	if err != nil {
		return err
	}

	startTimePayload := workflowInfo.SearchAttributes.IndexedFields["TemporalScheduledStartTime"]
	var startTime time.Time
	err = converter.GetDefaultDataConverter().FromPayload(startTimePayload, &startTime)
	if err != nil {
		return err
	}

	// notification activities based on pseudo process stage: processing
	err = workflow.ExecuteActivity(ctx2, ScheduleEmail, scheduledByID, startTime, sd).Get(ctx, nil)

	if err != nil {
		logger.Error(u.ColorGreen, "ScheduleWorkflow:", u.ColorReset, "ExecuteActivity failed.", u.ColorReset, "Error", err)
		return err
	}

	logger.Info(u.ColorGreen, "ScheduleWorkflow:", u.ColorReset, "Complete", "-", workflow.GetInfo(ctx).WorkflowExecution.ID)

	return nil
}
