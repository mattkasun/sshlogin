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
	"time"

	"github.com/mattkasun/sshlogin"
	"github.com/spf13/cobra"
)

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post page line1 line2",
	Args:  cobra.ExactArgs(3),
	Short: "send post request to server",
	Long: `send post request to server
specifying page and json data`,
	Run: func(cmd *cobra.Command, args []string) {
		post(server, args[0], args[1], args[2], port)
	},
}

func init() {
	rootCmd.AddCommand(postCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func post(server, page, line1, line2 string, port int) {
	cookie := getCookie()
	client := http.Client{Timeout: time.Second}
	url := fmt.Sprintf("%s:%d/pages/%s", server, port, page)
	data := sshlogin.Data{
		Line1: line1,
		Line2: line2,
	}
	payload, err := json.Marshal(data)
	cobra.CheckErr(err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
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
	returnData := sshlogin.Data{}
	err = json.Unmarshal(body, &returnData)
	cobra.CheckErr(err)
	fmt.Println(string(body))
}
