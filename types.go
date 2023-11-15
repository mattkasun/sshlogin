package sshlogin

import "golang.org/x/crypto/ssh"

// Login contains login information
type Login struct {
	Message string        `json:"message"`
	Sig     ssh.Signature `json:"sig"`
	User    string        `json:"user"`
}

// Registration contains information to register a new user
type Registation struct {
	User string `json:"user"`
	Key  string `json:"key"`
}

// Data contains infomation to post to server
type Data struct {
	Line1 string
	Line2 string
}
