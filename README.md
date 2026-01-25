# sshlogin
POC for application login using ssh keys

client and server code to demonstrate using ssh keys for login
# Steps:
### Registration
    a. client registers with server by calling register endpoint with username and ssh publickey
    b. server stores user data (a simple map for POC only)
### Login
   a. client signs a string and calls login endpoint with string, ssh signature, and username
   b. server verifies signature and returns session cookie
### Application Enpoints
  a. client includes session cookie for calls to application endpoints

# Server Cmd
````
A server which supports user registeration and login via ssh

Usage:
  server [flags]

Flags:
  -p, --port int        server port (default 8080)
````
# Client Cmd
````
client to register/login/interact with app server using ssh

Usage:
client [flags] command

Flags:
  -p int
        server listening port (default 8080)
  -s string
        server url (default "http://localhost")

Commands:
        get             http get request
        login           server login using ssh key
        post            send post request to server
        register        server registration user with server

For help about a command
         client [command] -h
````
## register
````
register user with server

Usage:
client [global flags] register [flags] <username>

Global Flags:
  -p int
        server listening port (default 8080)
  -s string
        server url (default "http://localhost")

Command Flags
  -k string
        name of public ssh key (relative to $HOME/.ssh/) (default "id_ed25519.pub")

Examples:
client register myName
````
## login
````
login to an app server using ssh

Usage:
client [global flags] login [flags] <username>

Global Flags:
  -p int
        server listening port (default 8080)
  -s string
        server url (default "http://localhost")

Command Flags
  -k string
        name of private ssh key (relative to $HOME/.ssh/) (default "id_ed25519")

Examples:
client login myName
````
## get
```
send request to server

Usage:
client [global flags] get [flags] <path>

Global Flags:
  -p int
        server listening port (default 8080)
  -s string
        server url (default "http://localhost")

Command Flags

Examples:
client get ip
```
## post
````
send post to server

Usage:
client [global flags] post [flags] <path>

Global Flags:
  -p int
        server listening port (default 8080)
  -s string
        server url (default "http://localhost")

Command Flags

Examples:
client post lines hello world
````
 # Example session
 ````
> client register aUser
registration successfull for aUser
> client login aUser
login successful
> client get ip
ip address is 127.0.0.1:38542
> client post lines hello world
{"hello":"world"}
````
