package types

type UserCoordinates struct {
	LocationID int     `db:"id"`
	Latitude   float64 `db:"latitude"`
	Longitude  float64 `db:"longitude"`
}
