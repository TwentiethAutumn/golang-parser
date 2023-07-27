package parser_services

import (
	"sync"
	"yellot-parser/main/core/parser/parser_settings"
)

type Parser struct {
	Settings       *parser_settings.ParserSettings
	ProviderParser *ProviderParser
	OrderParser    *OrderParser
}

type ParserBuilder struct {
	parser *Parser
}

type ParserProviderBuilder struct {
	ParserBuilder
}

type ParserOrderBuilder struct {
	ParserBuilder
}

type ParserSettingsBuilder struct {
	ParserBuilder
}

func NewParserBuilder() *ParserBuilder {
	return &ParserBuilder{parser: &Parser{}}
}

func (p *ParserBuilder) Provider() *ParserProviderBuilder {
	return &ParserProviderBuilder{*p}
}

func (p *ParserBuilder) Order() *ParserOrderBuilder {
	return &ParserOrderBuilder{*p}
}

func (p *ParserBuilder) Settings() *ParserSettingsBuilder {
	return &ParserSettingsBuilder{*p}
}

func (p *ParserProviderBuilder) Is(parser ProviderParser) *ParserProviderBuilder {
	p.parser.ProviderParser = &parser
	return p
}

func (p *ParserOrderBuilder) Is(parser OrderParser) *ParserOrderBuilder {
	p.parser.OrderParser = &parser
	return p
}

func (p *ParserSettingsBuilder) Is(settings parser_settings.ParserSettings) *ParserSettingsBuilder {
	p.parser.Settings = &settings
	return p
}

func (p *ParserBuilder) Build() *Parser {
	if p.parser.Settings != nil {
		if p.parser.OrderParser != nil {
			p.parser.OrderParser.Settings = *p.parser.Settings
		}
		if p.parser.ProviderParser != nil {
			p.parser.ProviderParser.Settings = *p.parser.Settings
		}
	}

	return p.parser
}

func (p *Parser) Parse(wg *sync.WaitGroup) {
	defer wg.Done()

	var parseWg sync.WaitGroup

	if p.OrderParser != nil {
		go p.OrderParser.Parse(&parseWg)
		parseWg.Add(1)
	}

	if p.ProviderParser != nil {
		go p.ProviderParser.Parse(&parseWg)
		parseWg.Add(1)
	}

	parseWg.Wait()
}
