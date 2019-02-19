package main

import (
    "fmt"
    utils "github.com/Danceiny/go.utils"
    log "github.com/sirupsen/logrus"
    urllib "net/url"
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
