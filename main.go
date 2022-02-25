package main

import (
	"net/http"

	"./pkg/requests"
)

func main() {
	client := &http.Client{}

	request := requests.NewRequest(client)
}
