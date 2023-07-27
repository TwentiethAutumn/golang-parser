package parser_services

import "yellot-parser/main/models"

type OrderCard struct {
	Order       *models.Order
	Suggestions *[]models.Suggestion
}
