package sshlogin

import "golang.org/x/crypto/ssh"

// Login contains login information.
type Login struct {
	Message string        `json:"message"`
	Sig     ssh.Signature `json:"sig"`
	User    string        `json:"user"`
}

// Registration contains information to register a new user.
type Registration struct {
	User string `json:"user"`
	Key  []byte `json:"key"`
}

// Data contains information to post to server.
type Data struct {
	Line1 string
	Line2 string
}
