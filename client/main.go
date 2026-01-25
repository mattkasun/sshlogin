package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	defaultPort = 8080
)

func main() {
	log.SetFlags(log.Lshortfile)
	port := flag.Int("p", defaultPort, "server listening port")
	server := flag.String("s", "http://localhost", "server url")

	flag.Usage = func() {
		fmt.Println("client to register/login/interact with app server using ssh")
		fmt.Println("\nUsage:")
		fmt.Printf("%s [flags] command\n", os.Args[0])
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		fmt.Println("\nCommands:\n",
			"\tget\t\thttp get request\n",
			"\tlogin\t\tserver login using ssh key\n",
			"\tpost\t\tsend post request to server\n",
			"\tregister\tserver registration user with server",
		)
		fmt.Println("\nFor help about a command\n\t", os.Args[0], "[command] -h")
	}
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	switch cmd := flag.Args()[0]; cmd {
	case "get":
		get(*server, *port, flag.Args()[1:])
	case "login":
		login(*server, *port, flag.Args()[1:])
	case "post":
		post(*server, *port, flag.Args()[1:])
	case "register":
		register(*server, *port, flag.Args()[1:])
	default:
		flag.CommandLine.Usage()
	}
}
