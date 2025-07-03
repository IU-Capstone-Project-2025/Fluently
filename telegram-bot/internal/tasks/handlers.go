package tasks

import (
	"context"
	"encoding/json"

	"fluently/telegram-bot/pkg/logger"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// TaskMux is a multiplexer for Asynq tasks
type TaskMux struct {
	mux *asynq.ServeMux
}

// NewTaskMux creates a new task multiplexer
func NewTaskMux() *TaskMux {
	return &TaskMux{
		mux: asynq.NewServeMux(),
	}
}

// HandleFunc registers a handler function for the given task type
func (tm *TaskMux) HandleFunc(taskType string, handler func(context.Context, *asynq.Task) error) {
	tm.mux.HandleFunc(taskType, handler)
}

// Handle registers a handler for the given task type
func (tm *TaskMux) Handle(taskType string, handler asynq.Handler) {
	tm.mux.Handle(taskType, handler)
}

// GetServeMux returns the underlying ServeMux
func (tm *TaskMux) GetServeMux() *asynq.ServeMux {
	return tm.mux
}

// TaskHandler interface for task processing
type TaskHandler interface {
	HandleLessonReminderTask(ctx context.Context, task *asynq.Task) error
	HandleDailyNotificationTask(ctx context.Context, task *asynq.Task) error
	HandleGenerateLessonTask(ctx context.Context, task *asynq.Task) error
	HandleSyncProgressTask(ctx context.Context, task *asynq.Task) error
	HandleCleanupSessionsTask(ctx context.Context, task *asynq.Task) error
}

// DefaultTaskHandler provides default implementations
type DefaultTaskHandler struct{}

// HandleLessonReminderTask handles lesson reminder tasks
func (h *DefaultTaskHandler) HandleLessonReminderTask(ctx context.Context, task *asynq.Task) error {
	var payload LessonReminderPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Log.Error("Failed to unmarshal lesson reminder payload", zap.Error(err))
		return err
	}

	logger.Log.Info("Processing lesson reminder task",
		zap.Int64("telegram_id", payload.TelegramID),
		zap.String("reminder_type", payload.ReminderType))

	// TODO: Implement lesson reminder logic
	// This would send a Telegram message to remind user about lessons

	return nil
}

// HandleDailyNotificationTask handles daily notification tasks
func (h *DefaultTaskHandler) HandleDailyNotificationTask(ctx context.Context, task *asynq.Task) error {
	var payload DailyNotificationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Log.Error("Failed to unmarshal daily notification payload", zap.Error(err))
		return err
	}

	logger.Log.Info("Processing daily notification task",
		zap.Int64("telegram_id", payload.TelegramID),
		zap.String("notification_type", payload.NotificationType))

	// TODO: Implement daily notification logic
	// This would send daily facts, motivation, or tips to users

	return nil
}

// HandleGenerateLessonTask handles lesson generation tasks
func (h *DefaultTaskHandler) HandleGenerateLessonTask(ctx context.Context, task *asynq.Task) error {
	var payload GenerateLessonPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Log.Error("Failed to unmarshal generate lesson payload", zap.Error(err))
		return err
	}

	logger.Log.Info("Processing generate lesson task",
		zap.Int64("telegram_id", payload.TelegramID),
		zap.String("user_id", payload.UserID),
		zap.String("cefr_level", payload.CEFRLevel))

	// TODO: Implement lesson generation logic
	// This would call the API to generate a new lesson and notify the user

	return nil
}

// HandleSyncProgressTask handles progress synchronization tasks
func (h *DefaultTaskHandler) HandleSyncProgressTask(ctx context.Context, task *asynq.Task) error {
	var payload SyncProgressPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Log.Error("Failed to unmarshal sync progress payload", zap.Error(err))
		return err
	}

	logger.Log.Info("Processing sync progress task",
		zap.Int64("telegram_id", payload.TelegramID),
		zap.String("user_id", payload.UserID))

	// TODO: Implement progress sync logic
	// This would sync progress data with the backend API

	return nil
}

// HandleCleanupSessionsTask handles session cleanup tasks
func (h *DefaultTaskHandler) HandleCleanupSessionsTask(ctx context.Context, task *asynq.Task) error {
	logger.Log.Info("Processing cleanup sessions task")

	// TODO: Implement session cleanup logic
	// This would clean up expired sessions and temporary data

	return nil
}
