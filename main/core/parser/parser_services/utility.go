package parser_services

import (
	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"yellot-parser/main/models"
)

func GetStartPoint(node *goquery.Document, startPointTag *CustomTag) int {
	if startPointTag == nil {
		return 1
	}

	if startPointTag.IsHaveContainer() {
		if startPointTag.IsHaveAttribute() {
			startPoint, _ := node.
				Find(startPointTag.Container).
				Find(startPointTag.Tag).
				Attr(startPointTag.Attribute)

			i, err := strconv.Atoi(startPoint)

			if err == nil {
				return i
			} else {
				return 1
			}
		} else {
			startPoint := node.
				Find(startPointTag.Container).
				Find(startPointTag.Tag).
				Text()

			i, err := strconv.Atoi(startPoint)

			if err == nil {
				return i
			} else {
				return 1
			}
		}
	} else {
		if startPointTag.IsHaveAttribute() {
			startPoint, _ := node.
				Find(startPointTag.Tag).
				Attr(startPointTag.Attribute)

			i, err := strconv.Atoi(startPoint)

			if err == nil {
				return i
			} else {
				return 1
			}
		} else {
			startPoint := node.
				Find(startPointTag.Tag).
				Text()

			i, err := strconv.Atoi(startPoint)

			if err == nil {
				return i
			} else {
				return 1
			}
		}
	}
}

func GetEndPoint(node *goquery.Document, endPointTag *CustomTag) int {
	if endPointTag == nil {
		return 1
	}

	max := 1
	if endPointTag.IsHaveContainer() {
		if endPointTag.IsHaveAttribute() {
			node.Find(endPointTag.Container).Each(func(i int, s *goquery.Selection) {
				title, _ := s.Find(endPointTag.Tag).Attr(endPointTag.Attribute)
				i, err := strconv.Atoi(title)
				if err == nil {
					if i > max {
						max = i
					}
				}
			})
			return max
		} else {
			node.Find(endPointTag.Container).Each(func(i int, s *goquery.Selection) {
				title := s.Find(endPointTag.Tag).Text()
				i, err := strconv.Atoi(title)
				if err == nil {
					if i > max {
						max = i
					}
				}
			})
			return max
		}
	} else {
		if endPointTag.IsHaveAttribute() {
			endPoint, _ := node.
				Find(endPointTag.Tag).
				Attr(endPointTag.Attribute)

			i, err := strconv.Atoi(endPoint)

			if err == nil {
				return i
			} else {
				return 1
			}
		} else {
			endPoint := node.
				Find(endPointTag.Tag).
				Text()

			i, err := strconv.Atoi(endPoint)
			if err == nil {
				return i
			} else {
				return 1
			}
		}
	}
}

func GetParameterByTag(tag *CustomTag, htmlDoc *goquery.Document) string {
	if tag == nil {
		return ""
	}
	if tag.IsHaveContainer() {
		if tag.IsHaveAttribute() {
			var attr, _ = htmlDoc.
				Find(tag.Container).
				Find(tag.Tag).
				Children().
				Remove().
				End().
				Attr(tag.Attribute)
			return attr
		} else {
			return htmlDoc.
				Find(tag.Container).
				Find(tag.Tag).
				Children().
				Remove().
				End().
				Text()
		}
	} else {
		if tag.IsHaveAttribute() {
			var attr, _ = htmlDoc.
				Find(tag.Tag).
				Children().
				Remove().
				End().
				Attr(tag.Attribute)
			return attr
		} else {
			return htmlDoc.Find(tag.Tag).
				Children().
				Remove().
				End().
				Text()
		}
	}
}

func FormatString(s string) string {
	return strings.Trim(regexp.MustCompile(`\s+`).ReplaceAllString(s, " "), " ")
}

func GetResource(db *gorm.DB, address string) (*models.Resource, error) {
	var resource models.Resource

	result := db.Where(models.Resource{Address: address}).FirstOrCreate(&resource)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &resource, nil
}

func StripHtmlRegex(s string) string {
	r := regexp.MustCompile(`<.*?>`)
	return r.ReplaceAllString(s, "")
}
