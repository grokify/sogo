package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"

	"github.com/grokify/sogo/database/kvs"
	"github.com/grokify/sogo/database/kvs/redis"
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

	client := redis.NewClient(kvs.Config{
		Host:        "127.0.0.1",
		Port:        6379,
		Password:    "",
		CustomIndex: 0})

	data := client.GetOrDefaultString(context.Background(), opts.Key, "")

	fmt.Printf("Data: %v\n", data)

	fmt.Println("DONE")
}
