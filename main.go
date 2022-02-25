package main

import (
	"net/http"
)

func main() {
	var client *http.Client
	var clientCookie *http.Client
	client = &http.Client{}

	request := requests.newRequest(client)
}
