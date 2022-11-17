package site

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

const (
	Debug     = false
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"

	// =================BAIDU==================
	BaiduCookie = "BIDUPSID=3535F2B8A915447A4839A7DD194BA7B3; PSTM=1617718622; __yjs_duid=1_76a86308fbc1b29f8b2c852e8c9e24fd1620886035244; BD_UPN=123253; MCITY=-%3A; BDUSS=m5TRVBaQjk0NXF0SFpmMHZtN2ZWSTlyZVo2d0ZJTVhlU1BMTUpqbXRrZEl0WFZqSUFBQUFBJCQAAAAAAAAAAAEAAAC8QhUx1MbUqtauwfpoaQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEgoTmNIKE5jc; BDUSS_BFESS=m5TRVBaQjk0NXF0SFpmMHZtN2ZWSTlyZVo2d0ZJTVhlU1BMTUpqbXRrZEl0WFZqSUFBQUFBJCQAAAAAAAAAAAEAAAC8QhUx1MbUqtauwfpoaQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEgoTmNIKE5jc; BDSFRCVID=rY_OJeC62lCrte6jotU8bVRNE2SdnBRTH6aotxm4whxuChbLecyMEG0Pyf8g0KubiKd_ogKK0eOTHktF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF=JJkO_D_atKvDqTrP-trf5DCShUFsyMvlB2Q-XPoO3KJnMU74M4Qtj4Fp3-ON--QiW5cpoMbgylRp8P3y0bb2DUA1y4vp5MnqQeTxoUJ2fnRJEUcGqj5Ah--ebPRiJPQ9QgbWLpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0MC09j6KhDTPVKgTa54cbb4o2WbCQQq3O8pcN2b5oQT81jnnHatoh32jZ-pb4bb5vOPQKDpOUWfAkXpJvQnJjt2JxaqRCKhv-Sl5jDh3Me-AsLn6te6jzaIvy0hvctn5cShnc5MjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDNKtt5_jJJIsBP_8aJ7bHn7gbJK_-P4DeP6wexRZ5mAqoDQ6tbbnHRcH3noJ-lKw3MTA-55lba6naIQqa-cKDR38Kp58hPT00tb8Qnb43bRT2PPy5KJvfJo1BTrYhP-UyN3-Wh3725nlMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMonLafDD3fb7kbP6Eq4D_MfOtetJyaR3OoM7vWJ5WqR7jD5DbMP4qQfcrQxTn3Nrfb-OIJPTKShbXKxonKlLObbOGLM6EfNcZ0l8K3l02V-bRDDcfQJQDQt7JKPRMW20j0l7mWnvMsxA45J7cM4IseboJLfT-0bc4KKJxbnLWeIJEjjCKejbQDGADq6nfb5kXLn7J-J5HfbTkbITjhPrMKRrdWMT-0bFHWpO-bnK-jR7m3j7D3lDNQH5gKP6GLHn7_JjObPnVM-3d2boTM-DTbfb8WxQxtNR80DnjtpvhHRTobhnobUPUDMo9LUvWbgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj0DKLK-oj-D8lej7P; BAIDUID=9A9CC424161261B757A36C59CFED76E5:FG=1; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=36542_37557_37519_36884_37627_36786_37539_37500_26350_37343_37461; BAIDUID_BFESS=9A9CC424161261B757A36C59CFED76E5:FG=1; BDSFRCVID_BFESS=rY_OJeC62lCrte6jotU8bVRNE2SdnBRTH6aotxm4whxuChbLecyMEG0Pyf8g0KubiKd_ogKK0eOTHktF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF_BFESS=JJkO_D_atKvDqTrP-trf5DCShUFsyMvlB2Q-XPoO3KJnMU74M4Qtj4Fp3-ON--QiW5cpoMbgylRp8P3y0bb2DUA1y4vp5MnqQeTxoUJ2fnRJEUcGqj5Ah--ebPRiJPQ9QgbWLpQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0MC09j6KhDTPVKgTa54cbb4o2WbCQQq3O8pcN2b5oQT81jnnHatoh32jZ-pb4bb5vOPQKDpOUWfAkXpJvQnJjt2JxaqRCKhv-Sl5jDh3Me-AsLn6te6jzaIvy0hvctn5cShnc5MjrDRLbXU6BK5vPbNcZ0l8K3l02V-bIe-t2XjQhDNKtt5_jJJIsBP_8aJ7bHn7gbJK_-P4DeP6wexRZ5mAqoDQ6tbbnHRcH3noJ-lKw3MTA-55lba6naIQqa-cKDR38Kp58hPT00tb8Qnb43bRT2PPy5KJvfJo1BTrYhP-UyN3-Wh3725nlMKoaMp78jR093JO4y4Ldj4oxJpOJ5JbMonLafDD3fb7kbP6Eq4D_MfOtetJyaR3OoM7vWJ5WqR7jD5DbMP4qQfcrQxTn3Nrfb-OIJPTKShbXKxonKlLObbOGLM6EfNcZ0l8K3l02V-bRDDcfQJQDQt7JKPRMW20j0l7mWnvMsxA45J7cM4IseboJLfT-0bc4KKJxbnLWeIJEjjCKejbQDGADq6nfb5kXLn7J-J5HfbTkbITjhPrMKRrdWMT-0bFHWpO-bnK-jR7m3j7D3lDNQH5gKP6GLHn7_JjObPnVM-3d2boTM-DTbfb8WxQxtNR80DnjtpvhHRTobhnobUPUDMo9LUvWbgcdot5yBbc8eIna5hjkbfJBQttjQn3hfIkj0DKLK-oj-D8lej7P; delPer=0; BD_CK_SAM=1; PSINO=3; BA_HECTOR=8581ala48k0la181aga50fm81hm8prs1e; ZFY=MW5tgLZfxnXk8HLzWL5vKMTfHTBgdsBmZudM7p1:BE:Aw:C; sugstore=1; H_PS_645EC=6985pYm8uvzjRIBPLtrwALaCau4%2Fxj7uefaYB7sGfR5qYLEGfRZ8Eyb%2F1oc"
	BaiduUrl    = "https://www.baidu.com"
	BaiduDomain = "www.baidu.com"
	BaiduAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	BaiduSearch = BaiduUrl + "/s?wd=%s" + "&usm=3&rsv_idx=2&rsv_page=1"
	BaiduFrom   = "百度"

	// =================BING==================
	BingCoolkie = "MUID=16AE9B2FE9DC6568044F8B3EE8F2640C; MUIDB=16AE9B2FE9DC6568044F8B3EE8F2640C; SRCHD=AF=BDVEHC; SRCHUID=V=2&GUID=7321ED22410B41459C503CB6D2628196&dmnchg=1; _UR=QS=0&TQS=0; imgv=flts=20220704; _tarLang=default=zh-Hans; _TTSS_IN=hist=WyJlbiIsImF1dG8tZGV0ZWN0Il0=; _TTSS_OUT=hist=WyJ6aC1IYW5zIl0=; _HPVN=CS=eyJQbiI6eyJDbiI6MjQsIlN0IjoyLCJRcyI6MCwiUHJvZCI6IlAifSwiU2MiOnsiQ24iOjI0LCJTdCI6MCwiUXMiOjAsIlByb2QiOiJIIn0sIlF6Ijp7IkNuIjoyNCwiU3QiOjEsIlFzIjowLCJQcm9kIjoiVCJ9LCJBcCI6dHJ1ZSwiTXV0ZSI6dHJ1ZSwiTGFkIjoiMjAyMi0wOS0wNlQwMDowMDowMFoiLCJJb3RkIjowLCJHd2IiOjAsIkRmdCI6bnVsbCwiTXZzIjowLCJGbHQiOjAsIkltcCI6Nzd9; ANIMIA=FRE=1; MMCASM=ID=3FEB6F3855CC49E584F2DA61F6E5E44C; ZHCHATSTRONGATTRACT=TRUE; _SS=SID=11649FAAD94F6D9716388DF6D8296CB4&PC=U316; SRCHS=PC=U316; ABDEF=V=13&ABDV=13&MRB=0&MRNB=1668441959593; SUID=M; _EDGE_S=SID=11649FAAD94F6D9716388DF6D8296CB4&ui=zh-cn; SRCHUSR=DOB=20210406&T=1668472757000&TPC=1668472758000; ZHCHATWEAKATTRACT=TRUE; ipv6=hit=1668476360196&t=4; ZHLASTACTIVECHAT=0; ZHSEARCHCHATSTATUS=STATUS=0; SNRHOP=I=&TS=; RECSEARCH=SQs=[{\"q\":\"giac%20%E4%B8%8A%E6%B5%B7\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"sessioncachesize\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"rset\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"xx%3Ainitiatingheapoccupancypercent\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"%E7%AE%A1%E7%90%86%E7%9A%84%E5%B8%B8%E8%AF%86\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"%E7%AE%A1%E7%90%86%E7%9A%84%E5%B8%B8%E8%AF%86%20%E5%BE%B7%E9%B2%81%E5%85%8B\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"%E7%AE%A1%E7%90%86\"%2C\"c\":1%2C\"ad\":true}%2C{\"q\":\"goquery%20%E5%BE%AA%E7%8E%AF\"%2C\"c\":1%2C\"ad\":false}%2C{\"q\":\"yuanbiguo\"%2C\"c\":1%2C\"ad\":false}]; SRCHHPGUSR=SRCHLANGV2=zh-Hans&BRW=W&BRH=S&CW=1396&CH=435&DPR=2&UTC=480&DM=0&WTS=63804069557&HV=1668474094&BZA=0&SRCHLANG=zh-Hans&SW=1440&SH=900&PV=11.2.3&EXLTT=6&SCW=1381&SCH=1408&PRVCW=1396&PRVCH=764"
	BingUrl     = "https://cn.bing.com"
	BingDomain  = "cn.bing.com"
	BingAccept  = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	BingSearch  = BingUrl + "/search?q=%s" + "&PC=U316&FORM=CHROMN"
	BingFrom    = "Bing"
)

type SearchEngine interface {
	Search() (result *EntityList)
	urlWrap() (url string)
	toEntityList() (entityList *EntityList)
	send() (resp *Resp, err error)
}

type Baidu struct {
	Req  Req
	resp Resp
}

type Bing struct {
	Req  Req
	resp Resp
}

type JsonResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data *EntityList `json:"data"`
}

type Req struct {
	Q         string
	url       string
	userAgent string
	http.Cookie
}

type Resp struct {
	code int
	body string
	doc  *goquery.Document
}

type EntityList struct {
	Index int      `json:"index"`
	Size  int      `json:"size"`
	List  []Entity `json:"list"`
}

type Entity struct {
	Title    string `json:"title"`
	Host     string `json:"host"`
	Url      string `json:"url"`
	SubTitle string `json:"subTitle"`
	From     string `json:"from"`
}

func init() {
	fmt.Printf("site init,Debug:%v, UserAgent:%v\n", Debug, UserAgent)
}

func LoadConf() {
	fmt.Println("load site conf")
}
