package main

import (
	"encoding/json"
	"fmt"
	"gosearch/module/site"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"strconv"
)

// Helper function to check if a string is in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

type EntityWrapper struct { //注意此处
	entity []site.Entity
	by     func(p, q *site.Entity) bool
}

func (ew EntityWrapper) Len() int { // 重写 Len() 方法
	return len(ew.entity)
}
func (ew EntityWrapper) Swap(i, j int) { // 重写 Swap() 方法
	ew.entity[i], ew.entity[j] = ew.entity[j], ew.entity[i]
}
func (ew EntityWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return ew.by(&ew.entity[i], &ew.entity[j])
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/health", health)
	http.HandleFunc("/search", search)
	site.LoadConf()
	go func() {
		log.Println("go in")
		defer func() {
			if err := recover(); err != nil {
				log.Println("go err:", err)
			}
		}()
		log.Println("go out")
	}()
	//handle定义请求访问该服务器里的/health路径，就有下面health去处理，health一般为健康检查
	err := http.ListenAndServe(fmt.Sprintf(":%v", site.GetPort()), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// 定义handle处理函数，只要该health被调用，就会写入ok
func health(w http.ResponseWriter, request *http.Request) {
	log.Println(request.URL)
	_ = request.ParseForm()
	log.Println(request.Form.Get("user"))
	_, _ = io.WriteString(w, "ok")
}

func search(w http.ResponseWriter, request *http.Request) {
	log.Println(request.URL)
	_ = request.ParseForm()

    authHeader := request.Header.Get("Authorization")
    if len(site.AuthHeader) > 0 && authHeader != site.AuthHeader {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	q := request.Form.Get("q")
	q = url.QueryEscape(q)
	filter := request.Form.Get("filter")
	filter = url.QueryEscape(filter)
	page, err := strconv.Atoi(request.Form.Get("page"))
	if err != nil {
		page = 0
	}

	var engine []string

	if engineStr := request.Form.Get("engine"); len(engineStr) == 0 {
		engine = []string{
			"Baidu",
			"Bing",
			// "Google",
			// "Wx",
		}
	} else {
		engine = strings.Split(engineStr, ",")
	}

	var timeout time.Duration
	if timeoutStr := request.Form.Get("timeout"); timeoutStr != "" {
		timeout, _ = time.ParseDuration(timeoutStr)
	} else {
		timeout = site.MaxTimeout
	}

	if site.GetDebug() {
		log.Printf("查询内容: %s\n", q)
		log.Printf("引擎: %v\n", engine)
		log.Printf("超时: %v\n", timeout)
	}

	start := time.Now().UnixNano()
	jsonResult := &site.JsonResult{Code: 200, Data: &site.EntityList{
		Index: 0,
		Size:  0,
		List:  []site.Entity{},
	}}

	array, unsupported := site.GetByNames(engine, q, page, filter)
	if array == nil {
		w.WriteHeader(400)
		jsonResult.Code = -1
		jsonResult.Msg = fmt.Sprintf("不支持的搜索引擎: %s", unsupported)
		v, _ := json.Marshal(jsonResult)
		_, _ = w.Write(v)
		return
	}

	cLen := len(array)
	c := make(chan *site.EntityList, cLen)
	for _, engine := range array {
		go func(engine site.SearchEngine) {
			c <- engine.Search()
		}(engine)
	}

	results := []*site.EntityList{}
	timeoutAfter := time.After(timeout)
	timeoutError := false

outer:
	for {
		select {
		case result := <-c:
			results = append(results, result)

			if len(result.List) != 0 && site.GetDebug() {
				log.Println("收到: " + result.List[0].From)
			}

			cLen--
			if cLen == 0 {
				close(c)
				break outer
			}
		case <-timeoutAfter:
			timeoutError = true
			break outer
		}
	}

	// processedURLs := make(map[string]bool)
	// urlFromCount := make(map[string]map[string]bool)

	// for _, result := range results {
	// 	for i, entity := range result.List {
	// 		//初始化自然排序
	// 		entity.PositionScore = (len(result.List) - i) * site.GetPositionWeight(entity.From)
	// 		entity.SearchScore = site.GetSearchScore(entity.From)
	// 		entity.DomainScore = site.GetDomainScore(entity.Host)
	// 		entity.Score = entity.PositionScore + entity.SearchScore + entity.DomainScore

			
	// 		// 检查 URL 是否已经存在于 processedURLs 中
	// 		if !processedURLs[entity.Url] {
	// 			// 如果 URL 没有被处理过，添加到 jsonResult.Data.List 中
	// 			jsonResult.Data.List = append(jsonResult.Data.List, entity)
	// 			// 标记 URL 为已处理
	// 			processedURLs[entity.Url] = true
	// 			jsonResult.Data.Size += 1
	// 		}
	// 	}
	// }

	processedEntities := make(map[string]site.Entity) // Store entities by URL

	for _, result := range results {
		for i, entity := range result.List {
			// Initialize scoring
			positionScore := (len(result.List) - i) * site.GetPositionWeight(entity.From)
			searchScore := site.GetSearchScore(entity.From)
			domainScore := site.GetDomainScore(entity.Host)
			initialScore := positionScore + searchScore + domainScore

			if existingEntity, exists := processedEntities[entity.Url]; exists {
				// Add new source to existing entity's From list, if it's not already there
				entity_list := strings.Split(existingEntity.From, ",")
				if !contains(entity_list, entity.From) {
					existingEntity.From = strings.Join(append(entity_list, entity.From), ",")
					// Update score with additional weight
					existingEntity.Score += initialScore
					processedEntities[entity.Url] = existingEntity
				}
			} else {
				entity.Score = initialScore
				processedEntities[entity.Url] = entity
			}
		}
	}

	// Update jsonResult.Data.List with processedEntities values
	for _, entity := range processedEntities {
		jsonResult.Data.List = append(jsonResult.Data.List, entity)
		jsonResult.Data.Size += 1
	}

	// sort score
	sort.Sort(EntityWrapper{jsonResult.Data.List, func(p, q *site.Entity) bool {
		// Score 递减排序
		return p.Score > q.Score
	}})

	// 构造返回
	jsonResult.Cost = (time.Now().UnixNano() - start) / 1e6
	if timeoutError {
		jsonResult.Code = 202
	}
	body, err := json.Marshal(jsonResult)
	if err != nil {
		jsonResult.Code = -1
		jsonResult.Msg = err.Error()
		w.WriteHeader(500)
		v, _ := json.Marshal(jsonResult)
		_, _ = w.Write(v)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}
