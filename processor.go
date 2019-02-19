package main

import (
    "fmt"
    "github.com/Danceiny/gocelery"
    log "github.com/sirupsen/logrus"
    "math/rand"
    "regexp"
    "time"
)

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
        url = reg.ReplaceAllString(url, fmt.Sprint(rand.Intn(64)))
        if task.SubUrls2 != nil && len(task.SubUrls2) > 0 {
            for _, url2 := range task.SubUrls2 {
                url2 = reg.ReplaceAllString(url2, fmt.Sprint(rand.Intn(64)))
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

func Add(a, b int) int {
    return a + b
}

func StartWorker() {
    // create broker and backend
    celeryBroker := gocelery.NewRedisCeleryBroker(CELERY_BROKER_HOST, CELERY_BROKER_PORT, 0, CELERY_BROKER_PASSWORD)
    celeryBackend := gocelery.NewRedisCeleryBackend(CELERY_BACKEND_HOST, CELERY_BACKEND_PORT, 0, CELERY_BACKEND_PASSWORD)

    // Configure with 2 celery workers
    celeryServer, _ := gocelery.NewCeleryServer(celeryBroker, celeryBackend, 10)

    // worker.add name reflects "add" task method found in "worker.py"
    // this worker uses args
    celeryServer.Register("runTask", &CrawlerTask{})
    celeryServer.Register("add", Add)
    celeryServer.StartWorker()
}

func main() {
    StartWorker()
}
