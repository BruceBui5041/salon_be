package models

import (
	"salon_be/common"

	"gorm.io/gorm"
)

type Location struct {
	common.SQLModel `json:",inline"`
	UserType        string  `json:"user_type" gorm:"column:user_type;type:ENUM('service_man','customer');not null"`
	Latitude        float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,8);not null"`
	Longitude       float64 `json:"longitude" gorm:"column:longitude;type:decimal(11,8);not null"`
	Accuracy        float32 `json:"accuracy" gorm:"column:accuracy;type:float"`
}

func (Location) TableName() string {
	return "locations"
}

func (l *Location) Mask(isAdmin bool) {
	l.GenUID(common.DBTypeLocation)
}

func (l *Location) AfterFind(tx *gorm.DB) (err error) {
	l.Mask(false)
	return
}

// New BookingLocation model for initial locations
type BookingLocation struct {
	common.SQLModel      `json:",inline"`
	BookingID            uint32   `json:"booking_id" gorm:"column:booking_id;not null;uniqueIndex"`
	Booking              *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	CustomerLatitude     float64  `json:"customer_lat" gorm:"column:customer_lat;type:decimal(10,8);not null"`
	CustomerLongitude    float64  `json:"customer_long" gorm:"column:customer_long;type:decimal(11,8);not null"`
	DestinationLatitude  *float64 `json:"destination_latitude" gorm:"column:destination_latitude;type:decimal(10,8)"`
	DestinationLongitude *float64 `json:"destination_longitude" gorm:"column:destination_longitude;type:decimal(11,8)"`
	ServiceManLatitude   float64  `json:"service_man_lat" gorm:"column:service_man_lat;type:decimal(10,8);not null"`
	ServiceManLongitude  float64  `json:"service_man_long" gorm:"column:service_man_long;type:decimal(11,8);not null"`
	InitialDistance      float64  `json:"initial_distance" gorm:"column:initial_distance;type:decimal(10,2);not null"`
}

func (BookingLocation) TableName() string {
	return "booking_locations"
}

func (bl *BookingLocation) Mask(isAdmin bool) {
	bl.GenUID(common.DBTypeBookingLocation)
}

func (bl *BookingLocation) AfterFind(tx *gorm.DB) (err error) {
	bl.Mask(false)
	return
}

// New DistanceTracking model
type DistanceTracking struct {
	common.SQLModel      `json:",inline"`
	BookingID            uint32    `json:"booking_id" gorm:"column:booking_id;not null;index"`
	Booking              *Booking  `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	ServiceManLocationID uint32    `json:"service_man_location_id" gorm:"column:service_man_location_id;not null"`
	ServiceManLocation   *Location `json:"service_man_location,omitempty" gorm:"foreignKey:ServiceManLocationID"`
	UserLocationID       uint32    `json:"user_location_id" gorm:"column:user_location_id;not null"`
	UserLocation         *Location `json:"user_location,omitempty" gorm:"foreignKey:UserLocationID"`
	Distance             float64   `json:"distance" gorm:"column:distance;type:decimal(10,2);not null"`
}

func (DistanceTracking) TableName() string {
	return "distance_trackings"
}

func (dt *DistanceTracking) Mask(isAdmin bool) {
	dt.GenUID(common.DBTypeDistanceTracking)
}

func (dt *DistanceTracking) AfterFind(tx *gorm.DB) (err error) {
	dt.Mask(false)
	return
}
