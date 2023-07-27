package instance

import (
	"gorm.io/gorm"
	"yellot-parser/main/core/parser/parser_services"
	"yellot-parser/main/core/parser/parser_settings"
	"yellot-parser/main/logger"
)

func MetaloobrabotchikiParser(iLogger logger.ILogger, db *gorm.DB) *parser_services.Parser {
	parser := parser_services.NewParserBuilder().
		Provider().Is(parser_services.ProviderParser{
		Postfix: "company/",
		StartPointTag: parser_services.NewTagBuilder().
			Contained().In("li.page-item.active").
			Build().SetTag("span"),
		EndPointTag: parser_services.NewTagBuilder().
			Contained().In("li.page-item").
			Attribute().Is("data-ci-pagination-page").
			Build().SetTag("a"),
		DetailCardTag: parser_services.NewTagBuilder().
			Contained().In("div.btn_card").
			Attribute().Is("href").
			Build().SetTag("a"),
		CompanyNameTag: parser_services.NewTagBuilder().
			Contained().In(".title .col-md-10").
			Build().SetTag("h1"),
		TaxpayerTag: parser_services.NewTagBuilder().
			Build().SetTag(".contact-card p:nth-child(1)"),
		CompanyAddressTag: parser_services.NewTagBuilder().
			Build().SetTag(".contact-card p:nth-child(2)"),
		CompanyDescTag: parser_services.NewTagBuilder().
			Contained().In(".company_card_text").
			Build().SetTag("p"),
		CompanyServicesTag: parser_services.NewTagBuilder().
			Contained().In(".company_card_service").
			Build().SetTag("li a"),
	}).Order().Is(parser_services.OrderParser{
		Postfix: "orders/",
		DetailOrderTag: parser_services.NewTagBuilder().
			Contained().In("div.order_card").
			Attribute().Is("href").
			Build().SetTag("a.btn-order"),
		OrderDetailNameTag: parser_services.NewTagBuilder().
			Contained().In(".title .col-10").
			Build().SetTag("h1"),
		OrderDescTag: parser_services.NewTagBuilder().
			Contained().In("div.card-text").
			Build().SetTag("p"),
	}).Settings().Is(parser_settings.ParserSettings{
		BaseUrl:  "https://metalloobrabotchiki.ru/",
		Logger:   iLogger,
		Database: db,
	}).Build()

	return parser
}
