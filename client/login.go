package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mattkasun/sshlogin"
	"golang.org/x/crypto/ssh"
)

func login(server string, port int, args []string) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	key := fs.String("k", "id_ed25519", "name of private ssh key (relative to $HOME/.ssh/)")
	fs.Usage = func() {
		usage(fs, "login", "username", "login to an app server using ssh",
			os.Args[0]+" login myName")
	}
	fs.Parse(args)

	if len(fs.Args()) != 1 {
		log.Println("wrong num args")
		fs.Usage()
		return
	}
	signed, err := signMessage(*key)
	if err != nil {
		fmt.Println("sign key", err)
		return
	}
	login := sshlogin.Login{
		Message: "hello world",
		Sig:     *signed,
		User:    fs.Arg(0),
	}
	payload, err := json.Marshal(login)
	if err != nil {
		fmt.Println("marshal payload:", err)
		return
	}
	url := fmt.Sprintf("%s:%d/login", server, port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("http request", err)
		return
	}
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println("post", url, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("bad response", response.Status)
		return
	}
	saveCookie(response.Cookies())
	fmt.Println("login successful")
}

func signMessage(key string) (*ssh.Signature, error) {
	private, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/" + key)
	if err != nil {
		return nil, fmt.Errorf("read ssh key %s %w", key, err)
	}
	signer, err := ssh.ParsePrivateKey(private)
	if err != nil {
		return nil, fmt.Errorf("parse key: %s %w", key, err)
	}
	sig, err := signer.Sign(rand.Reader, []byte("hello world"))
	if err != nil {
		return nil, fmt.Errorf("sign message %s %w", key, err)
	}
	return sig, nil
}
