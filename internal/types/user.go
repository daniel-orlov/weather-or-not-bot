package types

type UserCoordinates struct {
	UserID    int     `db:"user_id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}
