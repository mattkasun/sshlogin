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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register username",
	Args:  cobra.ExactArgs(1),
	Short: "register user with server",
	Long: `register user with server

	server name can be specified with -s --server flag
	server port can be specified with -p --port flag`,

	Run: func(cmd *cobra.Command, args []string) {
		key, err := cmd.Flags().GetString("pubkey")
		cobra.CheckErr(err)
		server, err := cmd.Flags().GetString("server")
		cobra.CheckErr(err)
		port, err := cmd.Flags().GetInt("port")
		cobra.CheckErr(err)
		if err := register(args[0], server, key, port); err != nil {
			fmt.Println("registation failed")
			cobra.CheckErr(err)
		}
	},
}

func register(user, server, key string, port int) error {
	c := &http.Client{Timeout: 1 * time.Second}
	pub, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/" + key)
	if err != nil {
		return fmt.Errorf("read pub key %w", err)
	}
	register := struct {
		User string
		Key  string
	}{
		User: user,
		Key:  string(pub),
	}
	payload, err := json.Marshal(register)
	if err != nil {
		return fmt.Errorf("marshal %v", err)
	}
	url := fmt.Sprintf("%s:%d/register", server, port)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}
	resp, err := c.Do(request)
	if err != nil {
		return fmt.Errorf("post %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status error %s %s ", resp.Status, string(body))
	}
	fmt.Println(string(body), "for", user)
	return nil
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	registerCmd.Flags().StringP("pubkey", "k", "id_ed25519.pub", "path to public key relative to $HOME/.ssh")
}
