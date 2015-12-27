package main

import (
	"fmt"
	"time"

	"github.com/pestkam/scraper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	// Parsing City List
	session, err := mgo.Dial("192.168.100.6")
	if err != nil {
		panic(err)
	}
	scrapCity := scraper.NewScraper(1, 3)
	scrapCity.AddLink("http://2gis.ru/countries/global/")
	go scrapCity.RunCrawler()
	for result := range scrapCity.Results {
		if result.Err != nil {
			fmt.Println(result.Err)
			continue
		}
		ParseCity(result, session)
	}
	//Parsing Rubrics List
	scrapRubric := scraper.NewScraper(1, 3)

	colCity := session.DB("2Gis").C("city")
	RubricProfile := bson.M{}
	date := time.Now().Add(-time.Hour * 24 * 14)
	iter := colCity.Find(bson.M{"$or": [2]bson.M{bson.M{"lastUpdate": bson.M{`$exists`: false}}, bson.M{"lastUpdate": bson.M{`$lt`: date}}}}).Iter()
	for iter.Next(&RubricProfile) {
		scrapRubric.AddLink(fmt.Sprint(RubricProfile["url"]))
	}
	session.Clone()
	go scrapRubric.RunCrawler()
	for result := range scrapRubric.Results {
		if result.Err != nil {
			fmt.Println(result.Err)
			continue
		}
		ParseRubrics(result, session)
	}
	//
	//Parsing SubRubrics List
	scrapSubRubric := scraper.NewScraper(4, 3)
	colRubrics := session.DB("2Gis").C("rubrics")
	SubRubricProfile := bson.M{}
	date = time.Now().Add(-time.Hour * 24 * 14)
	iter = colRubrics.Find(bson.M{"$or": [2]bson.M{bson.M{"lastUpdate": bson.M{`$exists`: false}}, bson.M{"lastUpdate": bson.M{`$lt`: date}}}}).Iter()
	for iter.Next(&SubRubricProfile) {
		scrapSubRubric.AddLink(fmt.Sprint(SubRubricProfile["url"]))
	}
	go scrapSubRubric.RunCrawler()
	for result := range scrapSubRubric.Results {
		if result.Err != nil {
			fmt.Println(result.Err)
			continue
		}
		ParseSubRubrics(result, session)
	}

	// Parsing Company List
	scrapCompanyList := scraper.NewScraper(1, 3)
	colSubRubrics := session.DB("2Gis").C("subrubrics")
	CompanyListProfile := bson.M{}
	date = time.Now().Add(-time.Hour * 24 * 14)
	iter = colSubRubrics.Find(bson.M{"$or": [2]bson.M{bson.M{"lastUpdate": bson.M{`$exists`: false}}, bson.M{"lastUpdate": bson.M{`$lt`: date}}}}).Iter()
	for iter.Next(&CompanyListProfile) {
		scrapCompanyList.AddLink(fmt.Sprint(CompanyListProfile["url"]))
	}
	go scrapCompanyList.RunCrawler()
	for result := range scrapCompanyList.Results {
		if result.Err != nil {
			fmt.Println(result.Err)
			continue
		}
		ParseCompanyList(result, session)
	}

	//Parsing Company Info
	iter = colCity.Find(bson.M{"lastUpdate": bson.M{`$exists`: true}}).Iter()
	cityList := []string{}
	for iter.Next(&RubricProfile) {
		cityList = append(cityList, RubricProfile["latincityname"].(string))
	}
	for _, latincityname := range cityList {
		colCompanyList := session.DB("2Gis").C(latincityname)
		scrapCompanyInfo := scraper.NewScraper(4, 3)
		CompanyProfile := bson.M{}
		date := time.Now().Add(-time.Hour * 24 * 14)
		iter := colCompanyList.Find(bson.M{"$or": [2]bson.M{bson.M{"lastUpdate": bson.M{`$exists`: false}}, bson.M{"lastUpdate": bson.M{`$lt`: date}}}}).Iter()
		for iter.Next(&CompanyProfile) {
			scrapCompanyInfo.AddLink(fmt.Sprint(CompanyProfile["url"]))
		}
		go scrapCompanyInfo.RunCrawler()
		for result := range scrapCompanyInfo.Results {
			if result.Err != nil {
				fmt.Println(result.Err)
				continue
			}
			ParseCompanyInfo(result, session)
		}
	}

}
