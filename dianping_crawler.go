package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
	log "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

var DianpingCrawler *Dianping

func init() {
	DianpingCrawler = &Dianping{}
}

var shopReg = regexp.MustCompile(`/shop/([0-9]+)`)
var storage = &redisstorage.Storage{
	Address:  fmt.Sprintf("%s:%s", CELERY_BACKEND_HOST, CELERY_BACKEND_PORT),
	Password: CELERY_BACKEND_PASSWORD,
	DB:       0,
	Prefix:   "creeper",
}

type Dianping struct {
}

type Shop struct {
	Url       string
	Id        string
	CrawledAt *time.Time
	Contact   string
	Title     string
	Images    []string
}

func (*Dianping) crawl(task *CrawlerTask) interface{} {
	log.Infof("run dianping crawler...")
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: www.dianping.com
		colly.AllowedDomains(task.Url.Host),
		// Turn on asynchronous requests
		colly.Async(true),
	)
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})

	var result = make(map[string]interface{})
	result["urls"] = make([]string, 0, 100)
	var urlReg = regexp.MustCompile(fmt.Sprintf("%s%s%s",
		task.Url, Arr2RegOr(task.SubUrls, "/", ""), Arr2RegOr(task.SubUrls2, "/", "")))

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absUrl := e.Request.AbsoluteURL(link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		var shouldVisit = false
		// Print link
		matchArr := shopReg.FindStringSubmatch(link)
		// log.Infof("Link found: text: %q, link: %s, absUrl: %q, matchArr: %q",
		//     e.Text, link, absUrl, matchArr)
		if len(matchArr) == 2 || urlReg.Match([]byte(absUrl)) {
			result["urls"] = append(result["urls"].([]string), absUrl)
			shouldVisit = true
		}
		if shouldVisit {
			c.Visit(absUrl)
		} else {
			log.Warningf("why not visit me? url: %v, href: %s", e.Request.URL, link)
		}
	})
	c.OnHTML("div[class=page]", func(e *colly.HTMLElement) {
		log.Infof("分页url: %q", e.Request.URL)
		// shop := &Shop{}
		// shop. = e.ChildAttr("a[data-event-action=title]", "href")
		links := e.ChildAttrs("a[href]", "href")
		for _, link := range links {
			log.Infof("visit page: %s", link)
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})
	c.OnHTML("div[class=tit]", func(e *colly.HTMLElement) {
		log.Infof("listing页面url: %q", e.Request.URL)
		// shop := &Shop{}
		// shop. = e.ChildAttr("a[data-event-action=title]", "href")
		links := e.ChildAttrs("a[href]", "href")
		for _, link := range links {
			log.Infof("visit page: %s", link)
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %s", r.URL.String())
		// cookie := c.Cookies(r.URL.String())
	})
	c.OnError(func(response *colly.Response, e error) {
		log.Errorf("response: %v, error: %v", response, e)
	})

	log.Infof("start visiting %q", task.urls())
	// StartWorker scraping on https://hackerspaces.org
	for _, url := range task.urls() {
		c.Visit(url)
	}
	return result
}
