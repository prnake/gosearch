package site

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var BingP = EndPoint{
	// =================BING==================
	Cookie: "",
	Url:    "https://cn.bing.com",
	Domain: "cn.bing.com",
	Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	Search: "https://cn.bing.com" + "/search?q=%s&first=%d" + "&PC=U316&FORM=CHROMN",
	From:   "Bing",

	Transport: GetTransport(),
}

func (bing *Bing) Enable() (enable bool) {
	return GetEnable(BingP.Domain)
}

func (bing *Bing) Search() (result *EntityList) {
	bing.Req.url = bing.urlWrap()
	log.Printf("bing req.url: %s\n", bing.Req.url)
	resp := &Resp{}
	resp, _ = bing.send()
	bing.resp = *resp
	result = bing.toEntityList()
	return result
}

func (bing *Bing) urlWrap() (url string) {
	return fmt.Sprintf(BingP.Search, bing.Req.Q, bing.Req.Page*10)
}

func (bing *Bing) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}

	if bing.resp.doc != nil {
		// Find the review items
		bing.resp.doc.Find("ol#b_results>li[class*=b_algo]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("h2").Remove().Text()
			url := s.Find("a").Remove().AttrOr("href", "")
			subTitle := s.Text()
			entity := Entity{From: BingP.From}
			entity.Title = title
			entity.SubTitle = subTitle
			entity.Url = url
			host := strings.ReplaceAll(url, "http://", "")
			host = strings.ReplaceAll(host, "https://", "")
			entity.Host = strings.Split(host, "/")[0]
			entityList.List = append(entityList.List, entity)
		})
		entityList.Size = len(entityList.List)
	}
	return entityList
}

func (bing *Bing) send() (resp *Resp, err error) {

	client := &http.Client{
		Transport: &BingP.Transport,
	}
	//提交请求
	request, err := http.NewRequest("GET", bing.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", BingP.Domain)
	request.Header.Add("Cookie", BingP.Cookie)
	request.Header.Add("Accept", BingP.Accept)

	return SendDo(client, request)
}
