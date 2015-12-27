package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/PuerkitoBio/goquery"
	"github.com/pestkam/scraper"
)

func ParseCompanyList(resp scraper.Response, session *mgo.Session) {
	domain := getShemaAndDomain(resp.Request.URL.String())
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {

		colSubRubrics := session.DB("2Gis").C("subrubrics")

		var subRubricProfile = bson.M{}
		change := mgo.Change{
			Update:    bson.M{"$set": bson.M{"lastUpdate": time.Now()}},
			ReturnNew: true,
		}
		_, err = colSubRubrics.Find(bson.M{"url": resp.Request.URL.String()}).Apply(change, &subRubricProfile)

		cityname := subRubricProfile["cityname"]
		latincityname := subRubricProfile["latincityname"]
		rubric := subRubricProfile["rubric"]
		subrubric := subRubricProfile["subrubric"]
		colCompanyCity := session.DB("2Gis").C(fmt.Sprint(latincityname))

		if doc.Find("a.miniCard__headerTitle").Text() != "" {
			doc.Find("a.miniCard__headerTitle").Each(func(i int, item *goquery.Selection) {
				companyName := item.Text()
				url, _ := item.Attr("href")
				companyProfile := bson.M{"$set": bson.M{"url": domain + url, "name": companyName, "subrubric": subrubric, "rubric": rubric, "cityname": cityname, "latincityname": latincityname}}
				_, err = colCompanyCity.Upsert(bson.M{"url": domain + url}, companyProfile)
				if err != nil {
					fmt.Println(err)
				}
			})
		} else if doc.Find("a.mediaMiniCard__link").Text() != "" {
			doc.Find("a.mediaMiniCard__link").Each(func(i int, item *goquery.Selection) {
				companyName := item.Text()
				url, _ := item.Attr("href")
				companyProfile := bson.M{"$set": bson.M{"url": domain + url, "name": companyName, "subrubric": subrubric, "rubric": rubric, "cityname": cityname, "latincityname": latincityname}}
				_, err = colCompanyCity.Upsert(bson.M{"url": domain + url}, companyProfile)
				if err != nil {
					fmt.Println(err)
				}
			})
		}
	}
}
