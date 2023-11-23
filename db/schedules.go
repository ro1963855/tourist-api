package database

import "time"

type UserSchedule struct {
	ID           uint      `gorm:"primaryKey"`
	ScheduleName string    `gorm:"column:schedule_name"`
	UserID       uint      `gorm:"foreignKey:ID"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func CreateSchedule(userID uint, scheduleName string) (*UserSchedule, error) {
	newUserSchedule := &UserSchedule{UserID: userID, ScheduleName: scheduleName}

	result := DB.Create(&newUserSchedule)
	return newUserSchedule, result.Error
}
