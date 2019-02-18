package main

import (
    "fmt"
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

func (task *CrawlerTask) ResultFilename() string {
    y, m, d := time.Now().Date()
    return fmt.Sprintf("%s_%d_%d_%d.txt", task.SiteName, y, m, d)
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

func (task *CrawlerTask) urls() []string {
    var reg = regexp.MustCompile("{uint}")
    var total = task.urlsCount()
    var urls = make([]string, total)
    var cnt = 0
    for _, url := range task.SubUrls {
        reg.ReplaceAllString(url, fmt.Sprint(rand.Int()))
        if task.SubUrls2 != nil && len(task.SubUrls2) > 0 {
            for _, url2 := range task.SubUrls2 {
                reg.ReplaceAllString(url2, fmt.Sprint(rand.Int()))
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
func (task *CrawlerTask) RunTask() (interface{}, error) {
    log.Infof("start running task: %s...", task.SiteName)
    var ret interface{}
    if task.SiteName == "dianping" {
        ret = DianpingCrawler.crawl(task)
    }
    log.Info("task executed")
    return ret, nil
}
