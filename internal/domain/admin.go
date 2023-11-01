package domain

import "context"

type AdminDashboard struct {
	Users     int64 `json:"users"`
	Items     int64 `json:"items"`
	Locations int64 `json:"locations"`
}
type AdminUsecase interface {
	Dashboard(ctx context.Context) (*AdminDashboard, error)
}
