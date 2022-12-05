package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bhuvana-chinnadurai/cli-http-workerpool-go/client"
	"github.com/bhuvana-chinnadurai/cli-http-workerpool-go/workerpool"
)

func main() {
	var urls []string
	var maxParallelCount int
	flag.Func("urls", "List of URLs to make http requests", func(flagValue string) error {
		urls = strings.Fields(flagValue)
		return nil
	})
	flag.Int("parallel", maxParallelCount, "No of Parallel requests that can be made")

	flag.Parse()

	if len(urls) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	wp := workerpool.New(maxParallelCount)

	go func() {
		for _, url := range urls {
			wp.URLs <- url
		}
		close(wp.URLs)
	}()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx, client.New(ctx))

	for r := range wp.Results {
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			continue
		}
		if result, ok := r.Value.(*client.Result); ok && result != nil {
			fmt.Printf(" \n \t '%s':\t  '%s'\t  ", result.URLAddress, result.MD5Hash)
		}
	}

}
