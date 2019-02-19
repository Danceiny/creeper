package main

import (
    "github.com/Danceiny/go.fastjson"
    "reflect"
    "time"
)

const CAPITALS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type TStruct interface {
    GetId() string
}

type Shop struct {
    Id        string               `json:"id"`
    Title     string               `json:"title"`
    Contacts  []string             `json:"contacts"`
    Url       string               `json:"url"`
    CrawledAt *time.Time           `json:"crawled_at"`
    Images    []string             `json:"images"`
    Attr      *fastjson.JSONObject `json:"attr"`
}

func (s *Shop) GetId() string {
    return s.Id
}

func (s Shop) getShopFields() []string {
    t := reflect.TypeOf(s)
    c := t.NumField()
    var ret = make([]string, c)
    for i := 0; i < c; i++ {
        ret[i] = t.Field(i).Tag.Get("json")
    }
    return ret
}
