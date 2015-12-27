package main

import (
	"regexp"
	"strconv"
)

func getShemaAndDomain(url string) string {
	getDomainRegexp := regexp.MustCompile(`^(http://[\w.]+)`)
	resultRegexp := getDomainRegexp.FindStringSubmatch(url)
	return resultRegexp[0]
}

func getCompanyCount(itemDesk string) int {
	countRegexp := regexp.MustCompile(`\d+`)
	resultRegexp := countRegexp.FindStringSubmatch(itemDesk)
	countCompany, err := strconv.Atoi(resultRegexp[0])
	if err != nil {
		return 0
	}
	return countCompany
}

func getCuttetURL(url string) string {
	cutURLRegexp := regexp.MustCompile(`\/tab.*$`)
	return cutURLRegexp.ReplaceAllString(url, "")
}

func getCuttedCompanyURL(url string) string {
	cutURLRegexp := regexp.MustCompile(`\?(http://.+)$`)
	resultRegexp := cutURLRegexp.FindAllStringSubmatch(url, 1)
	if len(resultRegexp) >= 1 {
		if len(resultRegexp[0]) >= 1 {
			return resultRegexp[0][1]
		}
		return ""
	}
	return ""

}

func getCityNameFromURL(url string) string {
	cityNameRegexp := regexp.MustCompile(`http://.+?/(\w+?)/firm`)
	resultRegexp := cityNameRegexp.FindStringSubmatch(url)
	return resultRegexp[1]
}
