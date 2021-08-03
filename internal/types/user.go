package types

// UserCoordinates contains user location data.
type UserCoordinates struct {
	LocationID int    `db:"id"`
	Latitude   string `db:"latitude"`
	Longitude  string `db:"longitude"`
}
