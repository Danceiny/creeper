package main

import (
	"fmt"
	urllib "net/url"
	"regexp"
)

func main() {
	u, _ := urllib.Parse("http://www.dianping.com")
	task := &CrawlerTask{
		"dianping",
		u,
		[]string{"shanghai"},
		[]string{"education"},
	}
	regstr := fmt.Sprintf("%s%s%s",
		task.Url, Arr2RegOr(task.SubUrls, "/", ""), Arr2RegOr(task.SubUrls2, "/", "$"))
	fmt.Println(regstr)
	var urlReg *regexp.Regexp = regexp.MustCompile(regstr)
	b := urlReg.Match([]byte(`http://www.dianping.com/shanghai/education`))
	fmt.Println(b)
}
