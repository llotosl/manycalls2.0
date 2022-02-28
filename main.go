package main

import (
	"fmt"
	"html/template"
	"log"
	"manycalls/pkg/services"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/main/index.html")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	phonesStr := r.PostFormValue("phones")
	phones := strings.Split(phonesStr, " ")
	captchaToken := r.PostFormValue("captchaToken")
	proxy := r.PostFormValue("proxy")
	fmt.Println(phonesStr)

	if phonesStr != "" {
		mru := services.NewMailRu(captchaToken)
		for i := range phones {
			go mru.Call(phones[i], proxy, strconv.Itoa(i))
			time.Sleep(1 * time.Second)
		}
	}

	tmpl.Execute(w, nil)
}

func handleRequest() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Launched")
	handleRequest()
}
