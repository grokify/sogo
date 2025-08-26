package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grokify/sogo/database/kvs"
	"github.com/grokify/sogo/database/kvs/redis"
)

func main() {
	client := redis.NewClient(kvs.Config{
		Host:        "127.0.0.1",
		Port:        6379,
		Password:    "",
		CustomIndex: 0})

	key := "hello"

	for i, val := range []string{"world", "monde", "世界", "ప్రపంచ"} {
		err := client.SetString(context.Background(), key, val)
		if err != nil {
			slog.Error(err.Error())
		} else {
			slog.Info("successful write",
				"key", i+1,
				"set", val,
				"get", client.GetOrDefaultString(context.Background(), key, ""),
				"is_equal", (val == client.GetOrDefaultString(context.Background(), key, "")),
			)
			/*
				fmt.Printf("(%v) KEY [%v] SET [%v] GET [%v] EQ [%v]\n",
					i+1,
					key,
					val,
					client.GetOrEmptyString(key),
					val == client.GetOrEmptyString(key))
			*/
		}
	}

	fmt.Println("DONE")
}
