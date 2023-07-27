package models

import "gorm.io/plugin/optimisticlock"

type ParseFile struct {
	Id          int
	DownloadUrl string
	FileName    string
	//Resource    Resource
	//Order       Order
	Version optimisticlock.Version
}
