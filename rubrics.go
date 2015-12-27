package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/PuerkitoBio/goquery"
	"github.com/pestkam/scraper"
)

func ParseRubrics(resp scraper.Response, session *mgo.Session) {
	domain := getShemaAndDomain(resp.Request.URL.String())
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		colCity := session.DB("2Gis").C("city")
		colRubrics := session.DB("2Gis").C("rubrics")

		var cityProfile = bson.M{}
		change := mgo.Change{
			Update:    bson.M{"$set": bson.M{"lastUpdate": time.Now()}},
			ReturnNew: true,
		}

		_, err = colCity.Find(bson.M{"url": resp.Request.URL.String()}).Apply(change, &cityProfile)
		cityname := cityProfile["cityname"]
		latincityname := cityProfile["latincityname"]
		doc.Find("a.rubricsList__listItemLinkTitle").Each(func(i int, item *goquery.Selection) {
			rubric := item.Text()
			href, _ := item.Attr("href")
			rubricProfile := bson.M{"$set": bson.M{"url": domain + href, "rubric": rubric, "cityname": cityname, "latincityname": latincityname}}
			_, err = colRubrics.Upsert(bson.M{"url": domain + href}, rubricProfile)
			if err != nil {
				fmt.Println(err)
			}
		})
	}
}
