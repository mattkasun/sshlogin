# sshlogin
POC for application login using ssh keys

client and server code to demonstrate using ssh keys for login
# Steps:
### Registration
    a. client registers with server by calling register endpoint with username and ssh publickey
    b. server stores user data (a simple map for POC only)
### Login
   a. client calls hello endpoint and receives random string to sign
   b. client signs the string and calls login endpoint with string, ssh signature, and username
   c. server verifies signature and returns session cookie
### Application Enpoints
  a. client includes session cookie for calls to application endpoints

# Server Cmd
````
A server which supports user registeration and login via ssh

Usage:
  server [flags]

Flags:
      --config string   config file (default is $HOME/.server.yaml)
  -h, --help            help for server
  -p, --port int        server port (default 8080)
````
# Client Cmd
````
client code to register/login app server using ssh

Usage:
  client [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  get         http get request
  help        Help about any command
  login       A brief description of your command
  post        send post request to server
  register    register user with server
  version     show version information

Flags:
      --config string   config file (default is $HOME/.client.yaml)
  -h, --help            help for client
      --port int        port to connect to (default 8080)
      --server string   server url (default "http://localhost")
````
## register
````
register user with server

Usage:
  client register username [flags]

Flags:
  -h, --help         help for register
  -k, --key string   path to public key relative to $HOME/.ssh (default "id_ed25519.pub")

Global Flags:
      --config string   config file (default is $HOME/.client.yaml)
      --port int        port to connect to (default 8080)
      --server string   server url (default "http://localhost")
````
## login
````
login to an app server with ssh

        a string is retrieved from server, signed with ssh private key;
string, signature and user name is forwarded to server which will
return a session cookie for calls to protected endpoints

Usage:
  client login user [flags]

Flags:
  -h, --help         help for login
  -k, --key string   name of private ssh key: relative to $HOME/.ssh (default "id_ed25519")

Global Flags:
      --config string   config file (default is $HOME/.client.yaml)
      --port int        port to connect to (default 8080)
      --server string   server url (default "http://localhost")
````
## get
```
send http get request

Ex: ./client get ip
must be logged in to the server

Usage:
  client get page [flags]

Flags:
  -h, --help   help for get

Global Flags:
      --config string   config file (default is $HOME/.client.yaml)
      --port int        port to connect to (default 8080)
      --server string   server url (default "http://localhost")
```
## post
````
send post request to server specifying page and data in form of key value pairs

        Ex: ./server post lines hello world

Usage:
  client post page [key value ...] [flags]

Flags:
  -h, --help   help for post

Global Flags:
      --config string   config file (default is $HOME/.client.yaml)
      --port int        port to connect to (default 8080)
      --server string   server url (default "http://localhost")
````
 # Example session
 ````
~/ssh-login/app/client (master)> ./client register aUser
registration successfull for aUser
~/ssh-login/app/client (master)> ./client login aUser
login successful
~/ssh-login/app/client (master)> ./client get ip
ip address is 127.0.0.1:38542
~/ssh-login/app/client (master)> ./client post lines hello world
[hello world]
{"hello":"world"}
````
