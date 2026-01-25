package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const cookiefile = "/sshlogin.cookie"

func getCookie() *http.Cookie {
	cookie := &http.Cookie{}
	file, err := os.ReadFile(os.TempDir() + cookiefile)
	if err != nil {
		fmt.Println("cookie error", err)
		return cookie
	}
	_ = json.Unmarshal(file, cookie)
	return cookie
}

func saveCookie(cookies []*http.Cookie) {
	found := false
	for _, c := range cookies {
		if c.Name == "sshlogin" {
			found = true
			cookie, err := json.Marshal(*c)
			if err != nil {
				fmt.Println("marshal cookie", err)
				return
			}
			os.WriteFile(os.TempDir()+cookiefile, cookie, 0)
			break
		}
	}
	if !found {
		fmt.Println("no cookie from server")
		return
	}
}
