package main

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/devilcove/cookie"
	"github.com/mattkasun/sshlogin"
	"golang.org/x/crypto/ssh"
)

const (
	cookieName = "sshlogin"
	cookieAge  = 300
	stringSize = 14
)

var users = map[string][]byte{}

func setupRouter() http.Handler {
	if err := cookie.New(cookieName, cookieAge); err != nil {
		slog.Error("cookie error", "error", err)
		return nil
	}
	router := http.NewServeMux()
	router.HandleFunc("/{$}", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "hello world")
	})
	router.HandleFunc("GET /hello", hello)
	router.HandleFunc("POST /login", login)
	router.HandleFunc("POST /register", register)
	restricted := http.NewServeMux()
	router.Handle("/pages/", http.StripPrefix("/pages", auth(restricted)))
	restricted.HandleFunc("GET /ip", ip)
	restricted.HandleFunc("POST /lines", lines)
	return Logger(router)
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := cookie.Get(r, cookieName); err != nil {
			slog.Error("cookie", "error", err)
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func hello(w http.ResponseWriter, _ *http.Request) {
	b := make([]byte, stringSize)
	rand.Read(b)
	w.Write(b)
}

func login(w http.ResponseWriter, r *http.Request) {
	var login sshlogin.Login
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read login data "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &login); err != nil {
		http.Error(w, "unable to read login data "+err.Error(), http.StatusBadRequest)
		return
	}
	pub, ok := users[login.User]
	if !ok {
		http.Error(w, "no such user", http.StatusBadRequest)
		return
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pub)
	if err != nil {
		http.Error(w, "unable to parse key "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := pubKey.Verify([]byte(login.Message), &login.Sig); err != nil {
		http.Error(w, "unable to verify sig "+err.Error(), http.StatusUnauthorized)
		return
	}
	if err := cookie.Save(w, cookieName, []byte(login.User)); err != nil {
		slog.Error("cookie save", "error", err)
	}
	io.WriteString(w, "login successful")
}

func register(w http.ResponseWriter, r *http.Request) {
	var reg sshlogin.Registration
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read registration data "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &reg); err != nil {
		http.Error(w, "unable to read registration data "+err.Error(), http.StatusBadRequest)
		return
	}
	_, ok := users[reg.User]
	if ok {
		http.Error(w, "username is taken", http.StatusBadRequest)
	}
	users[reg.User] = reg.Key
	io.WriteString(w, "registration successful")
}

func ip(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, r.RemoteAddr)
}

func lines(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read data "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "unable to read data "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(body)
}
