package database

import (
	"errors"
	"time"
	"tourist-api/utils"
)

type Location struct {
	ID            uint      `gorm:"primaryKey"`
	PlaceID       string    `gorm:"column:place_id"`
	LocationName  string    `gorm:"column:location_name"`
	Longitude     float64   `gorm:"column:longitude"`
	Latitude      float64   `gorm:"column:latitude"`
	Rating        float64   `gorm:"column:rating"`
	CoverImageUrl string    `gorm:"column:cover_image_url"`
	TotalReviews  int       `gorm:"column:total_reviews"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func CreateLocation(newLocation Location) (Location, error) {
	result := DB.Create(&newLocation)
	return newLocation, result.Error
}

func GetLocationByPlaceID(placeID string) (Location, error) {
	var location Location

	result := DB.Where("place_id = ?", placeID).Find(&location)
	return location, result.Error
}

func GetLocationByLocationID(locationID uint) (Location, error) {
	var location Location

	result := DB.Where("id = ?", locationID).Find(&location)
	return location, result.Error
}

func FindTripsLocationsRelation(tripID uint, locationID uint) TripsLocations {
	var tripsLocations TripsLocations
	DB.Where("trip_id = ? AND location_id = ?", tripID, locationID).Find(&tripsLocations)
	return tripsLocations
}

func BindLocationToTrip(userID uint, tripID uint, location Location) error {
	trip := GetUserTrip(userID, tripID)

	if trip.ID == 0 {
		return &utils.NotFoundError{
			Trace: errors.New("cannot find trip"),
		}
	}

	tripsLocations := FindTripsLocationsRelation(tripID, location.ID)
	if tripsLocations.ID != 0 {
		return &utils.ConflictError{
			Trace: errors.New("already exist"),
		}
	}

	err := DB.Model(&trip).Association("Locations").Append(&location)
	return err
}

func UnbindLocationFromTrip(userID uint, tripID uint, locationID uint) error {
	trip := GetUserTrip(userID, tripID)

	if trip.ID == 0 {
		return &utils.NotFoundError{
			Trace: errors.New("cannot find trip"),
		}
	}

	tripsLocations := FindTripsLocationsRelation(tripID, locationID)
	if tripsLocations.ID == 0 {
		return &utils.ConflictError{
			Trace: errors.New("not bind before"),
		}
	}

	DB.Delete(&tripsLocations)
	return nil
}
