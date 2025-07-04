package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// Task types
const (
	TaskSendLessonReminder    = "send:lesson_reminder"
	TaskSendDailyNotification = "send:daily_notification"
	TaskGenerateLesson        = "generate:lesson"
	TaskSyncProgress          = "sync:progress"
	TaskCleanupSessions       = "cleanup:sessions"
)

// Scheduler handles task scheduling using Asynq
type Scheduler struct {
	client *asynq.Client
	server *asynq.Server
}

// NewScheduler creates a new task scheduler
func NewScheduler(redisAddr, redisPassword string, redisDB int) *Scheduler {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}

	client := asynq.NewClient(redisOpt)

	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	return &Scheduler{
		client: client,
		server: server,
	}
}

// LessonReminderPayload represents lesson reminder task payload
type LessonReminderPayload struct {
	UserID         int64  `json:"user_id"`
	TelegramID     int64  `json:"telegram_id"`
	ReminderType   string `json:"reminder_type"` // "daily", "streak", "comeback"
	StreakDays     int    `json:"streak_days"`
	LastLessonDate string `json:"last_lesson_date"`
}

// DailyNotificationPayload represents daily notification task payload
type DailyNotificationPayload struct {
	UserID           int64  `json:"user_id"`
	TelegramID       int64  `json:"telegram_id"`
	NotificationType string `json:"notification_type"` // "fact", "motivation", "tip"
	CustomMessage    string `json:"custom_message"`
}

// GenerateLessonPayload represents lesson generation task payload
type GenerateLessonPayload struct {
	UserID         string `json:"user_id"`
	TelegramID     int64  `json:"telegram_id"`
	CEFRLevel      string `json:"cefr_level"`
	WordsPerLesson int    `json:"words_per_lesson"`
}

// SyncProgressPayload represents progress sync task payload
type SyncProgressPayload struct {
	UserID       string                 `json:"user_id"`
	TelegramID   int64                  `json:"telegram_id"`
	ProgressData map[string]interface{} `json:"progress_data"`
}

// ScheduleLessonReminder schedules a lesson reminder
func (s *Scheduler) ScheduleLessonReminder(userID, telegramID int64, reminderType string, delay time.Duration) error {
	payload := LessonReminderPayload{
		UserID:       userID,
		TelegramID:   telegramID,
		ReminderType: reminderType,
		StreakDays:   0,
	}

	task, err := NewLessonReminderTask(payload)
	if err != nil {
		return fmt.Errorf("failed to create lesson reminder task: %w", err)
	}

	info, err := s.client.Enqueue(task, asynq.ProcessIn(delay), asynq.Queue("default"))
	if err != nil {
		zap.L().With(
			zap.Error(err),
			zap.Int64("telegram_id", telegramID),
			zap.String("reminder_type", reminderType),
		).Error("Failed to schedule lesson reminder")
		return err
	}

	zap.L().With(
		zap.String("task_id", info.ID),
		zap.Int64("telegram_id", telegramID),
		zap.Duration("delay", delay),
	).Info("Scheduled lesson reminder")
	return nil
}

// ScheduleDailyNotification schedules a daily notification
func (s *Scheduler) ScheduleDailyNotification(userID, telegramID int64, notificationType string, scheduleTime time.Time) error {
	payload := DailyNotificationPayload{
		UserID:           userID,
		TelegramID:       telegramID,
		NotificationType: notificationType,
	}

	task, err := NewDailyNotificationTask(payload)
	if err != nil {
		return fmt.Errorf("failed to create daily notification task: %w", err)
	}

	info, err := s.client.Enqueue(task, asynq.ProcessAt(scheduleTime), asynq.Queue("default"))
	if err != nil {
		zap.L().With(
			zap.Error(err),
			zap.Int64("telegram_id", telegramID),
			zap.String("notification_type", notificationType),
		).Error("Failed to schedule daily notification")
		return err
	}

	zap.L().With(
		zap.String("task_id", info.ID),
		zap.Int64("telegram_id", telegramID),
		zap.Time("schedule_time", scheduleTime),
	).Info("Scheduled daily notification")
	return nil
}

