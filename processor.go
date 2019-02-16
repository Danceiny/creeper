package main

import (
	"fmt"
	"github.com/Danceiny/gocelery"
	log "github.com/sirupsen/logrus"
	urllib "net/url"
)

func interfaces2strings(input []interface{}, elementIsString bool) []string {
	output := make([]string, len(input))
	if elementIsString {
		for i, o := range input {
			output[i] = o.(string)
		}
	} else {
		for i, o := range input {
			output[i] = fmt.Sprint(o)
		}
	}
	return output
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
		task.SubUrls = interfaces2strings(subUrls.([]interface{}), true)
	}
	if subUrls2, ok := kwargs["subUrls2"]; !ok {
		return fmt.Errorf("undefined kwarg subUrls2")
	} else {
		task.SubUrls2 = interfaces2strings(subUrls2.([]interface{}), true)
	}
	return nil
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
