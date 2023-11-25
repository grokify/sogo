package main

import (
	"fmt"
	"log"

	"github.com/grokify/xgo/database/document"
	"github.com/grokify/xgo/database/document/redis"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Key string `short:"k" long:"key" description:"Storage key" required:"true"`
}

func main() {
	opts := &Options{}
	_, err := flags.Parse(opts)
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(document.Config{
		Host:        "127.0.0.1",
		Port:        6379,
		Password:    "",
		CustomIndex: 0})

	data := client.GetOrEmptyString(opts.Key)

	fmt.Printf("Data: %v\n", data)

	fmt.Println("DONE")
}
