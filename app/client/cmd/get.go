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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get page",
	Args:  cobra.ExactArgs(1),
	Short: "http get request",
	Long: `send http get request 

must be logged in to the server`,
	Run: func(cmd *cobra.Command, args []string) {
		get(server, args[0], port)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func get(server, page string, port int) {
	cookie := getCookie()
	client := http.Client{Timeout: time.Second}
	url := fmt.Sprintf("%s:%d/pages/%s", server, port, page)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	cobra.CheckErr(err)
	request.AddCookie(cookie)
	response, err := client.Do(request)
	cobra.CheckErr(err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	cobra.CheckErr(err)
	if response.StatusCode != http.StatusOK {
		fmt.Printf("status error %s %s", response.Status, string(body))
		return
	}
	fmt.Println("ip address is", string(body))
}

func getCookie() *http.Cookie {
	cookie := &http.Cookie{}
	file, err := os.ReadFile(os.TempDir() + "/sshlogin.cookie")
	if err != nil {
		return cookie
	}
	_ = json.Unmarshal(file, &cookie)
	return cookie
}
