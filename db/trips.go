package database

import (
	"errors"
	"time"
)

type Trip struct {
	ID        uint       `gorm:"primaryKey"`
	TripName  string     `gorm:"column:trip_name"`
	UserID    uint       `gorm:"foreignKey:ID"`
	Locations []Location `gorm:"many2many:trips_locations"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
}

type TripsLocations struct {
	ID         uint      `gorm:"primaryKey"`
	TripID     uint      `gorm:"primaryKey"`
	LocationID uint      `gorm:"primaryKey"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

type TripDetailFlatten struct {
	ID            uint    `json:"ID"`
	PlaceID       string  `json:"PlaceID"`
	TripName      string  `json:"TripName"`
	LocationID    uint    `json:"LocationID"`
	LocationName  string  `json:"LocationName"`
	Longitude     float64 `json:"Longitude"`
	Latitude      float64 `json:"Latitude"`
	Rating        float32 `json:"Rating"`
	CoverImageURL string  `json:"CoverImageURL"`
	TotalReviews  int     `json:"TotalReviews"`
	TagID         uint    `json:"TagID"`
	TagName       string  `json:"TagName"`
	Color         string  `json:"Color"`
}

type TripDetail struct {
	ID        uint                 `json:"ID"`
	TripName  string               `json:"TripName"`
	Locations []TripDetailLocation `json:"Locations"`
}

type TripDetailLocation struct {
	ID            uint                    `json:"LocationID"`
	PlaceID       string                  `json:"PlaceID"`
	LocationName  string                  `json:"LocationName"`
	Longitude     float64                 `json:"Longitude"`
	Latitude      float64                 `json:"Latitude"`
	Rating        float32                 `json:"Rating"`
	CoverImageURL string                  `json:"CoverImageURL"`
	TotalReviews  int                     `json:"TotalReviews"`
	Tags          []TripDetailLocationTag `json:"Tags"`
}

type TripDetailLocationTag struct {
	ID      uint   `json:"TagID"`
	TagName string `json:"TagName"`
	Color   string `json:"Color"`
}

func CreateTrip(userID uint, tripName string) (Trip, error) {
	newTrip := Trip{UserID: userID, TripName: tripName}

	result := DB.Create(&newTrip)
	return newTrip, result.Error
}

func GetTrips(userID uint) ([]Trip, error) {
	var trips []Trip

	result := DB.Where("user_id = ?", userID).Find(&trips)
	return trips, result.Error
}

func GetUserTrip(userID uint, tripID uint) Trip {
	var trip Trip
	DB.Where("user_id = ? AND id = ?", userID, tripID).Find(&trip)
	return trip
}

// FIXME: 想要直接拿到不是 Flatten 的結果，但中間表的 locations_tags.color 找不到方式可以抓取到值
func GetTripDetail(userID uint, tripID uint) (TripDetail, error) {
	var tripDetailFlatten []TripDetailFlatten
	result := DB.Model(&Trip{}).
		Select(`
			trips.id,
			trips.trip_name,
			locations.id as location_id,
			locations.place_id,
			locations.location_name,
			locations.longitude,
			locations.latitude,
			locations.rating,
			locations.cover_image_url,
			locations.total_reviews,
			tags.id as tag_id,
			tags.tag_name,
			locations_tags.color
		`).
		Joins("LEFT JOIN trips_locations ON trips.id = trips_locations.trip_id").
		Joins("LEFT JOIN locations ON trips_locations.location_id = locations.id").
		Joins("LEFT JOIN locations_tags ON locations.id = locations_tags.location_id AND locations_tags.trip_id = trips.id").
		Joins("LEFT JOIN tags ON tags.id = locations_tags.tag_id").
		Where("trips.id = ? AND trips.user_id = ?", tripID, userID).
		Scan(&tripDetailFlatten)

	tripDetail, err := tripDetailFlattenToTripDetail(tripDetailFlatten)
	if err != nil {
		return tripDetail, err
	}

	return tripDetail, result.Error
}

func tripDetailFlattenToTripDetail(tripDetailFlatten []TripDetailFlatten) (TripDetail, error) {
	var tripDetail TripDetail

	if len(tripDetailFlatten) == 0 {
		return tripDetail, errors.New("input cannot be empty")
	}

	tripDetail.ID = tripDetailFlatten[0].ID
	tripDetail.TripName = tripDetailFlatten[0].TripName
	tripDetail.Locations = make([]TripDetailLocation, 0)

	locationMap := make(map[uint]*TripDetailLocation, 0)
	for _, row := range tripDetailFlatten {
		if _, ok := locationMap[row.LocationID]; !ok && row.LocationID != 0 {
			locationMap[row.LocationID] = &TripDetailLocation{
				ID:            row.LocationID,
				PlaceID:       row.PlaceID,
				LocationName:  row.LocationName,
				Longitude:     row.Longitude,
				Latitude:      row.Latitude,
				Rating:        row.Rating,
				CoverImageURL: row.CoverImageURL,
			}

			locationMap[row.LocationID].Tags = make([]TripDetailLocationTag, 0)
		}

		if row.TagID != 0 {
			locationMap[row.LocationID].Tags = append(locationMap[row.LocationID].Tags, TripDetailLocationTag{
				ID:      row.TagID,
				TagName: row.TagName,
				Color:   row.Color,
			})
		}
	}

	for _, location := range locationMap {
		tripDetail.Locations = append(tripDetail.Locations, *location)
	}

	return tripDetail, nil
}
