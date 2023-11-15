/*
Copyright © 2023 Matthew R Kasun <mkasun@nusak.ca>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	sshlogin "github.com/mattkasun/ssh-login"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login user",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := cmd.Flags().GetString("key")
		cobra.CheckErr(err)
		message, err := hello(server, port)
		cobra.CheckErr(err)
		err = login(args[0], server, key, message, port)
		cobra.CheckErr(err)
		fmt.Println("login successful")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringP("key", "k", "id_ed25519", "name of private ssh key: relative to $HOME/.ssh")
}

func hello(server string, port int) ([]byte, error) {
	empty := []byte{}
	url := fmt.Sprintf("%s:%d", server, port)
	c := http.Client{Timeout: time.Second * 1}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return empty, fmt.Errorf("http request %w", err)
	}
	resp, err := c.Do(request)
	if err != nil {
		return empty, fmt.Errorf("get %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return empty, fmt.Errorf("read body %w", err)
	}
	return body, nil
}

func login(user, server, key string, message []byte, port int) error {
	client := http.Client{Timeout: time.Second}
	private, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/" + key)
	if err != nil {
		return fmt.Errorf("read key %w", err)
	}
	signer, err := ssh.ParsePrivateKey(private)
	if err != nil {
		return fmt.Errorf("parse private key %w", err)
	}
	sig, err := signer.Sign(rand.Reader, message)
	if err != nil {
		return fmt.Errorf("sign %w", err)
	}
	login := sshlogin.Login{
		Message: string(message),
		Sig:     *sig,
		User:    user,
	}
	payload, err := json.Marshal(login)
	if err != nil {
		return fmt.Errorf("marshal %w", err)
	}
	url := fmt.Sprintf("%s:%d/login", server, port)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("http request %w", err)
	}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("post %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status error %s %s", resp.Status, string(body))
	}
	found := false
	for _, c := range resp.Cookies() {
		if c.Name == "sshlogin" {
			found = true
			cookie, err := json.Marshal(*c)
			if err != nil {
				return fmt.Errorf("cookie error %w", err)
			}
			if err := saveCookie(cookie); err != nil {
				return fmt.Errorf("save cookie %w", err)
			}
		}
	}
	if !found {
		return fmt.Errorf("server did not return cookie")
	}
	return nil
}

func saveCookie(cookie []byte) error {
	return os.WriteFile(os.TempDir()+"/sshlogin.cookie", cookie, os.ModePerm)
}
