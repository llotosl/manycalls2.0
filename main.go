package main

import (
	"fmt"
	"log"
	"net/http"

	"./pkg/requests"
)

func main() {
	client := &http.Client{}
	request := requests.NewRequest(client)
	headers := map[string]string{}

	body, _, err := request.Get("https://google.com", headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
