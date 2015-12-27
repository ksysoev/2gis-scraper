package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/PuerkitoBio/goquery"
	"github.com/pestkam/scraper"
)

// type PaternCompanyProfile struct {
// 	Rubrics map[string]string
// }

type CompanyProfile struct {
	Name       string
	Address    string
	Worktime   string
	Web        []string
	Social     []string
	Email      string
	Telephones []string
}

func ParseCompanyInfo(resp scraper.Response, session *mgo.Session) {
	p := CompanyProfile{}
	latincityname := getCityNameFromURL(resp.Request.URL.String())

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		//Company Name
		if doc.Find("span.mediaCard__firmCardNameCut").Text() != "" {
			p.Name = doc.Find("span.mediaCard__firmCardNameCut").Text()
		} else if doc.Find("h1.firmCard__name").Text() != "" {
			p.Name = doc.Find("h1.firmCard__name").Text()
		} else if doc.Find("div.cardHeader__headerNameWrapper").Text() != "" {
			p.Name = doc.Find("div.cardHeader__headerNameWrapper").Text()
		}

		//Address
		if doc.Find("a.mediaCard__firmCardAddressName").Text() != "" {
			p.Address = doc.Find("a.mediaCard__firmCardAddressName").Text()
		} else if doc.Find("a.mediaAddress__address").Text() != "" {
			p.Address = doc.Find("a.mediaAddress__address").Text()
		} else if doc.Find("a.firmCard__addressLink").Text() != "" {
			p.Address = doc.Find("a.firmCard__addressLink").Text()
		} else if doc.Find("a.firmCard__geoNameLink").Text() != "" {
			p.Address = doc.Find("a.firmCard__geoNameLink").Text()
		}

		//Web sites
		doc.Find("div.contact__websites a").Each(func(i int, item *goquery.Selection) {
			if href, ok := item.Attr("href"); ok {
				p.Web = append(p.Web, getCuttedCompanyURL(href))
			}
		})

		// Social networks pages
		doc.Find("div.contact__socials a").Each(func(i int, item *goquery.Selection) {
			if href, ok := item.Attr("href"); ok {
				p.Social = append(p.Web, getCuttedCompanyURL(href))
			}
		})

		doc.Find("div.contact__phonesVisible > div > a > span").Each(func(i int, item *goquery.Selection) {
			p.Telephones = append(p.Telephones, item.Text())
		})
		// doc.Find("a.contact__phonesItemLink").Each(func(i int, item *goquery.Selection) {
		// 	p.Telephones = append(p.Telephones, item.Text())
		// })
		// doc.Find("a.mediaContacts__phonesNumber").Each(func(i int, item *goquery.Selection) {
		// 	p.Telephones = append(p.Telephones, item.Text())
		// })

		// doc.Find("a.rubricsList__listItemLinkTitle").Each(func(i int, item *goquery.Selection) {
		//
		// 	rubric := item.Text()
		//
		// 	p.Rubrics["http://2gis.ru"+href] = rubric
		// })

		colCompanyCity := session.DB("2Gis").C(fmt.Sprint(latincityname))

		companyProfile := bson.M{"$set": bson.M{"name": p.Name, "address": p.Address, "worktime": p.Worktime, "web": p.Web, "social": p.Social, "email": p.Email, "telephones": p.Telephones, "lastUpdate": time.Now()}}
		err = colCompanyCity.Update(bson.M{"url": resp.Request.URL.String()}, companyProfile)
		if err != nil {
			fmt.Println(err)
		}
	}
}
