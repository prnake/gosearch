package site

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var BaiduP = EndPoint{
	// =================BAIDU==================
	//Cookie = ""
	Cookie: "BIDUPSID=3535F2B8A915447A4839A7DD194BA7B3; PSTM=1617718622; __yjs_duid=1_76a86308fbc1b29f8b2c852e8c9e24fd1620886035244; BD_UPN=123253; MCITY=-%3A; BDUSS=m5TRVBaQjk0NXF0SFpmMHZtN2ZWSTlyZVo2d0ZJTVhlU1BMTUpqbXRrZEl0WFZqSUFBQUFBJCQAAAAAAAAAAAEAAAC8QhUx1MbUqtauwfpoaQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEgoTmNIKE5jc; BDUSS_BFESS=m5TRVBaQjk0NXF0SFpmMHZtN2ZWSTlyZVo2d0ZJTVhlU1BMTUpqbXRrZEl0WFZqSUFBQUFBJCQAAAAAAAAAAAEAAAC8QhUx1MbUqtauwfpoaQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEgoTmNIKE5jc; BDSFRCVID=rY_OJeC62lCrte6jotU8bVRNE2SdnBRTH6aotxm4whxuChbLecyMEG0Pyf8g0KubiKd_ogKK0eOTHktF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF=JJkO_D_atKvDqTrP-trf5DCShUFsyMvlB2Q-XPoO3KJnMU74M4Qtj4Fp3-ON--QiW5cpoMbgylRp8P3y0bb2DUA1y4vp5MnqQeTxoUJ2fnRJEUcGqj5Ah--ebPRiJPQ9QgbWLpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0MC09j6KhDTPVKgTa54cbb4o2WbCQQq3O8pcN2b5oQT81jnnHatoh32jZ-pb4bb5vOPQKDpOUWfAkXpJvQnJjt2JxaqRCKhv-Sl5jDh3Me-AsLn6te6jzaIvy0hvctn5cShnc5MjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDNKtt5_jJJIsBP_8aJ7bHn7gbJK_-P4DeP6wexRZ5mAqoDQ6tbbnHRcH3noJ-lKw3MTA-55lba6naIQqa-cKDR38Kp58hPT00tb8Qnb43bRT2PPy5KJvfJo1BTrYhP-UyN3-Wh3725nlMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMonLafDD3fb7kbP6Eq4D_MfOtetJyaR3OoM7vWJ5WqR7jD5DbMP4qQfcrQxTn3Nrfb-OIJPTKShbXKxonKlLObbOGLM6EfNcZ0l8K3l02V-bRDDcfQJQDQt7JKPRMW20j0l7mWnvMsxA45J7cM4IseboJLfT-0bc4KKJxbnLWeIJEjjCKejbQDGADq6nfb5kXLn7J-J5HfbTkbITjhPrMKRrdWMT-0bFHWpO-bnK-jR7m3j7D3lDNQH5gKP6GLHn7_JjObPnVM-3d2boTM-DTbfb8WxQxtNR80DnjtpvhHRTobhnobUPUDMo9LUvWbgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj0DKLK-oj-D8lej7P; BAIDUID=9A9CC424161261B757A36C59CFED76E5:FG=1; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=36542_37557_37519_36884_37627_36786_37539_37500_26350_37343_37461; BAIDUID_BFESS=9A9CC424161261B757A36C59CFED76E5:FG=1; BDSFRCVID_BFESS=rY_OJeC62lCrte6jotU8bVRNE2SdnBRTH6aotxm4whxuChbLecyMEG0Pyf8g0KubiKd_ogKK0eOTHktF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF_BFESS=JJkO_D_atKvDqTrP-trf5DCShUFsyMvlB2Q-XPoO3KJnMU74M4Qtj4Fp3-ON--QiW5cpoMbgylRp8P3y0bb2DUA1y4vp5MnqQeTxoUJ2fnRJEUcGqj5Ah--ebPRiJPQ9QgbWLpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0MC09j6KhDTPVKgTa54cbb4o2WbCQQq3O8pcN2b5oQT81jnnHatoh32jZ-pb4bb5vOPQKDpOUWfAkXpJvQnJjt2JxaqRCKhv-Sl5jDh3Me-AsLn6te6jzaIvy0hvctn5cShnc5MjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDNKtt5_jJJIsBP_8aJ7bHn7gbJK_-P4DeP6wexRZ5mAqoDQ6tbbnHRcH3noJ-lKw3MTA-55lba6naIQqa-cKDR38Kp58hPT00tb8Qnb43bRT2PPy5KJvfJo1BTrYhP-UyN3-Wh3725nlMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMonLafDD3fb7kbP6Eq4D_MfOtetJyaR3OoM7vWJ5WqR7jD5DbMP4qQfcrQxTn3Nrfb-OIJPTKShbXKxonKlLObbOGLM6EfNcZ0l8K3l02V-bRDDcfQJQDQt7JKPRMW20j0l7mWnvMsxA45J7cM4IseboJLfT-0bc4KKJxbnLWeIJEjjCKejbQDGADq6nfb5kXLn7J-J5HfbTkbITjhPrMKRrdWMT-0bFHWpO-bnK-jR7m3j7D3lDNQH5gKP6GLHn7_JjObPnVM-3d2boTM-DTbfb8WxQxtNR80DnjtpvhHRTobhnobUPUDMo9LUvWbgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj0DKLK-oj-D8lej7P; delPer=0; BD_CK_SAM=1; PSINO=3; BA_HECTOR=8581ala48k0la181aga50fm81hm8prs1e; ZFY=MW5tgLZfxnXk8HLzWL5vKMTfHTBgdsBmZudM7p1:BE:Aw:C; sugstore=1; H_PS_645EC=6985pYm8uvzjRIBPLtrwALaCau4%2Fxj7uefaYB7sGfR5qYLEGfRZ8Eyb%2F1oc",
	Url:    "https://www.baidu.com",
	Domain: "www.baidu.com",
	Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	Search: "https://www.baidu.com" + "/s?wd=%s" + "&usm=3&rsv_idx=2&rsv_page=1",
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
	return fmt.Sprintf(BaiduP.Search, baidu.Req.Q)
}

func (baidu *Baidu) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 0, List: []Entity{}}

	if baidu.resp.doc != nil {
		// Find the review items
		//log.Printf("Review doc: %s\n", resp.doc.Text())
		baidu.resp.doc.Find("div[srcid]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("h3").Find("a").Text()
			url := s.AttrOr("mu", "")
			tpl := s.AttrOr("tpl", "")
			if tpl != "se_com_default" {
				return
			}
			subTitle := s.Find(".c-gap-top-small").Find("span").Text()
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
	//提交请求
	request, err := http.NewRequest("GET", baidu.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", BaiduP.Domain)
	request.Header.Add("Cookie", BaiduP.Cookie)
	request.Header.Add("Accept", BaiduP.Accept)

	return SendDo(client, request)

}
