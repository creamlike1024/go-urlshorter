package main

import (
	"fmt"
	"html/template"
	"main/jsonedit"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	gonanoid "github.com/matoous/go-nanoid/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env file")
	}
	_, err = strconv.Atoi(os.Getenv("LENGTH"))
	if err != nil {
		log.Fatal("Error parsing LENGTH")
	}
	if _, err := os.Stat(os.Getenv("DB_FILE")); err != nil {
		if os.IsNotExist(err) {
			jsonedit.InitJson(os.Getenv("DB_FILE"))
		}
	}

	log.Info("Server started on port " + os.Getenv("PORT"))
	// Add Route
	http.HandleFunc("/", pathHandler)
	http.HandleFunc("/"+os.Getenv("MANAGE_PATH"), manageHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	jsonFileContent, err := os.ReadFile(os.Getenv("DB_FILE"))
	if err != nil {
		log.Fatal("Error reading json file")
	}
	pathUrls := jsonedit.ParseJsonFile(jsonFileContent)
	for _, pathUrl := range pathUrls.Db {
		if r.URL.Path == pathUrl.Path {
			http.Redirect(w, r, pathUrl.URL, http.StatusSeeOther)
			return
		}
	}
	defaultHandler(w, r)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "无效的路径")
}

func manageHandler(w http.ResponseWriter, r *http.Request) {
	// 如果是 GET 请求，显示添加路径页面
	if r.Method == "GET" {
		t, _ := template.ParseFiles("manage.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		if r.Form["password"][0] == os.Getenv("PASSWORD") {
			if r.Form["url"][0] != "" && r.Form["path"][0] == "" {
				addPathHandler(w, r, r.Form["url"][0])
			} else if r.Form["url"][0] == "" && r.Form["path"][0] != "" {
				delPathHandler(w, r, r.Form["path"][0])
			} else {
				fmt.Fprintf(w, "你不能同时填写要添加的 Url 和要删除的 Path")
			}
		} else {
			fmt.Fprintf(w, "密码错误")
		}
	}
}

func addPathHandler(w http.ResponseWriter, r *http.Request, u string) {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		fmt.Fprintf(w, "URL 格式错误")
		return
	}
	len, _ := strconv.Atoi(os.Getenv("LENGTH"))
	Id := genId(len)
	jsonFileContent, err := os.ReadFile(os.Getenv("DB_FILE"))
	if err != nil {
		log.Fatal("Error reading json file")
	}
	for {
		if !jsonedit.CheckPath(Id, &jsonFileContent) {
			break
		}
	}
	Id = "/" + Id
	jsonedit.AddPath(Id, u, os.Getenv("DB_FILE"))
	fmt.Fprintf(w, "添加成功，路径为："+os.Getenv("DOMAIN")+Id)
	log.Info("Path:" + Id + " added")
}

func delPathHandler(w http.ResponseWriter, r *http.Request, p string) {
	err := jsonedit.DelPath("/"+p, os.Getenv("DB_FILE"))
	if err != nil {
		fmt.Fprintf(w, "路径不存在")
	} else {
		fmt.Fprintf(w, "删除成功")
		log.Info("Path:" + p + " deleted")
	}
}

func genId(len int) string {
	id, _ := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", len)
	return id
}
