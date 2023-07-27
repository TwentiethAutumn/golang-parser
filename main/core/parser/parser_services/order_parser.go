package parser_services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"sync"
	"yellot-parser/main/core/parser/parser_settings"
	"yellot-parser/main/models"
)

type OrderParser struct {
	Settings                parser_settings.ParserSettings
	ResourceCache           map[string]models.Resource
	SuggestionParser        *SuggestionParser
	Postfix                 string
	StartPointTag           *CustomTag
	EndPointTag             *CustomTag
	DetailOrderTag          *CustomTag
	OrderPubDateTag         *CustomTag
	OrderAddressTag         *CustomTag
	OrderStatusTag          *CustomTag
	OrderProcessingTypesTag *CustomTag
	OrderQuantityTag        *CustomTag
	OrderPriceTag           *CustomTag
	OrderExpDateTag         *CustomTag
	OrderDescTag            *CustomTag
	OrderFileUrlTag         *CustomTag
	OrderDetailNameTag      *CustomTag
}

func (p *OrderParser) Parse(wg *sync.WaitGroup) {
	defer wg.Done()
	p.ResourceCache = make(map[string]models.Resource)
	p.Settings.Postfix = p.Postfix
	htmlDoc, err := p.Settings.GetHtmlDocument()

	if err != nil {
		p.Settings.Logger.Warning(fmt.Sprintf("Url is broken: %s%s", p.Settings.BaseUrl, p.Postfix))
		return
	}

	p.Settings.StartPoint = GetStartPoint(htmlDoc, p.StartPointTag)
	p.Settings.EndPoint = GetEndPoint(htmlDoc, p.EndPointTag)
	var detailParseWg sync.WaitGroup
	var pushWg sync.WaitGroup

	for i := p.Settings.StartPoint; i <= p.Settings.EndPoint; i++ {
		p.Settings.Postfix = p.Postfix + strconv.Itoa(i)
		htmlDoc, err = p.Settings.GetHtmlDocument()
		if err != nil {
			p.Settings.Logger.Warning(fmt.Sprintf("Order Url is broken: %s%s", p.Settings.BaseUrl, p.Postfix))
			continue
		}

		htmlDoc.Find(p.DetailOrderTag.Container).Each(func(i int, s *goquery.Selection) {
			url, _ := s.Find(p.DetailOrderTag.Tag).Attr(p.DetailOrderTag.Attribute)
			detailChan := make(chan OrderCard)
			detailParseWg.Add(1)
			go p.detailParse(url, detailChan, &detailParseWg)

			for card := range detailChan {
				pushWg.Add(1)
				go func() {
					err := p.pushOrder(card, &pushWg)

					if err != nil {
						p.Settings.Logger.Warning(err.Error())
					}
				}()
			}

		})
	}
	detailParseWg.Wait()
	pushWg.Wait()

	p.Settings.Logger.Info("Order parsing complete on " + p.Settings.BaseUrl + p.Postfix)
}

func (p *OrderParser) detailParse(url string, dc chan OrderCard, wg *sync.WaitGroup) {
	defer func() {
		close(dc)
		wg.Done()
	}()

	var parser OrderParser
	parser.Settings.BaseUrl = url
	htmlDoc, err := parser.Settings.GetHtmlDocument()
	isWrong := false

	if err != nil {
		isWrong = true
		if strings.Index(url, "/") == 0 {
			url = url[1:]
		}
		parser.Settings.BaseUrl = p.Settings.BaseUrl + url
		htmlDoc, err = parser.Settings.GetHtmlDocument()
		if err != nil {
			isWrong = true
		} else {
			isWrong = false
		}

	}
	if isWrong {
		p.Settings.Logger.Warning(fmt.Sprintf("Provider Url is broken"))
		return
	}

	_, ok := p.ResourceCache[p.Settings.BaseUrl]

	if ok == false {
		res, err := GetResource(p.Settings.Database, p.Settings.BaseUrl)

		if res == nil {
			p.Settings.Logger.Warning("Order resource is nil")
		}
		if err != nil {
			p.Settings.Logger.Warning("Order error with resource")
		}

		p.ResourceCache[p.Settings.BaseUrl] = *res
	}

	resource, _ := p.ResourceCache[p.Settings.BaseUrl]

	order := models.Order{
		ResourceID:      resource.Id,
		PublicationDate: FormatString(GetParameterByTag(p.OrderPubDateTag, htmlDoc)),
		Address:         FormatString(GetParameterByTag(p.OrderAddressTag, htmlDoc)),
		Quantity:        FormatString(GetParameterByTag(p.OrderQuantityTag, htmlDoc)),
		Price:           FormatString(GetParameterByTag(p.OrderPriceTag, htmlDoc)),
		ExpirationDate:  FormatString(GetParameterByTag(p.OrderExpDateTag, htmlDoc)),
		Description:     FormatString(GetParameterByTag(p.OrderDescTag, htmlDoc)),
		DownloadFileUrl: FormatString(GetParameterByTag(p.OrderFileUrlTag, htmlDoc)),
		NameDetail:      FormatString(GetParameterByTag(p.OrderDetailNameTag, htmlDoc)),
		Status:          FormatString(GetParameterByTag(p.OrderStatusTag, htmlDoc)),
	}

	procTypes := FormatString(GetParameterByTag(p.OrderProcessingTypesTag, htmlDoc))
	types := strings.Split(procTypes, ",")
	order.ProcessingTypes = types

	if p.SuggestionParser == nil {
		dc <- OrderCard{
			Order:       &order,
			Suggestions: nil,
		}
		return
	}

	var suggestions []models.Suggestion

	selection := htmlDoc.Find(p.SuggestionParser.Container)

	if selection.Size() == 0 {
		dc <- OrderCard{
			Order:       &order,
			Suggestions: nil,
		}
		return
	}

	selection.Each(func(i int, s *goquery.Selection) {
		time := p.SuggestionParser.TimeTag.GetFromSelection(s)
		price := p.SuggestionParser.PriceTag.GetFromSelection(s)
		comment := p.SuggestionParser.CommentTag.GetFromSelection(s)
		email := p.SuggestionParser.EmailTag.GetFromSelection(s)
		phone := p.SuggestionParser.PhoneTag.GetFromSelection(s)
		suggestions = append(suggestions, models.Suggestion{
			Time:    FormatString(time),
			Price:   FormatString(price),
			Comment: FormatString(comment),
			Email:   FormatString(email),
			Phone:   FormatString(phone),
		})
	})

	dc <- OrderCard{
		Order:       &order,
		Suggestions: &suggestions,
	}
}

func (p *OrderParser) pushOrder(data OrderCard, wg *sync.WaitGroup) error {
	defer wg.Done()

	var givenOrder models.Order
	result := p.Settings.Database.Where(data.Order).FirstOrCreate(&givenOrder)

	if err := result.Error; err != nil {
		return err
	}
	/*
		if result.RowsAffected == 0 {
			return errors.New("no rows affected")
		}
	*/

	if data.Suggestions == nil {
		return nil
	}

	for _, element := range *data.Suggestions {
		element.OrderID = givenOrder.Id
		result := p.Settings.Database.Create(&element)

		if err := result.Error; err != nil {
			return err
		}
	}

	return nil
}
