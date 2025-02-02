package model

import "time"

type ScheduleID int64

type ScheduleSetting struct {
	FromTime        time.Time
	ToTime          time.Time
	Interval        time.Duration
	IntervalMinutes int
}

type Schedule struct {
	ID                   ScheduleID `json:"schedule_id"`
	UserID               UserID     `json:"user_id"`
	IntervalMinutes      int        `json:"interval_minutes"`
	LastNotificationTime time.Time  `json:"last_notification_time"`
}

func (s *Schedule) NextNotification() time.Time {
	return s.LastNotificationTime.Add(time.Duration(s.IntervalMinutes) * time.Minute)
}
