package site

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v2"
)

type EndPoint struct {
	Cookie string
	Url    string
	Domain string
	Accept string
	Search string
	From   string

	UseProxy bool
	Proxy    url.URL

	Transport http.Transport
}

var (
	UserAgent = "user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
)

var config map[string]interface{}

var endpoints = map[string]*EndPoint{
	BaiduP.From:  &BaiduP,
	WxP.From:     &WxP,
	BingP.From:   &BingP,
	GoogleP.From: &GoogleP,
}

type SearchEngine interface {
	Search() (result *EntityList)
	Enable() (enable bool)
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

type Google struct {
	Req  Req
	resp Resp
}

type Wx struct {
	Req  Req
	resp Resp
}

type JsonResult struct {
	Cost int64       `json:"cost"`
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
	Title         string `json:"title"`
	Host          string `json:"host"`
	Url           string `json:"url"`
	SubTitle      string `json:"subTitle"`
	From          string `json:"from"`
	Score         int    `json:"score"`
	PositionScore int    `json:"positionScore"`
	SearchScore   int    `json:"searchScore"`
	DomainScore   int    `json:"domainScore"`
}

func init() {
	log.Printf("site init, UserAgent:%v\n", UserAgent)
}

func LoadConf() {
	log.Println("load site conf")
	path := "configs/config.yml"
	fi, _ := os.Open(path)
	configData, err := io.ReadAll(fi)
	if err != nil {
		log.Fatal(err)
	}
	config = make(map[string]interface{})

	// 执行解析
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("config: %v\n", config)
	}

	serverConf := config["server"].(map[interface{}]interface{})
	if proxyStr := serverConf["proxy"]; proxyStr != nil {
		for _, e := range endpoints {
			e.setProxy(proxyStr.(string))
		}
	}

	if timeout := serverConf["timeout"]; timeout != nil {
		switch timeout := timeout.(type) {
		case int:
			MaxTimeout = time.Millisecond * time.Duration(timeout)
		case string:
			MaxTimeout, _ = time.ParseDuration(timeout)
		default:
			log.Fatalf("invalid tiemout: %v\n", timeout)
		}
	}

	search := config["search"].([]interface{})
	for _, mo := range search {
		m := mo.(map[interface{}]interface{})
		if proxy := m["proxy"]; proxy != nil {
			endpoints[m["name"].(string)].setProxy(proxy.(string))
		}
	}
}

func (e *EndPoint) setProxy(proxyStr string) {
	proxy, _ := url.Parse(proxyStr)

	e.UseProxy = true
	e.Proxy = *proxy

	e.Transport.Proxy = http.ProxyURL(proxy)
}

func nameToReq(name string, q string) SearchEngine {
	switch name {
	case "百度":
		return &Baidu{Req: Req{Q: q}}
	case "Bing":
		return &Bing{Req: Req{Q: q}}
	case "Google":
		return &Google{Req: Req{Q: q}}
	case "微信公众号":
		return &Wx{Req: Req{Q: q}}
	default:
		return nil
	}
}

func GetAllEnabled(q string) []SearchEngine {
	var enabled []SearchEngine
	for _, e := range endpoints {
		if GetEnable(e.Domain) {
			req := nameToReq(e.From, q)
			if req == nil {
				log.Fatalf("unknown search engine: %v\n", e.From)
			}
			enabled = append(enabled, req)
		}
	}
	return enabled
}

func GetDebug() bool {
	server := config["server"].(map[interface{}]interface{})
	return server["debug"].(bool)
}

func GetPort() int {
	server := config["server"].(map[interface{}]interface{})
	return server["port"].(int)
}

func GetTimeout() time.Duration {
	server := config["server"].(map[interface{}]interface{})
	return time.Duration(server["timeout"].(float32)) * time.Second
}

func GetSearchScore(name string) int {
	search := config["search"].([]interface{})
	for _, mo := range search {
		m := mo.(map[interface{}]interface{})
		if m["name"] == name {
			return m["score"].(int) * m["weight"].(int)
		}
	}
	return 0
}

func GetProxy(domain string) (url.URL, bool) {
	proxy := config["search"].([]interface{})
	for _, mo := range proxy {
		m := mo.(map[string]interface{})
		if m["domain"] == domain {
			url, _ := url.Parse(m["proxy"].(string))
			return *url, true
		}
	}
	return url.URL{}, false
}

func GetDomainScore(host string) int {
	site := config["site"].([]interface{})
	for _, mo := range site {
		m := mo.(map[interface{}]interface{})
		if m["domain"] == host {
			return m["score"].(int) * m["weight"].(int)
		}
	}
	return 0
}
func GetPositionWeight(domain string) int {
	search := config["search"].([]interface{})
	for _, mo := range search {
		m := mo.(map[interface{}]interface{})
		if m["name"] == domain {
			return m["positionWeight"].(int)
		}
	}
	return 1
}
func GetEnable(domain string) bool {
	search := config["search"].([]interface{})
	for _, mo := range search {
		m := mo.(map[interface{}]interface{})
		if m["domain"] == domain {
			return m["enable"].(bool)
		}
	}
	return false
}

func SendDo(client *http.Client, request *http.Request) (*Resp, error) {
	resp := &Resp{code: 200}

	//处理返回结果
	response, e := client.Do(request)
	if response == nil {
		resp.code = 500
		log.Printf("response nil: %v\n", e)
		return resp, nil
	}
	if response.StatusCode != 200 {
		resp.code = response.StatusCode
		log.Printf("status code error: %d %s\n", response.StatusCode, response.Status)
		return resp, nil
	}
	defer response.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println(err)
	}

	resp.code = 200
	resp.doc = doc

	return resp, nil
}
