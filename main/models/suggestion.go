package models

import "gorm.io/plugin/optimisticlock"

type Suggestion struct {
	Id      int
	OrderID int
	Time    string
	Price   string
	Comment string
	Email   string
	Phone   string
	Order   Order
	Version optimisticlock.Version
}
