package site

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	// "encoding/json"
	"github.com/PuerkitoBio/goquery"
)

var BaiduP = EndPoint{
	// =================BAIDU==================
	// This Cookie is the default setting of Baidu, do not change it.
	Cookie: "H_PS_PSSID=39938_39999_40040_40089; BA_HECTOR=008h84a4002g2gal0k018k0khre4s71ip7tvp1s; BAIDUID=D2559CEF3EB85352E8809DA2B1336662:FG=1; BIDUPSID=D2559CEF3EB85352B585025DD84EFAF1; PSTM=1704195774; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BD_UPN=12314753;",
	Url:    "https://www.baidu.com",
	Domain: "www.baidu.com",
	Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	Search: "https://www.baidu.com" + "/s?wd=%s&pn=%d" + "&usm=3&rsv_idx=2&rsv_page=1",
	From:   "百度",

	UseProxy:  false,
	Proxy:     url.URL{},
	Transport: GetTransport(),
}

func (baidu *Baidu) Enable() (enable bool) {
	return GetEnable(BaiduP.Domain)
}

func (baidu *Baidu) Search() (result *EntityList) {
	baidu.Req.url = baidu.urlWrap()
	log.Printf("baidu req.url: %s\n", baidu.Req.url)
	resp := &Resp{}
	resp, _ = baidu.send()
	baidu.resp = *resp
	result = baidu.toEntityList()
	return result
}

func (baidu *Baidu) urlWrap() (url string) {
	return fmt.Sprintf(BaiduP.Search, baidu.Req.Q, baidu.Req.Page * 10)
}

func (baidu *Baidu) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 0, List: []Entity{}}

	if baidu.resp.doc != nil {
		// Find the review items
		//log.Printf("Review doc: %s\n", resp.doc.Text())
		baidu.resp.doc.Find("div[srcid]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("h3").Find("a").Remove().Text()
			url := s.AttrOr("mu", "")
			tpl := s.AttrOr("tpl", "")
			if tpl != "se_com_default" {
				return
			}
			subTitle := s.Text()
			entity := Entity{From: BaiduP.From}
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

func (baidu *Baidu) send() (resp *Resp, err error) {

    client := &http.Client{
        Transport: &BaiduP.Transport,
    }

    request, err := http.NewRequest("GET", baidu.urlWrap(), nil)
    if err != nil {
        log.Println("Error creating Baidu request:", err)
        return nil, err
    }

    // Add headers to your Baidu request
    request.Header.Add("User-Agent", UserAgent)
    request.Header.Add("Host", BaiduP.Domain)
    request.Header.Add("Cookie", BaiduP.Cookie)
    request.Header.Add("Accept", BaiduP.Accept)

    return SendDo(client, request)
}