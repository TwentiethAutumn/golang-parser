package instance

import (
	"gorm.io/gorm"
	"yellot-parser/main/core/parser/parser_services"
	"yellot-parser/main/core/parser/parser_settings"
	"yellot-parser/main/logger"
)

func MetalPortalParser(iLogger logger.ILogger, db *gorm.DB) *parser_services.Parser {
	parser := parser_services.NewParserBuilder().
		Provider().Is(parser_services.ProviderParser{
		Postfix: "katalog?page=",
		StartPointTag: parser_services.NewTagBuilder().
			Contained().In(".list-inline-item").
			Build().SetTag("span"),
		EndPointTag: parser_services.NewTagBuilder().
			Contained().In("li.list-inline-item").
			Build().SetTag("a"),
		DetailCardTag: parser_services.NewTagBuilder().
			Contained().In(".card-body").
			Attribute().Is("href").
			Build().SetTag("h3 a"),
		CompanyNameTag: parser_services.NewTagBuilder().
			Contained().In(".card-body").
			Build().SetTag("h1 a"),
		CompanyDescTag: parser_services.NewTagBuilder().
			Contained().In(".card-body").
			Build().SetTag("p"),
		CompanyAddressTag: parser_services.NewTagBuilder().
			Contained().In(".col div.card.order").
			Build().SetTag(".card-title + div.row.mb-3 :nth-child(2)"),
		CompanyPhoneTag: parser_services.NewTagBuilder().
			Contained().In(".col div.card.order").
			Build().SetTag("div a.text-dark.border-0"),
		CompanyServicesTag: parser_services.NewTagBuilder().
			Build().SetTag(".row.mb-4 + .row.mb-4 > .col .card-body"),
	}).
		Order().Is(parser_services.OrderParser{
		Postfix: "zakazi?page=",
		StartPointTag: parser_services.NewTagBuilder().
			Contained().In(".list-inline-item").
			Build().SetTag("span"),
		EndPointTag: parser_services.NewTagBuilder().
			Contained().In("li.list-inline-item").
			Build().SetTag("a"),
		DetailOrderTag: parser_services.NewTagBuilder().
			Contained().In("div.card-body").
			Attribute().Is("href").
			Build().SetTag("h5 a"),
		OrderStatusTag: parser_services.NewTagBuilder().
			Contained().In("h1.h5.card-title").
			Build().SetTag("span.mr-2"),
		OrderPubDateTag: parser_services.NewTagBuilder().
			Build().SetTag(".mb-4 + .mb-3 .col-md-9"),
		OrderAddressTag: parser_services.NewTagBuilder().
			Build().SetTag(".mb-4 + .mb-3 + .mb-3 .col-md-9"),
		OrderProcessingTypesTag: parser_services.NewTagBuilder().
			Build().SetTag(".mb-4 + .mb-3 + .mb-3 + .mb-3 .col-md-9"),
		OrderQuantityTag: parser_services.NewTagBuilder().
			Build().SetTag(".mb-4 + .mb-3 + .mb-3 + .mb-3 + .mb-3 .col-md-9"),
		OrderPriceTag: parser_services.NewTagBuilder().
			Build().SetTag(".mb-4 + .mb-3 + .mb-3 + .mb-3 + .mb-3 + .mb-3 .col-md-9"),
		OrderExpDateTag: parser_services.NewTagBuilder().
			Build().SetTag(".mb-4 + .mb-3 + .mb-3 + .mb-3 + .mb-3 + .mb-3 + .mb-3 .col-md-9"),
		OrderDescTag: parser_services.NewTagBuilder().
			Contained().In("div.row.mb-3 div.col-md-9").
			Build().SetTag("p"),
		OrderFileUrlTag: parser_services.NewTagBuilder().
			Contained().In(".col-md-12").
			Attribute().Is("href").
			Build().SetTag("a"),
		OrderDetailNameTag: parser_services.NewTagBuilder().
			Contained().In("h1.h5.card-title").
			Build().SetTag("span.text-dark"),
	}).
		Settings().Is(
		parser_settings.ParserSettings{
			BaseUrl:  "https://metallportal.com/",
			Logger:   iLogger,
			Database: db,
		}).Build()

	return parser
}
