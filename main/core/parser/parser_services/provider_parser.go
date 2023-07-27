package parser_services

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"sync"
	"yellot-parser/main/core/parser/parser_settings"
	"yellot-parser/main/models"
)

type ProviderParser struct {
	Settings           parser_settings.ParserSettings
	Postfix            string
	StartPointTag      *CustomTag
	EndPointTag        *CustomTag
	DetailCardTag      *CustomTag
	CompanyNameTag     *CustomTag
	CompanyDescTag     *CustomTag
	CompanyServicesTag *CustomTag
	CompanyAddressTag  *CustomTag
	CompanyPhoneTag    *CustomTag
	MailTag            *CustomTag
	TaxpayerTag        *CustomTag
}

func (p *ProviderParser) Parse(wg *sync.WaitGroup) {
	defer wg.Done()
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
			p.Settings.Logger.Warning(fmt.Sprintf("Provider Url is broken: %s%s", p.Settings.BaseUrl, p.Postfix))
			continue
		}

		htmlDoc.Find(p.DetailCardTag.Container).Each(func(i int, s *goquery.Selection) {
			url, _ := s.Find(p.DetailCardTag.Tag).Attr(p.DetailCardTag.Attribute)

			detailChan := make(chan models.Provider)
			detailParseWg.Add(1)
			go p.detailParse(url, detailChan, &detailParseWg)

			for card := range detailChan {
				pushWg.Add(1)
				go func() {
					err := p.pushProvider(&card, &pushWg)
					if err != nil {
						p.Settings.Logger.Warning(err.Error())
					}
				}()
			}
		})
	}
	detailParseWg.Wait()
	pushWg.Wait()

	p.Settings.Logger.Info("Provider parsing complete on " + p.Settings.BaseUrl + p.Postfix)
}

func (p *ProviderParser) detailParse(url string, dc chan models.Provider, wg *sync.WaitGroup) {
	defer func() {
		close(dc)
		wg.Done()
	}()

	var parser ProviderParser
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

	resource, err := GetResource(p.Settings.Database, p.Settings.BaseUrl)
	if resource == nil {
		p.Settings.Logger.Warning("Provider resource is nil")
	}
	if err != nil {
		p.Settings.Logger.Warning("Provider error with resource")
	}

	provider := models.Provider{
		ResourceID:         resource.Id,
		CompanyName:        FormatString(GetParameterByTag(p.CompanyNameTag, htmlDoc)),
		CompanyDescription: FormatString(GetParameterByTag(p.CompanyDescTag, htmlDoc)),
		Address:            FormatString(GetParameterByTag(p.CompanyAddressTag, htmlDoc)),
		Phone:              FormatString(GetParameterByTag(p.CompanyPhoneTag, htmlDoc)),
		TypesOfServices:    strings.Split(FormatString(GetParameterByTag(p.CompanyServicesTag, htmlDoc)), ","),
		TaxpayerId:         FormatString(GetParameterByTag(p.TaxpayerTag, htmlDoc)),
		Email:              FormatString(GetParameterByTag(p.MailTag, htmlDoc)),
	}

	dc <- provider
}

func (p *ProviderParser) pushProvider(order *models.Provider, wg *sync.WaitGroup) error {
	defer wg.Done()

	result := p.Settings.Database.Create(order)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
