package site

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var GoogleP = EndPoint{
	// =================Google==================
	//Cookie : "",
	Cookie: "CONSENT=YES+srp.gws-20220523-0-RC1.zh-CN+FX; HSID=AIDFOIfZRXMjhRznJ; SSID=ATYvcUlmXr3mPFerl; APISID=fTLh3HcYiZa0Ch2l/AiukySr7MDg_GhSRo; SAPISID=-JpKVkIJZXDgucyp/AjPiMHGTZYbtqrQbt; __Secure-1PAPISID=-JpKVkIJZXDgucyp/AjPiMHGTZYbtqrQbt; __Secure-3PAPISID=-JpKVkIJZXDgucyp/AjPiMHGTZYbtqrQbt; SEARCH_SAMESITE=CgQIk5YB; SID=Qggrdf-HuR3Im3foqYJWVZDPEay8b5U0-O_E6L489Ppxqy2bY358yp5YaaIs5NjEgyW23g.; __Secure-1PSID=Qggrdf-HuR3Im3foqYJWVZDPEay8b5U0-O_E6L489Ppxqy2brNnvCnxg8UCPQbHqZClMmA.; __Secure-3PSID=Qggrdf-HuR3Im3foqYJWVZDPEay8b5U0-O_E6L489Ppxqy2bJ0s1kXYhOOAgur0ax0qp6w.; OTZ=6777442_24_24__24_; AEC=AakniGNW8IBpZgW_-IGOh45Otu8fCK-Ty1n0eSZhgs05euV7s2hhmWLhCQ; NID=511=k96f3xX7KMB1Wvhg478KURvWSAd53Y3rTkVG3bMb1FhtdJwQUbUiHJOnVTgFgW0Mad_1X5gKnWLMt0lPHl33nQdVCTzEiitFOC2dIicYusLP1zl_L6Wh9l-XO6x9VRh4ZxGlcu1bCmlEbBYvz2eL2ioCJ3RMGZtGcr6_1tdpGZH8DRcj3c8X6FxRJqcX5peACa9pGELYmZk4TfOYiUON7p4ht5MTkfA4-hmLUD_JQOT_bK4Aub5SzhXukCFBt5qx2UMkKKtheWlK6cRlR-EIqbC0ppccy2pUdT0YyfE; DV=c77sl5f1j0hXACm1QWzNqJ3oPLjLSRge9TXKlwa5UAEAAMBLUjyomhJKegAAABD63jYjJ8URJwAAAKJBqsVeeJWmDgAAAA; 1P_JAR=2022-11-22-00; SIDCC=AIKkIs0Jk8Fs8wt-yHi9yEM8vFzc6y1bmzVm7mYb3eELJ0t_5yixpCQToLaIDwlOnk3yTvFwKw; __Secure-1PSIDCC=AIKkIs2Sk9zwJs1YfCaHpcFU4ZAKHzDO7Atrzk3b-uXf384Xuc1nQ0U41MiX_VQcEBrJ0TSR1ZY; __Secure-3PSIDCC=AIKkIs1iF4en_Ln8Jonhw8Q79QB_dwbKwcw0khXXIX0vsIwOjBgxmO8dd6kfi10evURBz7YdDg",
	Url:    "https://www.google.com",
	Domain: "www.google.com",
	Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	Search: "https://www.google.com" + "/search?q=%s" + "&ie=UTF-8",
	From:   "Google",

	Transport: GetTransport(),
}

func (g *Google) Enable() (enable bool) {
	return GetEnable(GoogleP.Domain)
}

func (g *Google) Search() (result *EntityList) {
	g.Req.url = g.urlWrap()
	log.Printf("google req.url: %s\n", g.Req.url)
	resp := &Resp{}
	resp, _ = g.send()
	g.resp = *resp
	result = g.toEntityList()
	return result
}

func (g *Google) urlWrap() (url string) {
	return fmt.Sprintf(GoogleP.Search, g.Req.Q)
}

func (g *Google) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}

	if g.resp.doc != nil {
		// Find the review items
		//log.Printf("Review doc: %s\n", g.resp.doc.Text())
		g.resp.doc.Find("div[class=MjjYud]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("div[class=yuRUbf]").Find("h3").Text()
			if title == "" {
				return
			}
			url := s.Find("div[class=yuRUbf]").Find("a").AttrOr("href", "")
			subTitle := s.Find("div[class='Z26q7c UK95Uc']").Find("span").Text()
			entity := Entity{From: GoogleP.From}
			entity.Title = title
			entity.SubTitle = subTitle
			entity.Url = url
			host := strings.ReplaceAll(url, "http://", "")
			host = strings.ReplaceAll(host, "https://", "")
			entity.Host = strings.Split(host, "/")[0]
			entityList.List = append(entityList.List, entity)
		})
	}
	return entityList
}

func (g *Google) send() (resp *Resp, err error) {

	trProxy := &GoogleP.Transport

	client := &http.Client{
		Transport: trProxy,
	}
	//提交请求
	request, err := http.NewRequest("GET", g.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", GoogleP.Domain)
	request.Header.Add("Cookie", GoogleP.Cookie)
	request.Header.Add("Accept", GoogleP.Accept)
	request.Header.Add("authority", "www.google.com")
	//request.Header.Add("accept-encoding", "gzip, deflate, br")
	request.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")

	return SendDo(client, request)

}
