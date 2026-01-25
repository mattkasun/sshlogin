package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func get(server string, port int, args []string) {
	// log.Println("get args", args)
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	fs.Usage = func() {
		usage(fs, "get", "path", "send request to server", os.Args[0]+" get ip")
	}
	fs.Parse(args)
	// check for additional flags
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			fs.Parse(args[i : i+1])
			// fs.Parse(args) // reset fs.Args()
		}
	}
	if len(fs.Args()) != 1 {
		log.Println("fs args", len(fs.Args()), fs.Args())
		fs.Usage()
		log.Println("no args or wrong args", fs.Args())
		return
	}
	// log.Println("call get page")
	// getPage(server, fs.Args()[0], port)
	// }

	// func getPage(server, page string, port int) {
	cookie := getCookie()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// url := fmt.Sprintf("%s:%d/pages/%s", server, port, page)
	url := fmt.Sprintf("%s:%d/pages/%s", server, port, fs.Args()[0])
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.AddCookie(cookie)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if response.StatusCode != http.StatusOK {
		fmt.Printf("status error %s %s", response.Status, string(body))
		return
	}
	fmt.Println("ip address is", string(body))
}
