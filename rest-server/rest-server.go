package main

    import (
        "fmt"
        "golang-sample/github.com/julienschmidt/httprouter"
        "log"
        "net/http"
		"os"
		"encoding/json"
    )
	
	type Article struct {
		id int
		title string
		date string
		body string
		tags []string
	}
	
	type ArticleList struct {
		articles []Article
	}
	
	var articleList ArticleList

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

    func adduser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        // uid := r.FormValue("uid")
        uid := ps.ByName("uid")
        fmt.Fprintf(w, "you are add user %s", uid)
    }

    func getTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        uid := ps.ByName("tagName")
		date := ps.ByName("date")
		
		
        fmt.Fprintf(w, "you are get user %s", uid)
    }
	
	func loadArticles() ArticleList {
		articleFile := "articles.json"
		fin, err := os.Open(articleFile) // read only. That's all this needs
		if err != nil {
			fmt.Println(articleFile, err)
			return nil
		}
		defer fin.Close()
		buf := make([]byte, 4096)
		var articles articleList 
		n, _ := fin.Read(buf)
		err = json.Unmarshal(buf, &articles)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		
		fmt.Println("articles", len(articles))
		return articles
	}

    func main() {
        router := httprouter.New()
        router.GET("/", Index)
        router.GET("/hello/:name", Hello)

        router.GET("/user/:uid", getuser)
        router.POST("/adduser/:uid", adduser)
        router.DELETE("/deluser/:uid", deleteuser)
        router.PUT("/moduser/:uid", modifyuser)
		router.GET("/tag/{tagName}/{date}", getTags)
		
		articleList = loadArticles()

        log.Fatal(http.ListenAndServe(":8080", router))
    }