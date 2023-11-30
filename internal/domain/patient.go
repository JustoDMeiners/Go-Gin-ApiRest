package domain

import "time"

// Producto es una estructura que define ...
type Patient struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
	Address  string    `json:"address"`
	Dni      string    `json:"dni"`
	DateUp   time.Time `json:"dateup"`
}
