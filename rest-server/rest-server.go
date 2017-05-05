package main

    import (
        "fmt"
        "golang-sample/httprouter"
        "log"
        "net/http"
		"os"
		"time"
		"encoding/json"
    )
	
	type Article struct {
		Id 		string		`json:"id"`
		Title 	string		`json:"title"`
		Date 	string		`json:"date"`
		Body 	string		`json:"body"`
		Tags 	[]string	`json:"tags"`
	}
	
	// Article cache
	type ArticleList struct {
		Articles []Article
	}
	
	var articleList ArticleList
	// Map keyed with tag contains map keyed with date containing article ID list
	var tagDateMap map[string]map[string][]string
	
	type TagSummary struct {
		Tag 			string		`json:"tag"`
		Count 			int			`json:"count"`
		Articles 		[]string	`json:"articles"`
		RelatedTags 	[]string	`json:"related_tags"`
	}

	// Endpoint handler for POST /articles
    func addArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		var a Article
		err := decoder.Decode(&a)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//fmt.Fprintf(w, "incorrect request data %s", err.Error)
			fmt.Println(err)
			return
		}
		defer r.Body.Close()
		
		addArticleToTagMap(a)
		articleList.Articles = append(articleList.Articles, a)
		
        fmt.Fprintf(w, "added article %s", a.Id)
		fmt.Printf("added article %v\n", a)
    }


	// Endpoint handler for GET /articles/{id}
    func getArticles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        id := ps.ByName("id")
		a, ok := getArticle(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			//fmt.Fprintf(w, "article %s not found", id)
			return
		}

		b, err := json.Marshal(a)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Error marshalling json:", err)
			return
		}
		
		fmt.Fprintf(w, string(b))
    }

	// Endpoint handler for GET /tag/{tagName}/{date}
	// Returns count stats, article ID list, and related tags
    func getTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        tag := ps.ByName("tagName")
		date := ps.ByName("date")
		t, err := time.Parse("20060102", date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//fmt.Fprintf(w, "invalid date %s", date)
			return
		}
		d := t.Format("2006-01-02")
        fmt.Printf("looking for tag %s on %s\n", tag, d)
		
		dateMap := tagDateMap[tag]
		fmt.Printf("found %v dates for tag\n", len(dateMap))
		
		var idList []string
		ids := dateMap[d]
		// we only want the last 10 articles
		if len(ids) > 10 {
			idList = ids[len(ids)-10:]
		} else {
			idList = ids
		}
		if idList == nil {
			idList = make([]string,0)
		}
		fmt.Printf("found %v IDs for date\n", len(idList))

		// Populate a TagSummary struct, then make it JSON
		var summary TagSummary
		summary.Tag = tag
		summary.Count = len(idList)
		summary.Articles = idList
		rtSlice := make([]string, 0, 10)
		
		// collect unique related tags
		rtMap := make(map[string]bool) // map key ensures uniqueuness
		for _, id := range idList {
			a, ok := getArticle(id)
			if ok {
				for _, tag1 := range a.Tags {
					// ignore the search tag
					if tag1 != tag {
						rtMap[tag1] = true
					}
				}
			}
		}
		for k := range rtMap {
			rtSlice = append(rtSlice, k)
		}
		summary.RelatedTags = rtSlice
		
		b, err := json.Marshal(summary)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error marshalling json:", err)
			return
		}
		
		fmt.Fprintf(w, string(b))
		
    }
	
	// Loads the existing article list from a text file
	func loadArticles() ArticleList {
		var articles ArticleList
		articleFile := "articles.json"
		fin, err := os.Open(articleFile) // read only. That's all this needs
		if err != nil {
			fmt.Println(articleFile, err)
			return articles
		}
		defer fin.Close()
		buf := make([]byte, 4096)
		n, _ := fin.Read(buf)
		if n == 0 {
			fmt.Println("no articles in file")
			return articles
		}
		err = json.Unmarshal(buf[:n], &articles)
		if err != nil {
			fmt.Println("error parsing articles file: ", err)
			fmt.Println(buf[:n])
			return articles
		}
		
		fmt.Printf("loaded %v articles\n", len(articles.Articles))
		//if len(articles.Articles) > 0 {
		//	fmt.Printf("art 1: %v\n", articles.Articles[0])
		//}
		
		// Assemble tag-date map
		tagDateMap = make(map[string]map[string][]string)
		for _, a := range articles.Articles {
			addArticleToTagMap(a)
		}
		return articles
	}
	
	// Adds a single article to the tag/date map
	func addArticleToTagMap(a Article) {
		for _, t := range a.Tags {
			dMap, ok := tagDateMap[t]
			if !ok { // tag not already there
				dMap = map[string][]string{}
				tagDateMap[t] = dMap
			}
			idList, ok := dMap[a.Date]
			if !ok {
				// date not already there
				idList = make([]string, 0, 10)
				dMap[a.Date] = idList
			}
			idList = append(idList, a.Id)
			tagDateMap[t][a.Date] = idList
		}
	}
	
	// Finds an article by ID, from the local cache. If not found, return ok=false.
	func getArticle(id string) (a Article, ok bool) {
		for _, a := range articleList.Articles {
			if a.Id == id {
				return a, true
			}
		}
		ok = false
		return 
	}

    func main() {
        router := httprouter.New()

        router.POST("/articles", addArticle)
        router.GET("/articles/:id", getArticles)
		router.GET("/tag/:tagName/:date", getTags)
		
		articleList = loadArticles()

        log.Fatal(http.ListenAndServe(":8080", router))
    }