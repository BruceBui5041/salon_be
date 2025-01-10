package models

const BookingServiceVersionM2M = "m2mbooking_service_version"

// BookingServiceVersion represents the many-to-many relationship between Booking and ServiceVersion
type BookingServiceVersion struct {
	BookingID        uint32          `json:"booking_id" gorm:"primaryKey;column:booking_id"`
	ServiceVersionID uint32          `json:"service_version_id" gorm:"primaryKey;column:service_version_id"`
	Booking          *Booking        `json:"-" gorm:"foreignKey:BookingID;references:Id"`
	ServiceVersion   *ServiceVersion `json:"-" gorm:"foreignKey:ServiceVersionID;references:Id"`
}

func (BookingServiceVersion) TableName() string {
	return BookingServiceVersionM2M
}
