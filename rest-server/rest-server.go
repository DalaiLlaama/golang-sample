package main

    import (
        "fmt"
        "golang-sample/github.com/julienschmidt/httprouter"
        "log"
        "net/http"
		"os"
		"time"
		"encoding/json"
    )
	
	type Article struct {
		Id 		string
		Title 	string
		Date 	string
		Body 	string
		Tags 	[]string
	}
	
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

    func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	    fmt.Println("path", r.URL.Path)
        fmt.Fprint(w, "Welcome!\n")
    }

    func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
    }

    func getuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        uid := ps.ByName("uid")
        fmt.Fprintf(w, "you are get user %s", uid)
    }


    func modifyuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        uid := ps.ByName("uid")
        fmt.Fprintf(w, "you are modify user %s", uid)
    }

    func deleteuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        uid := ps.ByName("uid")
        fmt.Fprintf(w, "you are delete user %s", uid)
    }

    func addArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Printf("body %v", r.Body)
		decoder := json.NewDecoder(r.Body)
		var a Article
		err := decoder.Decode(&a)
		if err != nil {
			fmt.Println("error parsing request body", err)
			fmt.Fprintf(w, "incorrect request data %s", err.Error)
		}
		defer r.Body.Close()
		
		addArticleToTagMap(a)
		
        fmt.Fprintf(w, "added article %s", a.Id)
		fmt.Printf("added article %v/n", a)
    }

    func getTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        tag := ps.ByName("tagName")
		date := ps.ByName("date")
		t, err := time.Parse("20060102", date)
		if err != nil {
			fmt.Fprintf(w, "invalid date %s", date)
			return
		}
		d := t.Format("2006-01-02")
        fmt.Printf("looking for tag %s on %s\n", tag, d)
		
		dateMap := tagDateMap[tag]
		fmt.Printf("found %v dates for tag\n", len(dateMap))
		
		idList := dateMap[d]
		fmt.Printf("found %v IDs for date\n", len(idList))

		var summary TagSummary
		summary.Tag = tag
		summary.Count = len(idList)
		summary.Articles = idList
		rtSlice := make([]string, 0, 10)
		
		// collect unique related tags
		rtMap := make(map[string]bool)
		for _, id := range idList {
			a, ok := getArticle(id)
			if ok {
				for _, tag1 := range a.Tags {
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
			fmt.Println("Error marshalling json:", err)
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
        router.GET("/", Index)
        router.GET("/hello/:name", Hello)

        router.POST("/articles", addArticle)
        router.DELETE("/deluser/:uid", deleteuser)
        router.PUT("/moduser/:uid", modifyuser)
		router.GET("/tag/:tagName/:date", getTags)
		
		articleList = loadArticles()

        log.Fatal(http.ListenAndServe(":8080", router))
    }