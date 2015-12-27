package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/PuerkitoBio/goquery"
	"github.com/pestkam/scraper"
)

func ParseSubRubrics(resp scraper.Response, session *mgo.Session) {
	domain := getShemaAndDomain(resp.Request.URL.String())
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {

		colRubrics := session.DB("2Gis").C("rubrics")
		colSubRubrics := session.DB("2Gis").C("subrubrics")
		var rubricProfile = bson.M{}
		change := mgo.Change{
			Update:    bson.M{"$set": bson.M{"lastUpdate": time.Now()}},
			ReturnNew: false,
		}

		_, err = colRubrics.Find(bson.M{"url": resp.Request.URL.String()}).Apply(change, &rubricProfile)
		cityname := rubricProfile["cityname"]
		latincityname := rubricProfile["latincityname"]
		rubric := rubricProfile["rubric"]

		doc.Find("li.rubricsList__listItem").Each(func(i int, item *goquery.Selection) {
			subrubric := item.Find("a.rubricsList__listItemLinkTitle").Text()
			href, _ := item.Find("a.rubricsList__listItemLinkTitle").Attr("href")
			countCompany := getCompanyCount(item.Find("span.rubricsList__listItemDescription").Text())
			URL := domain + href

			currentLink := getCuttetURL(URL)
			if countCompany > 12 {
				var lastPage int
				lastPage = int(countCompany / 12)
				if (countCompany % 12) > 0 {
					lastPage++
				}
				for currentPage := 1; currentPage <= lastPage; currentPage++ {
					url := fmt.Sprintf("%s/page/%d", currentLink, currentPage)
					subrubricProfile := bson.M{"$set": bson.M{"url": url, "subrubric": subrubric, "rubric": rubric, "cityname": cityname, "latincityname": latincityname}}
					_, err = colSubRubrics.Upsert(bson.M{"url": url}, subrubricProfile)
					if err != nil {
						fmt.Println(err)
					}
				}
			} else {
				subrubricProfile := bson.M{"$set": bson.M{"url": currentLink, "subrubric": subrubric, "rubric": rubric, "cityname": cityname, "latincityname": latincityname}}
				_, err = colSubRubrics.Upsert(bson.M{"url": currentLink}, subrubricProfile)
				if err != nil {
					fmt.Println(err)
				}
			}
		})
	}
}
