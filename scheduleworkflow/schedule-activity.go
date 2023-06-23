package scheduleworkflow

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

func ScheduleEmail(ctx context.Context, scheduleByID string, startTime time.Time, sd ScheduleDetails) error {

	activity.GetLogger(ctx).Info(ColorGreen, "ScheduleEmail:", ColorBlue, "Activity job run, scheduleID:", scheduleByID,
		"startTime:", startTime, "email:", sd.Email, ColorReset)

	sd.Id = scheduleByID
	err := SendEmailNotification(ctx, EmailNotificationStageRunning, sd)
	if err != nil {
		activity.GetLogger(ctx).Info("%sScheduleEmail:%s Failed to send email,%s %v", ColorGreen, ColorRed, ColorReset, err)
		return err
	}
	return nil
}
