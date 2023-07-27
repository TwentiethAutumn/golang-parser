package parser_settings

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"net/http"
	"yellot-parser/main/logger"
)

type ParserSettings struct {
	BaseUrl    string
	Postfix    string
	Logger     logger.ILogger
	Database   *gorm.DB
	StartPoint int
	EndPoint   int
}

func (settings *ParserSettings) GetHtmlDocument() (*goquery.Document, error) {
	var url = settings.BaseUrl

	if settings.Postfix != "" {
		url = fmt.Sprintf("%s%s", url, settings.Postfix)
	}

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	htmlDoc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	return htmlDoc, nil
}
