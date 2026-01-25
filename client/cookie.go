package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var dir string

func getCookie() *http.Cookie {
	cookie := &http.Cookie{}
	file, err := os.ReadFile(dir + "/cookie")
	if err != nil {
		return cookie
	}
	_ = json.Unmarshal(file, &cookie)
	return cookie
}

func saveCookie(cookies []*http.Cookie) {
	if dir == "" {
		var err error
		dir, err = os.MkdirTemp("", "sshlogin")
		if err != nil {
			log.Println("create temp dir", err)
			return
		}
	}
	found := false
	for _, c := range cookies {
		if c.Name == "sshlogin" {
			found = true
			cookie, err := json.Marshal(*c)
			if err != nil {
				fmt.Println("marshal cookie", err)
				return
			}
			os.WriteFile(dir+"/cookie", cookie, 0)
			break
		}
	}
	if !found {
		fmt.Println("no cookie from server")
		return
	}
}
