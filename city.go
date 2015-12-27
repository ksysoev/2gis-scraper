package main

import (
	"fmt"
	"regexp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/PuerkitoBio/goquery"
	"github.com/pestkam/scraper"
)

type CityInfo struct {
	CityName      string
	LatinCityName string
	URL           string
}

func ParseCity(resp scraper.Response, session *mgo.Session) {
	latinCityNameRegexp := regexp.MustCompile(`\w+$`)
	checkCountryRegexp := regexp.MustCompile(`^/`)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		c := session.DB("2Gis").C("city")
		doc.Find("li.world__listItem").Each(func(i int, item *goquery.Selection) {
			ProfileCity := CityInfo{}
			ProfileCity.CityName = item.Find("a.world__listItemName").Text()
			link, _ := item.Find("a.world__listItemName").Attr("href")
			ProfileCity.LatinCityName = latinCityNameRegexp.FindString(link)
			// if ProfileCity.LatinCityName != "gornoaltaysk" {
			// 	return
			// }

			if checkCountryRegexp.MatchString(link) {
				link = "http://2gis.ru" + link
			}
			ProfileCity.URL = link + "/rubrics"
			_, err = c.Upsert(bson.M{"latincityname": ProfileCity.LatinCityName}, bson.M{"$set": ProfileCity})
			if err != nil {
				fmt.Println(err)
			}
		})
	}
}
