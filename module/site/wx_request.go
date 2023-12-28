package site

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var WxP = EndPoint{
	// =================Wx==================
	Cookie: "",
	//Cookie : "SUID=F84DCC781539960A000000006235D4E8; SUV=1647695080857391; ssuid=6792259580; weixinIndexVisited=1; IPLOC=CN3100; ABTEST=0|1668775030|v1; JSESSIONID=aaa2PBWFjHNS8hQ8_tfpy; cd=1668913493&0f942166ea05ede01cfe88195d36508d; rd=tyllllllll20WBOSYTuBqQ2iuqV0WBOqAfJbLZllll9llllxVllll5@@@@@@@@@@; ld=6Zllllllll20WBOSYTuBqQ2DdNH0WBOqAfJbLZllll9lllllVklll5@@@@@@@@@@; LSTMV=217%2C66; LCLKINT=1482; PHPSESSID=udut3pen5cml0b9849o68jch40; SNUID=9D60E0542D28C24D1776D6A42DF58CE7; ariaDefaultTheme=undefined",
	Url:    "https://weixin.sogou.com",
	Domain: "weixin.sogou.com",
	Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	Search: "https://weixin.sogou.com" + "/weixin?type=2&s_from=input&query=%s&ie=utf8&_sug_=n&_sug_type_=",
	From:   "微信公众号",
	
	Transport: GetTransport(),
}

func (wx *Wx) Enable() (enable bool) {
	return GetEnable(WxP.Domain)
}

func (wx *Wx) Search() (result *EntityList) {
	wx.Req.url = wx.urlWrap()
	log.Printf("wx req.url: %s\n", wx.Req.url)
	resp := &Resp{}
	resp, _ = wx.send()
	wx.resp = *resp
	result = wx.toEntityList()
	return result
}

func (wx *Wx) urlWrap() (url string) {
	return fmt.Sprintf(WxP.Search, wx.Req.Q)
}

func (wx *Wx) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 0}
	entityList.List = []Entity{}

	if wx.resp.doc != nil {
		// Find the review items
		//log.Printf("Wx Review doc: %s\n", wx.resp.doc.Text())
		wx.resp.doc.Find("div[class='txt-box']").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("h3").Find("a").Text()
			url := s.Find("h3").Find("a").AttrOr("href", "")
			url = WxP.Url + url
			subTitle := s.Find("p[class='txt-info']").Text()

			entity := Entity{From: WxP.From}
			entity.Title = title
			entity.SubTitle = subTitle
			entity.Url = url
			host := s.Find("div[class='s-p']").Find("a").Text()
			entity.Host = strings.Split(host, "/")[0]
			entityList.List = append(entityList.List, entity)
		})
		entityList.Size = len(entityList.List)
	}
	return entityList
}

func (wx *Wx) send() (resp *Resp, err error) {

	client := &http.Client{
		Transport: &WxP.Transport,
	}
	//提交请求
	request, err := http.NewRequest("GET", wx.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", WxP.Domain)
	request.Header.Add("Cookie", WxP.Cookie)
	request.Header.Add("Accept", WxP.Accept)

	return SendDo(client, request)

}
