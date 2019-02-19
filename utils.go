package main

import (
    "context"
    "fmt"
    "github.com/Danceiny/go.fastjson"
    "github.com/gocolly/colly"
    "github.com/levigross/grequests"
    "github.com/sirupsen/logrus"
    "math/rand"
    "net/http"
    "net/url"
    "reflect"
    "sync/atomic"
)

func Arr2RegOr(arr []string, prefix string, suffix string) string {
    if arr == nil || len(arr) == 0 {
        return ""
    }
    var ret = prefix + "("
    var start = true
    for _, s := range arr {
        if start {
            ret += s
        } else {
            start = false
            ret += "|" + s
        }
    }
    return fmt.Sprintf("(%s)%s)", ret, suffix)
}

func GetProxyIp() string {
    resp, _ := grequests.Get(PROXY_SERVER_API+"/", nil)
    return resp.String()
}

func GetProxyUrl() *url.URL {
    return &url.URL{Scheme: "http", Host: GetProxyIp()}
}

func GetAllProxyUrl() []string {
    resp, err := grequests.Get(PROXY_SERVER_API+"_all/", nil)
    if err != nil {
        logrus.Error(err)
        return nil
    }
    addrs := fastjson.ParseArrayT(resp.String(), reflect.String).([]string)
    // for i, addr := range addrs {
    //     addrs[i] = "http://" + addr
    // }
    return addrs
}

// RoundRobinProxySwitcher creates a proxy switcher function which rotates
// ProxyURLs on every request.
// The proxy type is determined by the URL scheme. "http", "https"
// and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
func RoundRobinProxySwitcher(ProxyURLs ...string) (colly.ProxyFunc, error) {
    urls := make([]*url.URL, len(ProxyURLs))
    for i, u := range ProxyURLs {
        parsedU, err := url.Parse(u)
        if err != nil {
            return nil, err
        }
        urls[i] = parsedU
    }
    return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}

type roundRobinSwitcher struct {
    proxyURLs []*url.URL
    index     uint32
}

func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
    u := r.proxyURLs[r.index%uint32(len(r.proxyURLs))]
    atomic.AddUint32(&r.index, 1)
    ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, u.String())
    *pr = *pr.WithContext(ctx)
    return u, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
    b := make([]byte, rand.Intn(10)+10)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

var UserAgents = []string{
    "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
    "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16",
    "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0",
    "Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10"}

func RandomUserAgent() string {
    return UserAgents[rand.Intn(len(UserAgents))] + RandomString()
}
