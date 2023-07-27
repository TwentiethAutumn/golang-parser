package instance

import (
	"gorm.io/gorm"
	"yellot-parser/main/core/parser/parser_services"
	"yellot-parser/main/core/parser/parser_settings"
	"yellot-parser/main/logger"
)

func ObrabotkaNetArchiveParser(iLogger logger.ILogger, db *gorm.DB) *parser_services.Parser {
	parser := parser_services.NewParserBuilder().
		Order().Is(parser_services.OrderParser{
		Postfix: "orders/archive/?page=",
		StartPointTag: parser_services.NewTagBuilder().
			Contained().In(".pagination").
			Build().SetTag("li .active"),
		EndPointTag: parser_services.NewTagBuilder().
			Contained().In(".pagination li").
			Build().SetTag("a"),
		DetailOrderTag: parser_services.NewTagBuilder().
			Contained().In(".panel-heading").
			Attribute().Is("href").
			Build().SetTag("a"),
		OrderDetailNameTag: parser_services.NewTagBuilder().
			Contained().In(".page-header").
			Build().SetTag("h1"),
		OrderDescTag: parser_services.NewTagBuilder().
			Build().SetTag("dl.dl-horizontal dd:nth-child(12) p"),
		OrderAddressTag: parser_services.NewTagBuilder().
			Build().SetTag("dl.dl-horizontal dd:nth-child(6)"),
		OrderProcessingTypesTag: parser_services.NewTagBuilder().
			Build().SetTag("dl.dl-horizontal dd:nth-child(4) a"),
		OrderFileUrlTag: parser_services.NewTagBuilder().
			Contained().In(".col-md-4.col-md-push-8.order-files").
			Attribute().Is("href").
			Build().SetTag("a"),
		SuggestionParser: &parser_services.SuggestionParser{
			Container: ".panel.panel-default",
			TimeTag: parser_services.NewTagBuilder().
				Divided().By("br").Left().
				Build().SetTag(".col-sm-6:nth-child(1)"),
			CommentTag: parser_services.NewTagBuilder().
				Build().SetTag(".panel-body p"),
			PriceTag: parser_services.NewTagBuilder().
				Divided().By("br").Right().
				Build().SetTag(".col-sm-6:nth-child(1)"),
			EmailTag: parser_services.NewTagBuilder().
				Build().SetTag(".row .col-sm-6:nth-child(2) a:nth-child(5)"),
			PhoneTag: parser_services.NewTagBuilder().
				Build().SetTag(".row .col-sm-6:nth-child(2) a:nth-child(2)"),
		},
	}).Settings().Is(parser_settings.ParserSettings{
		BaseUrl:  "https://obrabotka.net/",
		Logger:   iLogger,
		Database: db,
	}).Build()

	return parser
}
