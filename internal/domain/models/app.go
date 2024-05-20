package models

type App struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	Secret string `db:"secret"`
}
