package main

import (
    "fmt"
    "github.com/Danceiny/go.fastjson"
    "github.com/Danceiny/gocelery"
    "github.com/gocolly/redisstorage"
)

var (
    DianpingCrawler *Dianping
    storage         *redisstorage.Storage
)

func initWorker() {
    DianpingCrawler = &Dianping{fastjson.GetEnvOrDefault("DIANPING_MAX_PAGE_NUMBER", 50).(int)}
    storage = &redisstorage.Storage{
        Address:  fmt.Sprintf("%s:%d", CELERY_BACKEND_HOST, CELERY_BACKEND_PORT),
        Password: CELERY_BACKEND_PASSWORD,
        DB:       0,
        Prefix:   "creeper",
    }
}

func StartWorker() {
    initWorker()
    // create broker and backend
    celeryBroker := gocelery.NewRedisCeleryBroker(CELERY_BROKER_HOST, CELERY_BROKER_PORT, 0, CELERY_BROKER_PASSWORD)
    celeryBackend := gocelery.NewRedisCeleryBackend(CELERY_BACKEND_HOST, CELERY_BACKEND_PORT, 0, CELERY_BACKEND_PASSWORD)

    // Configure with 2 celery workers
    celeryServer, _ := gocelery.NewCeleryServer(celeryBroker, celeryBackend, 10)

    // worker.add name reflects "add" task method found in "worker.py"
    // this worker uses args
    celeryServer.Register("runTask", &CrawlerTask{})
    celeryServer.StartWorker()
}
