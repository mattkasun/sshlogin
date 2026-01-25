package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func post(server string, port int, args []string) {
	fs := flag.NewFlagSet("port", flag.ExitOnError)
	fs.Usage = func() {
		usage(fs, "post", "path", "send post to server", os.Args[0]+" post lines hello world")
	}
	fs.Parse(args)

	cookie := getCookie()
	url := fmt.Sprintf("%s:%d/pages/%s", server, port, args[0])
	data := map[string]string{}
	for i := range args[1:] {
		if i%2 == 0 {
			continue
		}
		data[args[i]] = args[i+1]
	}
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal data", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("http request", http.MethodPost, url, err)
		return
	}
	request.AddCookie(cookie)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("http", url, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("status error:", response.Status)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("decode body", err)
		return
	}
	fmt.Println(string(body))
}