// ScheduleGenerateLesson schedules lesson generation
func (s *Scheduler) ScheduleGenerateLesson(userID string, telegramID int64, cefrLevel string, wordsPerLesson int, delay time.Duration) error {
	payload := GenerateLessonPayload{
		UserID:         userID,
		TelegramID:     telegramID,
		CEFRLevel:      cefrLevel,
		WordsPerLesson: wordsPerLesson,
	}

	task, err := NewGenerateLessonTask(payload)
	if err != nil {
		return fmt.Errorf("failed to create generate lesson task: %w", err)
	}

	info, err := s.client.Enqueue(task, asynq.ProcessIn(delay), asynq.Queue("critical"))
	if err != nil {
		zap.L().With(
			zap.Error(err),
			zap.Int64("telegram_id", telegramID),
			zap.String("user_id", userID),
		).Error("Failed to schedule lesson generation")
		return err
	}

	zap.L().With(
		zap.String("task_id", info.ID),
		zap.Int64("telegram_id", telegramID),
		zap.String("user_id", userID),
	).Info("Scheduled lesson generation")
	return nil
}

// ScheduleProgressSync schedules progress synchronization
func (s *Scheduler) ScheduleProgressSync(userID string, telegramID int64, progressData map[string]interface{}) error {
	payload := SyncProgressPayload{
		UserID:       userID,
		TelegramID:   telegramID,
		ProgressData: progressData,
	}

	task, err := NewSyncProgressTask(payload)
	if err != nil {
		return fmt.Errorf("failed to create sync progress task: %w", err)
	}

	info, err := s.client.Enqueue(task, asynq.Queue("low"))
	if err != nil {
		zap.L().With(
			zap.Error(err),
			zap.Int64("telegram_id", telegramID),
			zap.String("user_id", userID),
		).Error("Failed to schedule progress sync")
		return err
	}

	zap.L().With(
		zap.String("task_id", info.ID),
		zap.Int64("telegram_id", telegramID),
		zap.String("user_id", userID),
	).Debug("Scheduled progress sync")
	return nil
}

// ScheduleRecurringDailyNotifications schedules recurring daily notifications
func (s *Scheduler) ScheduleRecurringDailyNotifications(userID, telegramID int64, notificationTime time.Time) error {
	// Schedule for the next 7 days
	for i := 0; i < 7; i++ {
		scheduleTime := notificationTime.Add(time.Duration(i) * 24 * time.Hour)
		err := s.ScheduleDailyNotification(userID, telegramID, "daily", scheduleTime)
		if err != nil {
			zap.L().With(
				zap.Error(err),
				zap.Int64("telegram_id", telegramID),
				zap.Int("day", i),
			).Error("Failed to schedule recurring daily notification")
			return err
		}
	}

	zap.L().With(
		zap.Int64("telegram_id", telegramID),
		zap.Time("start_time", notificationTime),
	).Info("Scheduled recurring daily notifications")
	return nil
}

// CancelUserTasks cancels all tasks for a specific user
func (s *Scheduler) CancelUserTasks(telegramID int64) error {
	// Note: Asynq doesn't provide a direct way to cancel tasks by custom criteria
	// This would require implementing a custom task tracker or using Redis directly
	zap.L().With(zap.Int64("telegram_id", telegramID)).Info("Cancelling user tasks")
	// Implementation would depend on specific requirements
	return nil
}

// Task creation functions
func NewLessonReminderTask(payload LessonReminderPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskSendLessonReminder, data), nil
}

func NewDailyNotificationTask(payload DailyNotificationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskSendDailyNotification, data), nil
}

func NewGenerateLessonTask(payload GenerateLessonPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskGenerateLesson, data), nil
}

func NewSyncProgressTask(payload SyncProgressPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskSyncProgress, data), nil
}

// Close closes the scheduler
func (s *Scheduler) Close() error {
	if err := s.client.Close(); err != nil {
		zap.L().Error("Failed to close Asynq client", zap.Error(err))
		return err
	}
	return nil
}

// GetClient returns the Asynq client for manual task enqueueing
func (s *Scheduler) GetClient() *asynq.Client {
	return s.client
}

// GetServer returns the Asynq server for running workers
func (s *Scheduler) GetServer() *asynq.Server {
	return s.server
}
