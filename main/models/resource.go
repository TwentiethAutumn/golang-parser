package models

import "gorm.io/plugin/optimisticlock"

type Resource struct {
	Id          int
	Name        string
	Description string
	Address     string
	Version     optimisticlock.Version
}
