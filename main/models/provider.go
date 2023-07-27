package models

import (
	"fmt"
	"gorm.io/plugin/optimisticlock"
)

type Provider struct {
	Id                 int
	ResourceID         int
	CompanyName        string
	City               string
	Address            string
	CompanyDescription string
	TaxpayerId         string
	TypesOfServices    []string `gorm:"-"`
	Phone              string
	Email              string
	Resource           Resource
	Version            optimisticlock.Version
}

func (p Provider) String() string {
	return fmt.Sprintf("Provider\n"+
		"------------------\n"+
		"ID: %d\n"+
		"Company Name: %s\n"+
		"City: %s\n"+
		"Address: %s\n"+
		"Company Description: %s\n"+
		"Taxpayer Id: %s\n"+
		"Types of Services: %s\n"+
		"Phone: %s\n"+
		"Email: %s\n",
		p.Id,
		p.CompanyName,
		p.City,
		p.Address,
		p.CompanyDescription,
		p.TaxpayerId,
		p.TypesOfServices,
		p.Phone,
		p.Email)
}
