package database

import "time"

type Trip struct {
	ID        uint      `gorm:"primaryKey"`
	TripName  string    `gorm:"column:trip_name"`
	UserID    uint      `gorm:"foreignKey:ID"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func CreateTrip(userID uint, tripName string) (*Trip, error) {
	newTrip := &Trip{UserID: userID, TripName: tripName}

	result := DB.Create(&newTrip)
	return newTrip, result.Error
}

func GetTrips(userID uint) ([]Trip, error) {
	var trips []Trip

	result := DB.Where("user_id = ?", userID).Find(&trips)
	return trips, result.Error
}
