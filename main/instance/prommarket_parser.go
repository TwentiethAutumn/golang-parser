package instance

import (
	"gorm.io/gorm"
	"yellot-parser/main/core/parser/parser_services"
	"yellot-parser/main/core/parser/parser_settings"
	"yellot-parser/main/logger"
)

func PromMarketParser(iLogger logger.ILogger, db *gorm.DB) *parser_services.Parser {
	parser := parser_services.NewParserBuilder().
		Provider().Is(parser_services.ProviderParser{
		Postfix: "companies/?PAGEN_1=",
		StartPointTag: parser_services.NewTagBuilder().
			Contained().In(".uk-pagination.uk-flex-center").
			Build().SetTag("li .uk-active"),
		EndPointTag: parser_services.NewTagBuilder().
			Contained().In("ul.uk-pagination.uk-flex-center li").
			Build().SetTag("a"),
		DetailCardTag: parser_services.NewTagBuilder().
			Contained().In("div.card").
			Attribute().Is("href").
			Build().SetTag("a:nth-child(1)"),
		TaxpayerTag: parser_services.NewTagBuilder().
			Build().SetTag(".card-left > div:first-child span:nth-child(2)"),
		CompanyAddressTag: parser_services.NewTagBuilder().
			Build().SetTag(".address span:nth-child(2)"),
		CompanyNameTag: parser_services.NewTagBuilder().
			Build().SetTag(".page-title h1"),
		CompanyPhoneTag: parser_services.NewTagBuilder().
			Contained().In(".information").
			Build().SetTag("a"),
		CompanyDescTag: parser_services.NewTagBuilder().
			Build().SetTag(".detail-bottom__desc"),
		CompanyServicesTag: parser_services.NewTagBuilder().
			Build().SetTag(".uk-accordion-title:not(#reviews-link) > span"),
	}).Order().Is(parser_services.OrderParser{
		Postfix: "metalloobrabotka-na-zakaz/?PAGEN_2=",
		StartPointTag: parser_services.NewTagBuilder().
			Contained().In(".uk-pagination.uk-flex-center").
			Build().SetTag("li .uk-active a"),
		EndPointTag: parser_services.NewTagBuilder().
			Contained().In(".uk-pagination.uk-flex-center").
			Build().SetTag("li a"),
		DetailOrderTag: parser_services.NewTagBuilder().
			Contained().In(".card-middle").
			Attribute().Is("href").
			Build().SetTag("a:nth-child(1)"),
		OrderDetailNameTag: parser_services.NewTagBuilder().
			Build().SetTag(".page-title h1"),
		OrderExpDateTag: parser_services.NewTagBuilder().
			Build().SetTag("div.page-title + .order-details-list div:nth-child(1) .order-details-value:nth-child(1)"),
		OrderQuantityTag: parser_services.NewTagBuilder().
			Build().SetTag("div.page-title + .order-details-list div:nth-child(3) .order-details-value:nth-child(1)"),
		OrderAddressTag: parser_services.NewTagBuilder().
			Build().SetTag("div.page-title + .order-details-list div:nth-child(7) .order-details-value:nth-child(1)"),
		OrderPriceTag: parser_services.NewTagBuilder().
			Build().SetTag("div.page-title + .order-details-list div:nth-child(5) .order-details-value:nth-child(1)"),
		OrderFileUrlTag: parser_services.NewTagBuilder().
			Contained().In(".detail-bottom__file.download_files").
			Attribute().Is("href").
			Build().SetTag("a"),
		OrderDescTag: parser_services.NewTagBuilder().
			Build().SetTag(".detail-bottom__desc"),
		OrderStatusTag: parser_services.NewTagBuilder().
			Build().SetTag(".detail-top__desc-badges .open.lowercase"),
		OrderPubDateTag: parser_services.NewTagBuilder().
			Build().SetTag(".detail-top__desc-badges .open.lowercase + span"),
	}).Settings().Is(parser_settings.ParserSettings{
		BaseUrl:  "https://prom-market.com/",
		Logger:   iLogger,
		Database: db,
	}).Build()
	return parser
}
