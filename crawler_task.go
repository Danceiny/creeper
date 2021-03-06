package main

import (
    "fmt"
    utils "github.com/Danceiny/go.utils"
    log "github.com/sirupsen/logrus"
    "math/rand"
    urllib "net/url"
    "regexp"
    "time"
)

type CrawlerTask struct {
    SiteName string
    Url      *urllib.URL
    SubUrls  []string
    SubUrls2 []string
}

func (task *CrawlerTask) urlsCount() int {
    var l1, l2 int
    if task.SubUrls2 == nil {
        l2 = 0
    } else {
        l2 = len(task.SubUrls2)
    }
    if task.SubUrls == nil {
        l1 = 0
    } else {
        l1 = len(task.SubUrls)
    }
    if l1 == 0 {
        return 1
    } else if l1 != 0 && l2 == 0 {
        return l1
    } else if l1 != 0 && l2 != 0 {
        return l2
    }
    return 0
}

/**
kwargs -> Task的字典（k-v）形式
*/
func (task *CrawlerTask) ParseKwargs(kwargs map[string]interface{}) error {
    if siteName, ok := kwargs["siteName"]; !ok {
        return fmt.Errorf("undefined kwarg siteName")
    } else {
        task.SiteName = siteName.(string)
    }
    if url, ok := kwargs["url"]; !ok {
        return fmt.Errorf("undefined kwarg url")
    } else {
        url, err := urllib.Parse(url.(string))
        if err != nil {
            log.Warningf("url is error: %s", url)
            return err
        }
        task.Url = url
    }
    if subUrls, ok := kwargs["subUrls"]; !ok {
        return fmt.Errorf("undefined kwarg subUrls")
    } else {
        task.SubUrls = utils.ToStrings(subUrls.([]interface{}))
    }
    if subUrls2, ok := kwargs["subUrls2"]; !ok {
        return fmt.Errorf("undefined kwarg subUrls2")
    } else {
        task.SubUrls2 = utils.ToStrings(subUrls2.([]interface{}))
    }
    return nil
}

func (task *CrawlerTask) RunTask() (interface{}, error) {
    log.Infof("start running task: %s...", task.SiteName)
    var ret interface{}
    if task.SiteName == "dianping" {
        ret = DianpingCrawler.crawl(task)
    }
    log.Info("task executed")
    return ret, nil
}

func (task *CrawlerTask) ResultFilename() string {
    y, m, d := time.Now().Date()
    return fmt.Sprintf("%s_%d_%d_%d.txt", task.SiteName, y, m, d)
}

func (task *CrawlerTask) urls() []string {
    var reg = regexp.MustCompile("{uint}")
    var total = task.urlsCount()
    var urls = make([]string, total)
    var cnt = 0
    for _, url := range task.SubUrls {
        url = reg.ReplaceAllString(url, fmt.Sprint(rand.Intn(32)))
        if task.SubUrls2 != nil && len(task.SubUrls2) > 0 {
            for _, url2 := range task.SubUrls2 {
                url2 = reg.ReplaceAllString(url2, fmt.Sprint(rand.Intn(32)))
                urls[cnt] = fmt.Sprintf("%s/%s/%s", task.Url, url, url2)
                cnt++
            }
        } else {
            urls[cnt] = fmt.Sprintf("%s/%s", task.Url, url)
            cnt++
        }
    }
    return urls
}
