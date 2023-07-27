package models

import (
	"fmt"
	"gorm.io/plugin/optimisticlock"
)

type Order struct {
	Id              int
	ResourceID      int
	OrderNumber     string
	PublicationDate string
	Address         string
	Status          string
	ProcessingTypes []string `gorm:"-"`
	NameDetail      string
	Description     string
	Quantity        string
	Price           string
	ExpirationDate  string
	DownloadFileUrl string // TODO: fix to []string
	Resource        Resource
	Version         optimisticlock.Version
}

func (o Order) String() string {
	return fmt.Sprintf("Order\n"+
		"------------------\n"+
		"ID: %d\n"+
		"Order Number: %s\n"+
		"Publication Date: %s\n"+
		"Address: %s\n"+
		"Status: %s\n"+
		"Processing Types: %s\n"+
		"Name Detail: %s\n"+
		"Description: %s\n"+
		"Quantity: %s\n"+
		"Price: %s\n"+
		"Expiration Date: %s\n"+
		"Download Url: %s\n"+
		"Version: %d\n"+
		"------------------\n",
		o.Id,
		o.OrderNumber,
		o.PublicationDate,
		o.Address,
		o.Status,
		o.ProcessingTypes,
		o.NameDetail,
		o.Description,
		o.Quantity,
		o.Price,
		o.ExpirationDate,
		o.DownloadFileUrl,
		o.Version.Int64)
}
