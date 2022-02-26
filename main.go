package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/llotosl/manycalls2.0/pkg/requests"
)

func main() {
	client := &http.Client{}
	request := requests.NewRequest(client)
	headers := map[string]string{}

	data, contentType, err := requests.MakeBoundary("fdsQEFFJjffjgHkf", dataHead)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
	fmt.Println(contentType)
}
