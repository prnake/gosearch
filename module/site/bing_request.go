package site

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"crypto/rand"
    "encoding/hex"
	"github.com/PuerkitoBio/goquery"
)

var BingP = EndPoint{
	// =================BING==================
	Cookie: "",
	// Cookie: "MUID=07A703D7A9ED678C3DB21061A8346688;",
	Url:    "https://cn.bing.com",
	Domain: "cn.bing.com",
	Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	Search: "https://cn.bing.com" + "/search?q=%s&first=%d&filters=%s" + "&PC=U316&FORM=CHROMN",
	From:   "Bing",

	Transport: GetTransport(),
}

func (bing *Bing) Enable() (enable bool) {
	return GetEnable(BingP.Domain)
}

func (bing *Bing) Search() (result *EntityList) {
    retryCount := 3
    resultChan := make(chan *EntityList, retryCount)
    var wg sync.WaitGroup

    for retry := 0; retry < retryCount; retry++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            bing.Req.url = bing.urlWrap()
            log.Printf("bing req.url: %s\n", bing.Req.url)
            resp := &Resp{}
            resp, _ = bing.send()
            bing.resp = *resp
            result := bing.toEntityList()

            if result.Size > 0 {
                resultChan <- result
            }
        }()
    }

    go func() {
        wg.Wait()
        close(resultChan)
    }()

    for result = range resultChan {
        return result
    }

    log.Printf("Failed to get non-empty result after %d retries", retryCount)
    return result
}

func (bing *Bing) urlWrap() (url string) {
	return fmt.Sprintf(BingP.Search, bing.Req.Q, bing.Req.Page*10, bing.Req.Filter)
}

func (bing *Bing) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 0, List: []Entity{}}

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

func generateRandomHex(length int) (string, error) {
    byteLength := length / 2

    randomBytes := make([]byte, byteLength)

    _, err := rand.Read(randomBytes)
    if err != nil {
        return "", err
    }

    hexString := hex.EncodeToString(randomBytes)
	hexString = strings.ToUpper(hexString)

    return hexString, nil
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
	// request.Header.Add("Cookie", BingP.Cookie)
	randomHex, err := generateRandomHex(32)
	if err == nil {
		request.Header.Add("Cookie", "MUID="+randomHex+";")
	}

	request.Header.Add("Accept", BingP.Accept)

	return SendDo(client, request)
}
