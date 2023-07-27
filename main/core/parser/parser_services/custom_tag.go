package parser_services

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type CustomTag struct {
	Tag       string
	Container string
	Attribute string
	Divider   string
	LeftSide  bool
}

type TagBuilder struct {
	tag *CustomTag
}

type TagContainerBuilder struct {
	TagBuilder
}

type TagAttributeBuilder struct {
	TagBuilder
}

type TagDividerBuilder struct {
	TagBuilder
}

type TagDividePositionBuilder struct {
	TagDividerBuilder
}

func NewTagBuilder() *TagBuilder {
	return &TagBuilder{tag: &CustomTag{}}
}

func (b *TagBuilder) Contained() *TagContainerBuilder {
	return &TagContainerBuilder{*b}
}

func (b *TagBuilder) Attribute() *TagAttributeBuilder {
	return &TagAttributeBuilder{*b}
}

func (b *TagBuilder) Divided() *TagDividerBuilder {
	return &TagDividerBuilder{*b}
}

func (c *TagContainerBuilder) In(tag string) *TagContainerBuilder {
	c.tag.Container = tag
	return c
}

func (a *TagAttributeBuilder) Is(attr string) *TagAttributeBuilder {
	a.tag.Attribute = attr
	return a
}

func (d *TagDividerBuilder) By(tag string) *TagDividePositionBuilder {
	d.tag.Divider = tag
	return &TagDividePositionBuilder{*d}
}

func (dp *TagDividePositionBuilder) Left() *TagDividePositionBuilder {
	dp.tag.LeftSide = true
	return dp
}

func (dp *TagDividePositionBuilder) Right() *TagDividePositionBuilder {
	dp.tag.LeftSide = false
	return dp
}

func (b *TagBuilder) Build() *CustomTag {
	return b.tag
}

func (t *CustomTag) IsHaveContainer() bool {
	return t.Container != ""
}

func (t *CustomTag) IsHaveAttribute() bool {
	return t.Attribute != ""
}

func (t *CustomTag) IsHaveDivision() bool {
	return t.Divider != ""
}

func (t *CustomTag) SetTag(tag string) *CustomTag {
	t.Tag = tag
	return t
}

func (t *CustomTag) GetFromSelection(s *goquery.Selection) string {
	if t == nil {
		return ""
	}
	if t.IsHaveContainer() {
		if t.IsHaveAttribute() {
			var attr, _ = s.
				Find(t.Container).
				Find(t.Tag).
				Children().
				Remove().
				End().
				Attr(t.Attribute)
			return attr
		} else {
			if t.IsHaveDivision() {
				if t.LeftSide {
					htdoc, _ := goquery.OuterHtml(s.Find(t.Container).Find(t.Tag))
					index := strings.Index(htdoc, t.Divider)
					if index != -1 {
						htdoc = htdoc[:index-1]
					}
					return StripHtmlRegex(htdoc)
				}
				htdoc, _ := goquery.OuterHtml(s.Find(t.Container).Find(t.Tag))
				index := strings.Index(htdoc, t.Divider)
				if index != -1 {
					htdoc = htdoc[index-1:]
				}
				return StripHtmlRegex(htdoc)
			}
			return s.
				Find(t.Container).
				Find(t.Tag).
				Children().
				Remove().
				End().
				Text()
		}
	} else {
		if t.IsHaveAttribute() {
			var attr, _ = s.
				Find(t.Tag).
				Children().
				Remove().
				End().
				Attr(t.Attribute)
			return attr
		} else {
			if t.IsHaveDivision() {
				if t.LeftSide {
					htdoc, _ := goquery.OuterHtml(s.Find(t.Tag))
					index := strings.Index(htdoc, t.Divider)
					if index != -1 {
						htdoc = htdoc[:index-1]
					}
					return StripHtmlRegex(htdoc)
				}
				htdoc, _ := goquery.OuterHtml(s.Find(t.Tag))
				index := strings.Index(htdoc, t.Divider)
				if index != -1 {
					htdoc = htdoc[index-1:]
				}
				return StripHtmlRegex(htdoc)
			}
			return s.
				Find(t.Tag).
				Children().
				Remove().
				End().
				Text()
		}
	}
}
