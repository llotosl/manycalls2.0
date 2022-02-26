package main

import (
	"fmt"
	"github.com/llotosl/manycalls2.0/pkg/services"
)

func main() {
	mru := services.NewMailRu("325342fsd")
	err := mru.Call("79173642794")
	if err != nil {
		log.Fatal(err)
	}
}
