package models

import "time"

type File struct {
	Login    string
	Filename string
	UpdateAt time.Time
}
