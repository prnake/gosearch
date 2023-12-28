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
	//Cookie: = "MUID=16AE9B2FE9DC6568044F8B3EE8F2640C; MUIDB=16AE9B2FE9DC6568044F8B3EE8F2640C; SRCHD=AF=BDVEHC; SRCHUID=V=2&GUID=7321ED22410B41459C503CB6D2628196&dmnchg=1; _UR=QS=0&TQS=0; imgv=flts=20220704; _tarLang=default=zh-Hans; _TTSS_IN=hist=WyJlbiIsImF1dG8tZGV0ZWN0Il0=; _TTSS_OUT=hist=WyJ6aC1IYW5zIl0=; _HPVN=CS=eyJQbiI6eyJDbiI6MjQsIlN0IjoyLCJRcyI6MCwiUHJvZCI6IlAifSwiU2MiOnsiQ24iOjI0LCJTdCI6MCwiUXMiOjAsIlByb2QiOiJIIn0sIlF6Ijp7IkNuIjoyNCwiU3QiOjEsIlFzIjowLCJQcm9kIjoiVCJ9LCJBcCI6dHJ1ZSwiTXV0ZSI6dHJ1ZSwiTGFkIjoiMjAyMi0wOS0wNlQwMDowMDowMFoiLCJJb3RkIjowLCJHd2IiOjAsIkRmdCI6bnVsbCwiTXZzIjowLCJGbHQiOjAsIkltcCI6Nzd9; ANIMIA=FRE=1; MMCASM=ID=3FEB6F3855CC49E584F2DA61F6E5E44C; ZHCHATSTRONGATTRACT=TRUE; _SS=SID=11649FAAD94F6D9716388DF6D8296CB4&PC=U316; SRCHS=PC=U316; ABDEF=V=13&ABDV=13&MRB=0&MRNB=1668441959593; SUID=M; _EDGE_S=SID=11649FAAD94F6D9716388DF6D8296CB4&ui=zh-cn; SRCHUSR=DOB=20210406&T=1668472757000&TPC=1668472758000; ZHCHATWEAKATTRACT=TRUE; ipv6=hit=1668476360196&t=4; ZHLASTACTIVECHAT=0; ZHSEARCHCHATSTATUS=STATUS=0; SNRHOP=I=&TS=; RECSEARCH=SQs=[{\"q\":\"giac%20%E4%B8%8A%E6%B5%B7\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"sessioncachesize\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"rset\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"xx%3Ainitiatingheapoccupancypercent\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"%E7%AE%A1%E7%90%86%E7%9A%84%E5%B8%B8%E8%AF%86\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"%E7%AE%A1%E7%90%86%E7%9A%84%E5%B8%B8%E8%AF%86%20%E5%BE%B7%E9%B2%81%E5%85%8B\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"%E7%AE%A1%E7%90%86\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"goquery%20%E5%BE%AA%E7%8E%AF\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"yuanbiguo\"%2C\"c\":1%2C\"ad\":false}]; SRCHHPGUSR=SRCHLANGV2=zh-Hans&BRW=W&BRH=S&CW=1396&CH=435&DPR=2&UTC=480&DM=0&WTS=63804069557&HV=1668474094&BZA=0&SRCHLANG=zh-Hans&SW=1440&SH=900&PV=11.2.3&EXLTT=6&SCW=1381&SCH=1408&PRVCW=1396&PRVCH=764"
	Url:    "https://cn.bing.com",
	Domain: "cn.bing.com",
	Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	Search: "https://cn.bing.com" + "/search?q=%s" + "&PC=U316&FORM=CHROMN",
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
	return fmt.Sprintf(BingP.Search, bing.Req.Q)
}

func (bing *Bing) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}

	if bing.resp.doc != nil {
		// Find the review items
		//log.Printf("Review doc: %s\n", resp.doc.Text())
		bing.resp.doc.Find("ol#b_results>li[class=b_algo]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("div[class=b_title]>h2>a").Text()
			url := s.Find("div[class=b_title]>h2>a").AttrOr("href", "")
			subTitle := s.Find("div[class=b_caption]>p").Text()
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
