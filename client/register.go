package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mattkasun/sshlogin"
)

func usage(fs *flag.FlagSet, command, description, example string) {
	fmt.Println(description)
	fmt.Printf("Usage:\n\n")
	fmt.Printf("%s [global flags] %s [flags] <username>\n", os.Args[0], command)
	fmt.Println("\nGlobal Flags:")
	flag.PrintDefaults()
	fmt.Println("\nCommand Flags")
	fs.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println(example)
}

func register(server string, port int, args []string) {
	fs := flag.NewFlagSet("register", flag.ExitOnError)
	key := fs.String("k", "id_ed25519.pub", "name of public ssh key (relative to $HOME/.ssh/)")
	fs.Usage = func() {
		usage(fs, "register", "register user with server", os.Args[0]+" register myName")
	}
	fs.Parse(args)
	if len(fs.Args()) != 1 {
		fs.Usage()
		return
	}
	username := fs.Arg(0)
	pub, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/" + *key)
	if err != nil {
		fmt.Println("unable to read file", *key, err)
		return
	}
	register := sshlogin.Registration{
		User: username,
		Key:  pub,
	}
	payload, err := json.Marshal(register)
	if err != nil {
		fmt.Println("marshal payload:", err)
		return
	}
	url := fmt.Sprintf("%s:%d/register", server, port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("http request", err)
		return
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("http", http.MethodPost, url, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("status err", url, response.Status)
		return
	}
	fmt.Println("registration successful for user", username, "with ssh pub key", *key)
}
