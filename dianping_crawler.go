package main

import (
    "fmt"
    "github.com/Danceiny/go.fastjson"
    . "github.com/Danceiny/go.utils"
    "github.com/PuerkitoBio/goquery"
    "github.com/gocolly/colly"
    log "github.com/sirupsen/logrus"
    "math/rand"
    "os"
    "regexp"
    "strings"
    "sync"
    "time"
)

type Dianping struct {
    MaxPageNum int
}

var shopReg = regexp.MustCompile(`\.com/shop/([0-9]+)$`)

func (self *Dianping) crawl(task *CrawlerTask) interface{} {
    log.Infof("run dianping crawler...")
    var fn = task.ResultFilename()
    var err error
    var f *os.File
    if _, err := os.Stat(fn); os.IsNotExist(err) {
        f, err = os.Create(fn)
        // _, err = f.WriteString("fuck you")
        // PanicError(err)
    } else {
        f, err = os.OpenFile(fn, os.O_WRONLY|os.O_APPEND, 0666)
    }
    PanicError(err)
    // Instantiate default collector
    c := colly.NewCollector(
        // Visit only domains: www.dianping.com
        colly.AllowedDomains(task.Url.Host),
        // Turn on asynchronous requests
        colly.Async(OPEN_ASYNC_MODE),
        colly.AllowURLRevisit(),
    )
    err = c.Limit(&colly.LimitRule{
        DomainGlob:  "*",
        Parallelism: 2,
        RandomDelay: 10 * time.Second,
    })
    PanicError(err)
    // use proxy switcher
    if PROXYS == nil || len(PROXYS) == 0 {
        PROXYS = GetAllProxyUrl()
    }
    log.Infof("proxys size: %d, example: %s", len(PROXYS), PROXYS[0])
    var rp colly.ProxyFunc
    rp, err = RoundRobinProxySwitcher(PROXYS...)
    // PanicError(err)
    c.SetProxyFunc(rp)
    // add storage to the collector
    err = c.SetStorage(storage)
    if err != nil {
        panic(err)
    }

    var resultMutex = sync.RWMutex{}
    visitedUrls := make([]string, 0, 100)
    visitedShops := make(map[string]*Shop)

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
            resultMutex.Lock()
            visitedUrls = append(visitedUrls, absUrl)
            resultMutex.Unlock()
            shouldVisit = true
        }
        if shouldVisit {
            _ = c.Visit(absUrl)
        } else {
            log.Warningf("why not visit me? url: %v, href: %s", e.Request.URL, link)
        }
    })

    c.OnHTML("div[class=tit]", func(e *colly.HTMLElement) {
        // shop := &Shop{}
        // shop. = e.ChildAttr("a[data-event-action=title]", "href")
        links := e.ChildAttrs("a[href]", "href")
        for _, link := range links {
            absLink := e.Request.AbsoluteURL(link)
            log.Infof("访问详情页: %s", absLink)
            _ = c.Visit(absLink)
        }
    })

    c.OnHTML("span[data-phone]", func(e *colly.HTMLElement) {
        url := e.Request.URL.String()
        match := shopReg.FindStringSubmatch(url)
        var id string
        if len(match) == 2 {
            id = match[1]
        } else {
            log.Warningf("没有匹配到商户id， url: %s", url)
            return
        }
        contact := e.Attr("data-phone")
        dom := e.DOM.ParentsUntil("~")
        metaTags := dom.Find("meta")
        var shopName string
        shopName = dom.Find("div[class=shop-name]").First().
            ChildrenFiltered("h1").First().Text()
        log.Infof("抓取到商户[%s]电话：%s", shopName, contact)
        var shop *Shop
        var existed bool
        resultMutex.Lock()
        shop, existed = visitedShops[id]
        if !existed {
            shop = &Shop{Attr: &fastjson.JSONObject{}, Contacts: make([]string, 0, 2)}
            shop.Url = url
            shop.Id = id
            shop.Title = shopName
            now := time.Now()
            shop.CrawledAt = &now
            visitedShops[id] = shop
            metaTags.Each(func(_ int, s *goquery.Selection) {
                // Search for og:type meta tags
                property, _ := s.Attr("name")
                if strings.EqualFold(property, "Keywords") {
                    var1, _ := s.Attr("content")
                    shop.Attr.Put("keywords", var1)
                } else if strings.EqualFold(property, "Description") {
                    var1, _ := s.Attr("content")
                    shop.Attr.Put("description", var1)
                } else if strings.EqualFold(property, "location") {
                    var1, _ := s.Attr("content")
                    shop.Attr.Put("location", var1)
                }
            })
        }
        var b = true
        for _, c := range shop.Contacts {
            if c == contact {
                b = false
                break
            }
        }
        if b {
            shop.Contacts = append(shop.Contacts, contact)
        }
        content := append(fastjson.ToJSON(shop), '\r', '\n')
        _, _ = f.Write(content)
        resultMutex.Unlock()
    })

    // Before making a request print "Visiting ..."
    c.OnRequest(func(r *colly.Request) {
        r.Headers.Set("User-Agent", RandomUserAgent())
        log.Debugf("Visiting %s with proxy: %s, headers: %v",
            r.URL.String(), r.ProxyURL, r.Headers)
    })

    c.OnResponse(func(r *colly.Response) {
        log.Debugf("Visiting %s with proxy: %s, "+
            "responding headers: %v", r.Request.ProxyURL, r.Headers)
    })

    c.OnError(func(response *colly.Response, e error) {
        url := response.Request.URL.String()
        log.Errorf("url:[%s] respond status_code: %d, error: %v",
            url, response.StatusCode, e)
        self.reenter(c, url)
    })

    log.Infof("start visiting %q", task.urls())
    enter(c, task)
    if OPEN_ASYNC_MODE {
        c.Wait()
    }
    result := make([]string, 0, 1024)
    for _, v := range visitedShops {
        result = append(result, v.Contacts...)
    }
    f.Close()
    return result
}

func enter(c *colly.Collector, task *CrawlerTask) {
    var err error
    for _, url := range task.urls() {
        err = c.Visit(url)
        if err != nil {
            log.Error(err)
        }
    }
}

func (self *Dianping) reenter(c *colly.Collector, url string) {
    time.Sleep(10 * time.Second)
    url = regexp.MustCompile(`/p\d+`).ReplaceAllString(url,
        fmt.Sprintf("/p%d", rand.Intn(self.MaxPageNum)))
    _ = c.Visit(url)
}
