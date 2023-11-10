package scheduleworkflow

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"

  u "webapp/utils"
)

func ScheduleEmail(ctx context.Context, scheduleByID string, startTime time.Time, sd ScheduleDetails) error {

	activity.GetLogger(ctx).Info(u.ColorGreen, "ScheduleEmail:", u.ColorBlue, "Activity job run, scheduleID:", scheduleByID,
		"startTime:", startTime, "email:", sd.Email, u.ColorReset)

	sd.Id = scheduleByID
	err := SendEmailNotification(ctx, EmailNotificationStageRunning, sd)
	if err != nil {
		activity.GetLogger(ctx).Info("%sScheduleEmail:%s Failed to send email,%s %v", u.ColorGreen, u.ColorRed, u.ColorReset, err)
		return err
	}
	return nil
}
